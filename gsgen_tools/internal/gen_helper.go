package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

const ignoreNameTag = "`in:\"-\"`" //means ignore name, means inherit
const bsonIgnoreTag = "`bson:\"-\"`"
const maxFieldCount = 64 //脏标记为uint64,可容纳bit 0..63共64个字段

// checkStructField 检查是否合法的filed字段
func checkStructField(ctx *genContext, structNameIdent *ast.Ident, structType *ast.StructType, needDirty, withBson bool, ignoreCheckIdents []string) []*ast.Field {
	//定义检查
	out := checkStructFieldBase(structNameIdent.Name, structType, needDirty, withBson)
	if len(out) > maxFieldCount {
		panic(fmt.Sprintf("类型:%v, 最多只能有%d可导出的字段(因脏标记限制),现在有:%d", structNameIdent, maxFieldCount, len(out)))
	}
	//类型检查
	for _, field := range out {
		checkFiledTypeLegal(ctx, structNameIdent.Name, field.Names[0].Name, field.Type, needDirty, ignoreCheckIdents)
	}

	return out
}

// tagValue 安全获取字段的原始tag字面量(含反引号),字段无tag时返回""
func tagValue(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	return field.Tag.Value
}

// getFieldTag 获取字段指定key的tag值(返回带引号的字符串字面量,如 "_id"),bool表示该tag是否存在且非空
func getFieldTag(field *ast.Field, tagName string) (string, bool) {
	if field.Tag == nil || field.Tag.Value == "" {
		return "", false
	}
	tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
	val, ok := tag.Lookup(tagName)
	if !ok || val == "" {
		return "", false
	}
	return strconv.Quote(val), true
}

// isBsonIgnored 字段是否被标记为bson忽略(bson:"-")
func isBsonIgnored(field *ast.Field) bool {
	if field.Tag == nil || field.Tag.Value == "" {
		return false
	}
	tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
	val, ok := tag.Lookup("bson")
	return ok && val == "-"
}

// bsonFields 过滤掉 bson:"-" 的字段,返回参与bson序列化的字段
func bsonFields(fields []*ast.Field) []*ast.Field {
	out := make([]*ast.Field, 0, len(fields))
	for _, field := range fields {
		if isBsonIgnored(field) {
			continue
		}
		out = append(out, field)
	}
	return out
}

// isBasicType 根据类型名字判断是否基本类型
func isBasicType(typeStr string) bool {
	basicTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"string", "bool",
		"byte", "rune",
	}
	return slices.Contains(basicTypes, typeStr)
}

func isIgnoreSelectorExpr(typ *ast.SelectorExpr, ignoreCheckIdents []string) (bool, string) {
	ident, ok := typ.X.(*ast.Ident)
	if !ok {
		return false, ""
	}
	suffix := "." + typ.Sel.Name
	name := ident.Name + suffix

	for _, s := range ignoreCheckIdents {
		//
		if strings.Contains(s, ".") { //特定包的类型
			if strings.HasSuffix(s, name) {
				return true, strings.TrimSuffix(s, suffix)
			}
		} else { //特定包
			if strings.HasSuffix(s, ident.Name) {
				return true, s
			}
		}

	}

	return false, ""
}

// isBasicType1 根据类型是否基本类型
func isBasicType1(expr ast.Expr) bool {
	if ident, identOk := expr.(*ast.Ident); identOk {
		return isBasicType(ident.Name)
	}
	return false
}

// addImport 增加包名
func addImport(genFile *ast.File, imports ...string) {
	// 添加导入语句
	if len(imports) <= 0 {
		return
	}
	slices.Sort(imports)
	imports = slices.Compact(imports) //去重,避免生成重复import导致编译失败
	importSpecs := make([]ast.Spec, 0, len(imports))
	for _, importPath := range imports {
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(importPath),
			},
		})
	}
	var front = []ast.Decl{&ast.GenDecl{Tok: token.IMPORT, Specs: importSpecs}}
	genFile.Decls = append(front, genFile.Decls...)
}

// firstUpper 字符串首字母大写
func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// fieldNameToSetter 字段名字转Setter方法
func fieldNameToSetter(fieldName string, needSetter bool) string {
	if needSetter {
		return "Set" + firstUpper(fieldName)
	} else {
		return "set" + firstUpper(fieldName)
	}
}

// fieldNameToGetter 字段名字转Getter方法
func fieldNameToGetter(fieldName string) string {
	return "Get" + firstUpper(fieldName)
}

// fieldNameToBigFiled 字段名字转首字母大写
func fieldNameToBigFiled(fieldName string) string {
	return firstUpper(fieldName)
}

// recvS 构造方法的指针接收者列表: (s *T)
func recvS(structTypeExpr *ast.Ident) *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("s")},
				Type:  &ast.StarExpr{X: structTypeExpr},
			},
		},
	}
}

// notNil 构造二元表达式: <ident> != nil
func notNil(ident string) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  &ast.Ident{Name: ident},
		Op: token.NEQ,
		Y:  &ast.Ident{Name: "nil"},
	}
}

// buildUnnamedStruct 根据字段生成匿名struct
func buildUnnamedStruct(fields []*ast.Field) *ast.StructType {
	var structFieldList []*ast.Field
	for _, field := range fields {
		name := field.Names[0].Name //已提前检查
		structFieldList = append(structFieldList, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fieldNameToBigFiled(name))},
			Type:  field.Type,
			Tag:   field.Tag,
		})
	}
	return &ast.StructType{
		Fields: &ast.FieldList{
			List: structFieldList,
		},
	}
}

// buildUnnamedStructWithValue 根据字段生成匿名struct且赋值
func buildUnnamedStructWithValue(fields []*ast.Field) *ast.CompositeLit {
	ret := &ast.CompositeLit{
		Type: buildUnnamedStruct(fields),
	}
	for _, field := range fields {
		ret.Elts = append(ret.Elts, &ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent(field.Names[0].Name)})
	}
	return ret
}

// addHeader 增加头文件
func addHeader(buffer *bytes.Buffer, headAnnotations []string) {
	buffer.WriteString("// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.\n")
	fmt.Fprintf(buffer, "// gen_tools version: %s\n", Version)
	fmt.Fprintf(buffer, "// generate time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range headAnnotations {
		fmt.Fprintf(buffer, "%s\n", v)
	}
}

// printOutFile 输出文件
func printOutFile(fileSet *token.FileSet, astFile *ast.File, outPutFile string, headAnnotations []string) {
	// 打印修改后的源代码
	buffer := bytes.NewBuffer(nil)
	addHeader(buffer, headAnnotations)
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

// 检查基本的规范
func checkStructFieldBase(structName string, structType *ast.StructType, needDirty, withBson bool) []*ast.Field {
	var out []*ast.Field
	dirtyTag := false
	for _, field := range structType.Fields.List {
		if len(field.Names) > 1 { //每行只能定义1个
			panic(fmt.Sprintf("类型:%v,每行只能声明1个filed", structName))
		} else if len(field.Names) == 0 { //匿名类型必须配置ignoreNameTag或gsmodel.DirtyModel类型
			if tagValue(field) == ignoreNameTag {
				continue
			}
			if !needDirty {
				panic(fmt.Sprintf("类型:%v, 非dirty模式下不能有匿名类(dirty模式下则匿名字段必须为gsmodel.DirtyModel类型)", structName))
			}
			//类型检查
			typ := field.Type
			if selectExpr, selectOk := typ.(*ast.SelectorExpr); !selectOk {
				panic(fmt.Errorf("类型:%v, 匿名字段必须为gsmodel.DirtyModel类型 1", structName))
			} else if ident, identOk := selectExpr.X.(*ast.Ident); !identOk || ident.Name != "gsmodel" {
				panic(fmt.Errorf("类型:%v, 匿名字段必须为gsmodel.DirtyModel类型 2", structName))
			} else if selectExpr.Sel.Name != "DirtyModel" {
				panic(fmt.Errorf("类型:%v, 匿名字段必须为gsmodel.DirtyModel类型 3", structName))
			}
			//tag检查:匿名DirtyModel必须被bson忽略(与命名字段共用同一判定)
			if withBson {
				if !isBsonIgnored(field) {
					panic(fmt.Sprintf("类型:%v, 字段:%v, tag必须包含%v 当前是:%v", structName, "DirtyModel", bsonIgnoreTag, tagValue(field)))
				}
			}
			//避免重复定义
			if dirtyTag {
				panic(fmt.Errorf("类型:%v, 匿名字段必须为gsmodel.DirtyModel类型,且只有1个", structName))
			}
			dirtyTag = true
		} else { //只能有1个
			//导出检查-必须小写,不可导出
			fieldName := field.Names[0].Name
			if ast.IsExported(fieldName) {
				panic(fmt.Sprintf("类型:%v,字段:%v, 必须为非导出的(即小写)", structName, fieldName))
			}
			//tag检查
			if withBson && !isBsonIgnored(field) {
				_, found := getFieldTag(field, "bson")
				if !found {
					panic(fmt.Sprintf("类型:%v, 字段:%v, 需要生成bson, 但是缺少tag.bson:, 如需忽略请设置为:%v", structName, field.Names[0].Name, bsonIgnoreTag))
				}
			}
			out = append(out, field)
		}
	}

	if needDirty && !dirtyTag {
		panic(fmt.Sprintf("类型:%v, dirty模式下则匿名字段必须为gsmodel.DirtyModel类型", structName))
	}
	return out
}

// checkFiledTypeLegal 检查类型是否合法
func checkFiledTypeLegal(ctx *genContext, structName string, fieldName string, fieldType ast.Expr, needDirty bool, ignoreCheckIdents []string) {
	switch typ := fieldType.(type) {
	case *ast.Ident:
		if !isBasicType(typ.Name) {
			panic(fmt.Sprintf("类型:%v,字段:%v, 非基本类型只能用指针类型(为了规范),当前为:%v", structName, fieldName, typ.Name))
		}
	case *ast.ArrayType:
		panic(fmt.Sprintf("类型:%v,字段:%v, dirty改为gsmodel.DList,非dirty改为gsmodoel.AList【注:需要指针类型】", structName, fieldName))
	case *ast.MapType:
		panic(fmt.Sprintf("类型:%v,字段:%v, dirty改为gsmodel.DMap,非dirty改为gsmodoel.AMap【注:需要指针类型】", structName, fieldName))
	case *ast.StarExpr: //只能是包内的类型或gsmodel.DirtyModel/AList/DList/AMap/DMap
		checkFieldTypeLegalInStar(ctx, structName, fieldName, typ.X, needDirty, ignoreCheckIdents)
	default:
		panic(fmt.Sprintf("类型:%v,字段:%v, 类型不可用,必须是基本类型或指针类型", structName, fieldName))
	}
}

// checkFieldTypeLegalInStar 检查*指向的类型是否合法
func checkFieldTypeLegalInStar(ctx *genContext, structName string, fieldName string, fieldType ast.Expr, needDirty bool, ignoreCheckIdents []string) {
	switch typ := fieldType.(type) {
	case *ast.Ident:
	case *ast.ArrayType:
		panic(fmt.Sprintf("类型:%v,字段:%v, dirty改为gsmodel.DList,非dirty改为gsmodoel.AList【注:需要指针类型】", structName, fieldName))
	case *ast.MapType:
		panic(fmt.Sprintf("类型:%v,字段:%v, dirty改为gsmodel.DMap,非dirty改为gsmodoel.AMap【注:需要指针类型】", structName, fieldName))
	case *ast.StarExpr: //只能是包内的类型或gsmodel.DirtyModel/AList/DList/AMap/DMap
		panic(fmt.Sprintf("类型:%v,字段:%v,规范,不要使用**双重指针】", structName, fieldName))
	case *ast.IndexExpr:
		genType := mustGSList(structName, fieldName, typ, needDirty)
		checkFiledTypeLegal(ctx, structName, fieldName, genType, needDirty, ignoreCheckIdents)
	case *ast.IndexListExpr:
		genType1, genType2 := mustGSMap(structName, fieldName, typ, needDirty)
		checkFiledTypeLegal(ctx, structName, fieldName, genType1, needDirty, ignoreCheckIdents)
		checkFiledTypeLegal(ctx, structName, fieldName, genType2, needDirty, ignoreCheckIdents)
	case *ast.SelectorExpr:
		b, packageName := isIgnoreSelectorExpr(typ, ignoreCheckIdents)
		if !b {
			panic(fmt.Sprintf("类型:%v,字段:%v,不支持的外部类型", structName, fieldName))
		}
		ctx.markIgnorePackage(packageName)
	default:
		panic(fmt.Sprintf("类型:%v,字段:%v, 类型不可用", structName, fieldName))
	}
}

// mustGSList 必须是gsmodel.AList/DList
func mustGSList(parentName, fieldName string, indexExpr *ast.IndexExpr, needDirty bool) ast.Expr {
	listSelectExpr, listSelectExprOk := indexExpr.X.(*ast.SelectorExpr)
	if !listSelectExprOk {
		panic(fmt.Errorf("类型:%v,字段:%v, 1个泛型参数的的目前强制认为是gsmodel.AList/gsmodel.DList, 转换类型失败 1", parentName, fieldName))
	}
	listSelectExprXIdent, listSelectExprXIdentOk := listSelectExpr.X.(*ast.Ident)
	if !listSelectExprXIdentOk {
		panic(fmt.Errorf("类型:%v,字段:%v, 1个泛型参数的目前强制认为是gsmodel.AList/gsmodel.DList, 转换类型失败 2", parentName, fieldName))
	}
	if listSelectExprXIdent.Name != "gsmodel" {
		panic(fmt.Errorf("类型:%v,字段:%v, 1个泛型参数的目前强制认为是gsmodel.AList/gsmodel.DList, 转换类型失败 pkName:%v", parentName, fieldName, listSelectExprXIdent.Name))
	}
	if needDirty {
		if listSelectExpr.Sel.Name != "DList" {
			panic(fmt.Errorf("类型:%v,字段:%v, 1个泛型参数且dirty模式下强制认为是gsmodel.DList, 当前为:%v", parentName, fieldName, listSelectExpr.Sel.Name))
		}
	} else {
		if listSelectExpr.Sel.Name != "AList" {
			panic(fmt.Errorf("类型:%v,字段:%v, 1个泛型参数且非dirty模式下强制认为是gsmodel.AList,当前为:%v", parentName, fieldName, listSelectExpr.Sel.Name))
		}
	}
	return indexExpr.Index
}

// mustGSMap 必须要是gsmodel.AMap/DMap
func mustGSMap(parentName, fieldName string, indexListExpr *ast.IndexListExpr, needDirty bool) (ast.Expr, ast.Expr) {
	if len(indexListExpr.Indices) != 2 {
		panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数的的目前强制认为是gsmodel.AMap/gsmodel.DMap, 只允许两个泛型参数, 转换类型失败 1", parentName, fieldName))
	}
	listSelectExpr, listSelectExprOk := indexListExpr.X.(*ast.SelectorExpr)
	if !listSelectExprOk {
		panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数的的目前强制认为是gsmodel.AMap/gsmodel.DMap, 转换类型失败 2", parentName, fieldName))
	}
	listSelectExprXIdent, listSelectExprXIdentOk := listSelectExpr.X.(*ast.Ident)
	if !listSelectExprXIdentOk {
		panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数的目前强制认为是gsmodel.AMap/gsmodel.DMap, 转换类型失败 3", parentName, fieldName))
	}
	if listSelectExprXIdent.Name != "gsmodel" {
		panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数的目前强制认为是gsmodel.AMap/gsmodel.DMap, 转换类型失败 pkName:%v", parentName, fieldName, listSelectExprXIdent.Name))
	}
	if needDirty {
		if listSelectExpr.Sel.Name != "DMap" {
			panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数且dirty模式下强制认为是gsmodel.DMap, 当前为:%v", parentName, fieldName, listSelectExpr.Sel.Name))
		}
	} else {
		if listSelectExpr.Sel.Name != "AMap" {
			panic(fmt.Errorf("类型:%v,字段:%v, 多个泛型参数且非dirty模式下强制认为是gsmodel.AMap,当前为:%v", parentName, fieldName, listSelectExpr.Sel.Name))
		}
	}
	return indexListExpr.Indices[0], indexListExpr.Indices[1]
}
