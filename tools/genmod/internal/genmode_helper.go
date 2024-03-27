package internal

import (
	"fmt"
	"go/ast"
	"strings"
)

const bsonIgnoreTag = "`bson:\"-\"`"

func checkStructField(structNameIdent *ast.Ident, structType *ast.StructType, withMongo bool) (out []*ast.Field) {
	contain := false
	for _, field := range structType.Fields.List {
		isDirtyModel, needGenField := isLegalField(structNameIdent, field)
		if isDirtyModel && needGenField {
			panic(fmt.Sprintf("类型:%v, 字段:%v, 内部错误isDirtyModel和needGenField不能同时为true", structNameIdent, field))
		}
		if !isDirtyModel && !needGenField {
			panic(fmt.Sprintf("类型:%v, 字段:%v, 内部错误isDirtyModel和needGenField不能同时为flase", structNameIdent, field))
		}

		//
		if isDirtyModel {
			contain = true
			//生成mongo时候 tag必须设置为bsonIgnoreTag
			if withMongo {
				if field.Tag.Value != bsonIgnoreTag {
					panic(fmt.Sprintf("类型:%v, 字段:%v, 只能是%v 当前是:%v", structNameIdent, "DirtyModel", bsonIgnoreTag, field.Tag.Value))
				}
			}
		}
		//
		if needGenField {
			out = append(out, field)
			//生成mongo tag必须要有bson:
			if withMongo {
				_, found := getFieldTag(structNameIdent, field, "bson:")
				if !found {
					panic(fmt.Sprintf("类型:%v, 字段:%v, 需要生成mongo, 但是缺少tag.bson:, 如需忽略请设置为:%v", structNameIdent, field.Names[0].Name, bsonIgnoreTag))
				}
			}
		}
	}
	if !contain {
		panic(fmt.Sprintf("类型:%v, 必须包含DirtyModel", structNameIdent))
	}
	return out
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

func getFieldTag(structNameIdent *ast.Ident, field *ast.Field, tagName string) (string, bool) {
	if field.Tag == nil || field.Tag.Value == "" {
		return "", false
	}
	tags := strings.Split(strings.Trim(field.Tag.Value, "`"), " ")
	fieldName := field.Names[0].Name
	bsonTag := ""
	for _, tag := range tags {
		if strings.HasPrefix(tag, "bson:") {
			if bsonTag != "" {
				panic(fmt.Sprintf("类型:%v, 字段:%v, tag重复, tag:%v", structNameIdent, fieldName, bsonTag))
			} else {
				bsonTag = strings.TrimPrefix(tag, "bson:")
			}
		}
	}
	return bsonTag, bsonTag != "" && bsonTag != "\"\""
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
