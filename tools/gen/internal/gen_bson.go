package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

// generateBson 生成bson的marshal/unmarshal方法
func generateBson(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	genBsonMarshal(file, structTypeExpr, fields)
	genBsonUnmarshal(file, structTypeExpr, fields)
}

// genBsonMarshal bson的marshal
func genBsonMarshal(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	compositeLit := &ast.CompositeLit{
		Type: ast.NewIdent("bson.M"),
		Elts: []ast.Expr{},
	}

	for _, field := range fields {
		fieldName := field.Names[0].Name //前面已检查
		bsonTag, ok := getFieldTag(structTypeExpr, field, "bson:")
		if !ok {
			panic(fmt.Sprintf("类型:%v,字段:%v, 未找到tag.bson", structTypeExpr, fieldName))
		}
		elt := &ast.KeyValueExpr{
			Key:   &ast.BasicLit{Kind: token.STRING, Value: bsonTag},
			Value: &ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent(fieldName)},
		}
		compositeLit.Elts = append(compositeLit.Elts, elt)
	}

	f := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
				},
			},
		},
		Name: ast.NewIdent("MarshalBSON"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("[]byte")},
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names:  []*ast.Ident{ast.NewIdent("doc")},
								Type:   nil,
								Values: []ast.Expr{compositeLit},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent("bson.Marshal"),
							Args: []ast.Expr{ast.NewIdent("doc")},
						},
					},
				},
			},
		},
	}
	//all
	file.Decls = append(file.Decls, f)
}

// genBsonUnmarshal bson的Unmarshal
func genBsonUnmarshal(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	var setList []ast.Stmt
	for _, field := range fields {
		name := field.Names[0].Name //已提前检查
		setList = append(setList, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent("s." + fieldNameToSetter(name)),
				Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("doc"), Sel: ast.NewIdent(fieldNameToBigFiled(name))}},
			},
		})
	}

	f := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
				},
			},
		},
		Name: ast.NewIdent("UnmarshalBSON"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type:  ast.NewIdent("[]byte"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "doc"}},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: buildUnnamedStruct(fields),
						},
					},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun:  ast.NewIdent("bson.Unmarshal"),
							Args: []ast.Expr{ast.NewIdent("data"), &ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("doc")}},
						}},
					},
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"}, // nil值
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "err"},
								},
							},
						},
					},
					Else: nil,
				},
			},
		},
	}
	//setter
	f.Body.List = append(f.Body.List, setList...)
	//return
	f.Body.List = append(f.Body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})
	//all
	file.Decls = append(file.Decls, f)
}
