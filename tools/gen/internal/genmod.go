package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Gen 对外的生成接口
func Gen(dir string, fileSuffix []string, needSetter bool, needGetter bool) {
	//读取带处理文件列表
	targetFiles := readFileList(dir, fileSuffix)
	if len(targetFiles) == 0 {
		fmt.Printf("dir:%v not found file with suffix:%v\n", dir, fileSuffix)
		return
	}
	//开始处理
	for _, file := range targetFiles {
		genFile(file, needSetter, needGetter)
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

// genFile 生成文件
func genFile(sourceFile string, needSetter bool, needMongo bool) {
	// 解析源代码
	fileSet := token.NewFileSet()
	srcFile, err := parser.ParseFile(fileSet, sourceFile, nil, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse source fild:%vcode: %v", sourceFile, err))
	}
	outFile := strings.TrimSuffix(sourceFile, ".go") + ".gen.go"

	genAstFile := &ast.File{Name: srcFile.Name, Decls: []ast.Decl{}}
	//mongoAstFile := &ast.File{Name: srcFile.Name, Decls: []ast.Decl{}}

	ast.Inspect(srcFile, func(n ast.Node) bool {
		if genDecl, genDeclOk := n.(*ast.GenDecl); genDeclOk { //头文件
			if genDecl.Tok == token.IMPORT {
				genAstFile.Decls = append(genAstFile.Decls, genDecl)
			}
		} else if spec, specOk := n.(*ast.TypeSpec); specOk { //类型定义
			structType, structTypeOk := spec.Type.(*ast.StructType)
			if !structTypeOk {
				return true
			}
			//检查需要生成的Field
			fields := checkStructField(spec.Name, structType, needMongo)
			//开始生成
			generate(genAstFile, spec.Name, fields, needSetter)
		}
		return true
	})
	printOutFile(fileSet, genAstFile, outFile)
}

// generate 生成全部
func generate(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field, needSetter bool) {
	for idx, field := range fields {
		generateGetters(file, structTypeExpr, field)
		if needSetter {
			generateSetters(file, structTypeExpr, field, idx)
		}
	}
	generateClean(file, structTypeExpr)
}

// generateGetters 生成getter
func generateGetters(file *ast.File, structTypeExpr *ast.Ident, field *ast.Field) {
	//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
	fieldName := field.Names[0].Name
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
}

// generateSetters 生成setter
func generateSetters(file *ast.File, structTypeExpr *ast.Ident, field *ast.Field, idx int) {
	//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
	fieldName := field.Names[0].Name
	isBaseType := isBasicType1(field.Type)
	//setter-body
	var setterBody []ast.Stmt
	//不是基本类型先设置value的Parent
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
		}
	}
	//其他通用语句
	setterBody = append(setterBody,
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
		})
	//setter方法体
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

// generateClean 生成clean
func generateClean(file *ast.File, structTypeExpr *ast.Ident) {
	var cleanStructBody []ast.Stmt
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
}

// printOutFile 输出文件
func printOutFile(fileSet *token.FileSet, astFile *ast.File, outPutFile string) {
	// 打印修改后的源代码
	buffer := bytes.NewBuffer(nil)
	addHeader(buffer)
	//格式化 输出到buffer
	if err := format.Node(buffer, fileSet, astFile); err != nil {
		panic(fmt.Sprintf("Failed to format source code: %v", err))
	}
	//打开
	f, err := os.OpenFile(outPutFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file:%v|err:%v", outPutFile, err))
	}
	defer f.Close()
	//写入内容
	if _, err := f.Write(buffer.Bytes()); err != nil {
		panic(fmt.Sprintf("Failed to write file:%v|err:%v", outPutFile, err))
	}
}

// 增加头文件
func addHeader(buffer *bytes.Buffer) {
	buffer.WriteString("// Code generated by gg; DO NOT EDIT.\n")
	buffer.WriteString(fmt.Sprintf("// gg version: %s\n", Version))
	buffer.WriteString(fmt.Sprintf("// generate time: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	buffer.WriteString(fmt.Sprintf("// src code version: %s\n", ""))
	buffer.WriteString(fmt.Sprintf("// src code commit time : %s\n", ""))
}
