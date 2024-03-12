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

func Gen(dir string) {
	//读取文件
	targetFiles := getTargetFiles(dir)
	if len(targetFiles) == 0 {
		fmt.Printf("dir:%v not found file *.model.go\n", dir)
	}
	//挨个解析文件
	make(targetFiles)
}

func getTargetFiles(dir string) []string {
	var targetFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//只读取当前这层目录
		if info.IsDir() {
			return nil
		}
		//只读取.model.go文件
		if !strings.HasSuffix(info.Name(), ".model.go") {
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

func make(files []string) {
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}
		fset := token.NewFileSet()
		a, err := parser.ParseFile(fset, "", src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, v := range a.Decls {
			genDecl, ok := v.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, sv := range genDecl.Specs {
				tv, ok := sv.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := tv.Type.(*ast.StructType)
				if !ok {
					continue
				}
				//检查不能有可导出的字段，避免污染
				fields := checkAllFieldsNotExport(tv.Name.Name, structType)
				//生成getter和setter
				genGetterSetter(fields)
			}

		}
	}
}

func checkAllFieldsNotExport(structName string, structType *ast.StructType) []*ast.Field {
	var files []*ast.Field
	for _, v := range structType.Fields.List {
		if len(v.Names) == 0 {
			continue
		}
		if len(v.Names) > 1 {
			panic(fmt.Sprintf("names unexpect:%s", v.Names))
		}
		name := v.Names[0]
		if name.IsExported() {
			panic(fmt.Sprintf("field:%v.%v must not export:%v", structName, name))
		}
		files = append(files, v)
	}
	return files
}

func genGetterSetter(fields []*ast.Field) {

}
