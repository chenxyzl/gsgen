package internal

import (
	"fmt"
	"go/ast"
	"go/token"
)

// generateBson 生成bson的marshal/unmarshal方法
func generateBson(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field, needSetter bool) {
	genBsonMarshal(file, structTypeExpr, fields)
	genBsonUnmarshal(file, structTypeExpr, fields, needSetter)
	genBuildDirty(file, structTypeExpr, fields)
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
func genBsonUnmarshal(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field, needSetter bool) {
	var setList []ast.Stmt
	for _, field := range fields {
		name := field.Names[0].Name //已提前检查
		setList = append(setList, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent("s." + fieldNameToSetter(name, needSetter)),
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

// genBuildDirty bson的增量更新
func genBuildDirty(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	var dirtyList []ast.Stmt
	for idx, field := range fields {
		name := field.Names[0].Name //已提前检查
		bsonTag, ok := getFieldTag(structTypeExpr, field, "bson:")
		if !ok {
			panic(fmt.Sprintf("类型:%v,字段:%v, 未找到tag.bson", structTypeExpr, name))
		}
		dirtyBody := &ast.IfStmt{ //field设置自己的dirtyIdx
			Cond: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.Ident{Name: "dirty"},
					Op: token.AND,
					Y: &ast.BinaryExpr{
						X: &ast.BasicLit{
							Kind:  token.INT,
							Value: "1",
						}, // 左操作数
						Op: token.SHL, // 左移操作
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: fmt.Sprintf("%d", idx),
						}, // 右操作数
					},
				},
				Op: token.NEQ,
				Y: &ast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			},
			Body: &ast.BlockStmt{List: []ast.Stmt{}},
			Else: nil,
		}
		if isBasicType1(field.Type) {
			dirtyBody.Body.List = append(dirtyBody.Body.List, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: ast.NewIdent("mdata.AddSetDirtyM"),
					Args: []ast.Expr{
						ast.NewIdent("m"),
						&ast.CallExpr{
							Fun: ast.NewIdent("mdata.MakeBsonKey"),
							Args: []ast.Expr{
								&ast.BasicLit{Kind: token.STRING, Value: bsonTag},
								ast.NewIdent("preKey"),
							},
						},
						ast.NewIdent("s." + name),
					},
				},
			})
		} else {
			dirtyBody.Body.List = append(dirtyBody.Body.List, &ast.IfStmt{ //field设置自己的dirtyIdx
				Cond: &ast.BinaryExpr{
					X:  &ast.Ident{Name: "s." + name},
					Op: token.EQL,
					Y:  &ast.Ident{Name: "nil"},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("mdata.AddUnsetDirtyM"),
								Args: []ast.Expr{
									ast.NewIdent("m"),
									&ast.CallExpr{
										Fun: ast.NewIdent("mdata.MakeBsonKey"),
										Args: []ast.Expr{
											&ast.BasicLit{Kind: token.STRING, Value: bsonTag},
											ast.NewIdent("preKey"),
										},
									},
								},
							},
						},
					},
				},
				Else: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("s." + name + ".BuildDirty"),
								Args: []ast.Expr{
									ast.NewIdent("m"),
									&ast.CallExpr{
										Fun: ast.NewIdent("mdata.MakeBsonKey"),
										Args: []ast.Expr{
											&ast.BasicLit{Kind: token.STRING, Value: bsonTag},
											ast.NewIdent("preKey"),
										},
									},
								},
							},
						},
					},
				},
			})
		}
		dirtyList = append(dirtyList, dirtyBody)
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
		Name: ast.NewIdent("BuildDirty"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("m")},
						Type:  ast.NewIdent("bson.M"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("preKey")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "dirty"}},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent("s.GetDirty"),
							Args: []ast.Expr{},
						},
					},
				},
				&ast.IfStmt{ //field设置自己的dirtyIdx
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "dirty"},
						Op: token.EQL,
						Y:  &ast.Ident{Name: "0"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{},
							},
						},
					},
					Else: nil,
				},
			},
		},
	}
	//setter
	f.Body.List = append(f.Body.List, dirtyList...)
	f.Body.List = append(f.Body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("s"),
				Sel: ast.NewIdent("CleanDirty"),
			},
			Args: []ast.Expr{ast.NewIdent("false")},
		},
	})
	//return
	f.Body.List = append(f.Body.List, &ast.ReturnStmt{
		Results: []ast.Expr{},
	})
	//all
	file.Decls = append(file.Decls, f)
}
