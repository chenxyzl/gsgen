package test1

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gotest/tools/genmod/test1/internal"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestGetterSetter(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testFile := filepath.Clean(filepath.Join(currentDir, "../../../", "model/example.go"))
	testOutFile := filepath.Clean(filepath.Join(currentDir, "../../../", "model/example.gen.go"))

	// 解析源代码
	fset := token.NewFileSet()
	srcFile, err := parser.ParseFile(fset, testFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse source code: %v", err)
	}

	genFile := &ast.File{
		Name:  srcFile.Name,
		Decls: []ast.Decl{},
	}

	//
	addImport(genFile, []string{})

	// 找到需要生成 getter 和 setter 的 struct 类型
	ast.Inspect(srcFile, func(n ast.Node) bool {
		if spec, ok := n.(*ast.TypeSpec); ok {
			structType, ok := spec.Type.(*ast.StructType)
			if ok {
				//格式检查
				fields := checkModelStruct(spec.Name, structType)
				// 生成 getter 和 setter 方法
				generateGettersAndSetters(genFile, spec.Name, structType, fields)
			}
		}
		return true
	})

	outputFile, err := os.Create(testOutFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// 打印修改后的源代码
	buffer := bytes.NewBuffer(nil)
	addHeader(buffer)
	//格式化 输出到buffer
	err = format.Node(buffer, fset, genFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to format source code: %v", err))
	}
	err = os.WriteFile(testOutFile, buffer.Bytes(), 0444)
	if err != nil {
		panic(fmt.Sprintf("Failed to print modified source code: %v", err))
	}
	//err = printer.Fprint(outputFile, fset, genFile)
	//if err != nil {
	//	log.Fatalf("Failed to print modified source code: %v", err)
	//}
}

func checkModelStruct(structNameIdent *ast.Ident, structType *ast.StructType) (out []*ast.Field) {
	contain := false
	for _, field := range structType.Fields.List {
		isDirtyModel, needGenField := isLegalField(structNameIdent, field)
		if isDirtyModel {
			contain = true
		}
		if needGenField {
			out = append(out, field)
		}
	}
	if !contain {
		panic(fmt.Sprintf("类型:%v, 必须包含DirtyModel", structNameIdent))
	}
	return out
}

func generateGettersAndSetters(file *ast.File, structTypeExpr *ast.Ident, structType *ast.StructType, fields []*ast.Field) *ast.File {
	var cleanStructBody []ast.Stmt
	for idx, field := range fields {
		//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
		fieldName := field.Names[0].Name
		isBaseType := isBasicType1(field.Type)

		//getter
		file.Decls = append(file.Decls, &ast.FuncDecl{
			Name: ast.NewIdent("Get" + cases.Title(language.Und).String(fieldName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: field.Type,
						},
					},
				},
			},
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("s")},
						Type:  &ast.StarExpr{X: structTypeExpr},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("s"),
								Sel: ast.NewIdent(fieldName),
							},
						},
					},
				},
			},
		})

		//setter-body
		var setterBody []ast.Stmt
		if !isBaseType {
			setterBody = []ast.Stmt{
				&ast.IfStmt{ //field设置自己的dirtyIdx
					If:   0,
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "v"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: ast.NewIdent("v.SetParent"),
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.INT,
											Value: strconv.Itoa(idx),
										},
										&ast.BasicLit{
											Kind:  token.FUNC,
											Value: "s.UpdateDirty",
										},
									},
								},
							},
						},
					},
					Else: nil,
				},
				&ast.AssignStmt{ //赋值
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent(fieldName),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("v"),
					},
				},
				&ast.ExprStmt{ //更新当前的dirty
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent("UpdateDirty"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(idx),
							},
						},
					},
				},
			}
			//生成clean-field
			cleanStructBody = append(cleanStructBody, &ast.IfStmt{ //field设置自己的dirtyIdx
				If:   0,
				Init: nil,
				Cond: &ast.BinaryExpr{
					X:  &ast.Ident{Name: "s." + fieldName},
					Op: token.NEQ,
					Y:  &ast.Ident{Name: "nil"},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("s." + fieldName),
									Sel: ast.NewIdent("CleanDirty"),
								},
							},
						},
					},
				},
				Else: nil,
			})
		} else {
			setterBody = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent(fieldName),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("v"),
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s"),
							Sel: ast.NewIdent("UpdateDirty"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(idx),
							},
						},
					},
				},
			}
		}
		//setter
		file.Decls = append(file.Decls, &ast.FuncDecl{
			Name: ast.NewIdent("Set" + cases.Title(language.Und).String(fieldName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("v")},
							Type:  field.Type,
						},
					},
				},
			},
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("s")},
						Type:  &ast.StarExpr{X: structTypeExpr},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: setterBody,
			},
		})
	}
	//生成clean方法
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Name: ast.NewIdent("CleanDirty"),
		Type: &ast.FuncType{},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: append([]ast.Stmt{ //先clean自己
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s.DirtyModel"),
							Sel: ast.NewIdent("CleanDirty"),
						},
					},
				},
			},
				cleanStructBody..., //再clean-field
			),
		},
	})
	return file
}

func addImport(genFile *ast.File, imports []string) {
	// 添加导入语句
	if len(imports) > 0 {
		//importSpecs := make([]*ast.ImportSpec, 0, len(imports))
		importSpecs := make([]ast.Spec, 0, len(imports))
		for _, importPath := range imports {
			importSpecs = append(importSpecs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(importPath),
				},
			})
		}
		genFile.Decls = append(genFile.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		})
	}
}

func isBasicType(typeStr string) bool {
	basicTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"string", "bool",
		"byte", "rune",
	}
	for _, t := range basicTypes {
		if typeStr == t {
			return true
		}
	}
	return false
}
func isBasicType1(expr ast.Expr) bool {
	if ident, identOk := expr.(*ast.Ident); identOk {
		return isBasicType(ident.Name)
	}
	return false
}

func notSupportBasicType(typeStr string) bool {
	basicTypes := []string{
		"map", "slice",
	}
	for _, t := range basicTypes {
		if typeStr == t {
			return true
		}
	}
	return false
}

// getImplType 获取具体类型
// @return 具体类型
// @return 是否指针
func getImplType(expr ast.Expr) (*ast.Ident, bool) {
	var typeName *ast.Ident
	if selectExpr, selectOk := expr.(*ast.SelectorExpr); selectOk {
		return selectExpr.Sel, false
	} else if starExpr, starOk := expr.(*ast.StarExpr); starOk {
		if innerSelectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
			return innerSelectorExpr.Sel, true
		} else {
			return starExpr.X.(*ast.Ident), true
		}
	} else {
		typeName = expr.(*ast.Ident)
		return typeName, false
	}
}

// isLegalField 检查是否合法字段
// 只能是基本类型 MMap MList 或者包含mdata.DirtyModel实现的struct
// @return 是否是DirtyModel
// @return 是否需要生成的field
func isLegalField(structNameIdent *ast.Ident, field *ast.Field) (bool, bool) {
	if len(field.Names) > 1 { //
		panic(fmt.Sprintf("类型:%v,每行只能声明1个filed", structNameIdent))
	} else if len(field.Names) == 0 { //匿名类型必须只能是mdata.DirtyModel
		if selectExpr, selectOk := field.Type.(*ast.SelectorExpr); !selectOk {
			panic(fmt.Sprintf("类型:%v, 匿名字段必须为mdata.DirtyModel类型 field.Type:%v", structNameIdent, field.Type))
		} else if ident, identOk := selectExpr.X.(*ast.Ident); !identOk || ident.Name != "mdata" {
			panic(fmt.Sprintf("类型:%v, 匿名字段必须为mdata.DirtyModel类型 field.selectExpr.X:%v", structNameIdent, selectExpr.X))
		} else if selectExpr.Sel.Name != "DirtyModel" {
			panic(fmt.Sprintf("类型:%v, 匿名字段必须为mdata.DirtyModel类型 field.selectExpr.Sel.Name:%v", structNameIdent, selectExpr.Sel.Name))
		} else {
			return true, false
		}
	} else { //len(field.Names) == 1
		fieldName := field.Names[0].Name
		if ast.IsExported(fieldName) {
			panic(fmt.Sprintf("类型:%v,字段:%v, 必须为非导出的(即小写)", structNameIdent, fieldName))
		}
		if fieldIdent, fieldIdentOk := field.Type.(*ast.Ident); fieldIdentOk { //是*ast.Ident
			if !isBasicType(fieldIdent.Name) {
				panic(fmt.Sprintf("类型:%v,字段:%v, 非基本类型必须为指针类型", structNameIdent, fieldName))
			}
			return false, true
		} else if starExpr, starOk := field.Type.(*ast.StarExpr); starOk { //是*ast.StarExpr,返回具体类型
			if mlist, mlistOk := starExpr.X.(*ast.IndexExpr); mlistOk && isLegalMList(mlist) {
				return false, true
			} else if mmap, mmapOk := starExpr.X.(*ast.IndexListExpr); mmapOk && isLegalMMap(mmap) {
				return false, true
			}
			nestStarExprX, nestStarExprXOk := starExpr.X.(*ast.Ident)
			if !nestStarExprXOk {
				panic(fmt.Sprintf("类型:%v,字段:%v, 指针字段类型必须是当前包内类型,MMList,MMap", structNameIdent, fieldName))
			}
			if isBasicType(nestStarExprX.Name) {
				panic(fmt.Sprintf("类型:%v,字段:%v, 基本类型必须为值类型", structNameIdent, fieldName))
			}
			return false, true
		} else if selectExpr, selectOk := field.Type.(*ast.SelectorExpr); selectOk { //是mdata.MList或者是mdata.MMap
			ident, identOk := selectExpr.X.(*ast.Ident)
			if !identOk || ident.Name != "mdata" {
				panic(fmt.Sprintf("类型:%v,字段:%v, 不能用mdata.MList和mdata.MMap以外的类型", structNameIdent, fieldName))
			}
			if selectExpr.Sel.Name == "DirtyModel" {
				panic(fmt.Sprintf("类型:%v,字段:%v, mdata.DirtyModel必须是匿名字段", structNameIdent, fieldName))
			} else if selectExpr.Sel.Name != "MList" && selectExpr.Sel.Name != "MMap" {
				panic(fmt.Sprintf("类型:%v,字段:%v, 不能用mdata.MList和mdata.MMap以外的类型,当前是:%v", structNameIdent, fieldName, selectExpr.Sel.Name))
			}
			return false, true
		} else if _, mapTypeOk := field.Type.(*ast.MapType); mapTypeOk {
			panic(fmt.Sprintf("类型:%v,字段:%v, map请替换为mdata.MMap", structNameIdent, fieldName))
		} else if _, sliceTypeOk := field.Type.(*ast.ArrayType); sliceTypeOk {
			panic(fmt.Sprintf("类型:%v,字段:%v, slice请替换为mdata.MList", structNameIdent, fieldName))
		} else if mlist, mlistOk := field.Type.(*ast.IndexExpr); mlistOk && isLegalMList(mlist) {
			panic(fmt.Sprintf("类型:%v,字段:%v, mdate.MList必须为指针类型", structNameIdent, fieldName))
		} else if mmap, mmapOk := field.Type.(*ast.IndexListExpr); mmapOk && isLegalMMap(mmap) {
			panic(fmt.Sprintf("类型:%v,字段:%v, mdate.MMap必须为指针类型", structNameIdent, fieldName))
		} else {
			panic(fmt.Sprintf("类型:%v,字段:%v, 不支持的FieldType:%v", structNameIdent, fieldName, field.Type))
		}
	}
}
func isLegalMList(mlist *ast.IndexExpr) bool {
	if listSelectExpr, listSelectExprOk := mlist.X.(*ast.SelectorExpr); !listSelectExprOk {
		return false
	} else if listSelectExprXIdent, listSelectExprXIdentOk := listSelectExpr.X.(*ast.Ident); !listSelectExprXIdentOk {
		return false
	} else if listSelectExprXIdent.Name != "mdata" && listSelectExpr.Sel.Name != "MList" {
		return false
	}
	return true
}
func isLegalMMap(mmap *ast.IndexListExpr) bool {
	if mapSelectExpr, mapSelectExprOk := mmap.X.(*ast.SelectorExpr); !mapSelectExprOk {
		return false
	} else if mapSelectExprXIdent, mapSelectExprXIdentOk := mapSelectExpr.X.(*ast.Ident); !mapSelectExprXIdentOk {
		return false
	} else if mapSelectExprXIdent.Name != "mdata" && mapSelectExpr.Sel.Name != "MMap" {
		return false
	}
	return true
}

//func isLegalType(expr ast.Ident) bool {
//
//}
//
//func isBaseType(expr ast.Ident) bool {
//
//}

func addHeader(buffer *bytes.Buffer) {
	buffer.WriteString("// Code generated by genmod; DO NOT EDIT.\n")
	buffer.WriteString("// genmod version: 0.0.1\n")
	buffer.WriteString(fmt.Sprintf("// generate time: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	buffer.WriteString(fmt.Sprintf("// src code version: %s\n", ""))
	buffer.WriteString(fmt.Sprintf("// src code commit time : %s\n", ""))
}

func TestCopyData(t *testing.T) {
	a := internal.TestS{}
	a.SetX(1)
	b := a
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println("")
	a.SetX(2)
	b.SetX(3)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println("")
}
