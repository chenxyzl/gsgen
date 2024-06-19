package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Gen 对外的生成接口
func Gen(dir string, fileSuffix []string, exportSetter bool, exportBson bool, headAnnotations []string, ignoreCheckIdents []string) {
	//读取带处理文件列表
	targetFiles := readFileList(dir, fileSuffix)
	if len(targetFiles) == 0 {
		fmt.Printf("dir:%v not found file with suffix:%v\n", dir, fileSuffix)
		return
	}
	//开始处理
	for _, file := range targetFiles {
		genFile(file, exportSetter, exportBson, headAnnotations, ignoreCheckIdents)
	}
}

// readFileList 读取需要生成的文件列表
func readFileList(dir string, fileSuffix []string) []string {
	var targetFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//只读取当前这层目录
		if info.IsDir() {
			return nil
		}
		//只读取对应后缀的文件
		found := false
		for _, fs := range fileSuffix {
			if strings.HasSuffix(info.Name(), fs) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
		targetFiles = append(targetFiles, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return targetFiles
}

var usedGSModelStruct bool
var usedIgnoreCheckPackage []string

// genFile 生成文件
func genFile(sourceFile string, exportSetter, exportBson bool, headAnnotations []string, ignoreCheckIdents []string) {
	usedGSModelStruct = false
	needDirty := false
	usedIgnoreCheckPackage = nil
	if exportBson {
		needDirty = true
	}
	// 解析源代码
	fileSet := token.NewFileSet()
	srcFile, err := parser.ParseFile(fileSet, sourceFile, nil, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse source fild:%vcode: %v", sourceFile, err))
	}

	genAstFile := &ast.File{Name: srcFile.Name, Decls: []ast.Decl{}}
	bsonAstFile := &ast.File{Name: srcFile.Name, Decls: []ast.Decl{}}

	ast.Inspect(srcFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		//
		switch typ := n.(type) {
		case *ast.GenDecl:
			//if typ.Tok == token.IMPORT {
			//	genAstFile.Decls = append(genAstFile.Decls, typ)
			//	if exportBson {
			//		bsonAstFile.Decls = append(bsonAstFile.Decls, typ)
			//	}
			//}
		case *ast.TypeSpec:
			structType, structTypeOk := typ.Type.(*ast.StructType)
			if !structTypeOk {
				return true
			}
			//检查需要生成的Field
			fields := checkStructField(typ.Name, structType, needDirty, exportBson, ignoreCheckIdents)
			//
			generate(genAstFile, typ.Name, fields, exportSetter, needDirty)
			//bson 开始生成
			if exportBson {
				generateBson(bsonAstFile, typ.Name, fields, exportSetter)
			}
		case *ast.File, *ast.Ident, *ast.ImportSpec: //ignore
		default:
			//panic(fmt.Errorf("文件:%v 不允许有定义model以外的代码(头文件不好处理,解决后可能会放开限制)", srcFile.Name))
		}
		return true
	})

	//gen import
	genImportList := []string{"fmt", "encoding/json"}
	genImportList = append(genImportList, usedIgnoreCheckPackage...)
	if usedGSModelStruct {
		genImportList = append(genImportList, "github.com/chenxyzl/gsgen/gsmodel")
	}
	addImport(genAstFile, genImportList...)

	//bson import
	if exportBson {
		bsonImportLIst := []string{"go.mongodb.org/mongo-driver/bson"}
		bsonImportLIst = append(bsonImportLIst, usedIgnoreCheckPackage...)
		if usedGSModelStruct {
			bsonImportLIst = append(bsonImportLIst, "github.com/chenxyzl/gsgen/gsmodel")
		}
		addImport(bsonAstFile, bsonImportLIst...)
	}

	genFileName := strings.TrimSuffix(sourceFile, ".go") + ".gen.go"
	printOutFile(fileSet, genAstFile, genFileName, headAnnotations)
	fmt.Printf("输出：%v 完成\n", genFileName)

	if exportBson {
		bsonFileName := strings.TrimSuffix(sourceFile, ".go") + ".bson.go"
		printOutFile(fileSet, bsonAstFile, bsonFileName, headAnnotations)
		fmt.Printf("输出：%v 完成\n", bsonFileName)
	}
}

// generate 生成全部
func generate(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field, exportSetter bool, needDirty bool) {
	for idx, field := range fields {
		generateGetters(file, structTypeExpr, field)
		generateSetters(file, structTypeExpr, field, idx, exportSetter, needDirty)
	}
	genString(file, structTypeExpr, fields)
	genJsonMarshal(file, structTypeExpr, fields)
	genJsonUnmarshal(file, structTypeExpr, fields, exportSetter)
	genClone(file, structTypeExpr)
}

// generateClean 生成clean
func generateClean(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	var cleanStructBody []ast.Stmt
	for _, field := range fields {
		if isBasicType1(field.Type) {
			continue
		}
		name := field.Names[0].Name                            //已提前检查
		cleanStructBody = append(cleanStructBody, &ast.IfStmt{ //field设置自己的dirtyIdx
			Cond: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "s." + name},
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("s." + name),
								Sel: ast.NewIdent("CleanDirty"),
							},
							Args: []ast.Expr{},
						},
					},
				},
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
						Args: []ast.Expr{},
					},
				},
			},
				cleanStructBody...,
			),
		},
	})
}
