package test

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestGenMarshalBSON(t *testing.T) {
	fset := token.NewFileSet()
	file := &ast.File{
		Name: &ast.Ident{Name: "main"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("c")},
							Type:  &ast.StarExpr{X: ast.NewIdent("TestA")},
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
										Names: []*ast.Ident{ast.NewIdent("doc")},
										Type:  nil,
										Values: []ast.Expr{
											&ast.CompositeLit{
												Type: ast.NewIdent("bson.M"),
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key:   &ast.BasicLit{Kind: token.STRING, Value: `"_id"`},
														Value: &ast.SelectorExpr{X: ast.NewIdent("c"), Sel: ast.NewIdent("id")},
													},
													&ast.KeyValueExpr{
														Key:   &ast.BasicLit{Kind: token.STRING, Value: `"a"`},
														Value: &ast.SelectorExpr{X: ast.NewIdent("c"), Sel: ast.NewIdent("a")},
													},
													&ast.KeyValueExpr{
														Key:   &ast.BasicLit{Kind: token.STRING, Value: `"b"`},
														Value: &ast.SelectorExpr{X: ast.NewIdent("c"), Sel: ast.NewIdent("b")},
													},
												},
											},
										},
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
			},
		},
	}

	if err := printer.Fprint(os.Stdout, fset, file); err != nil {
		panic(err)
	}
}

func TestGenUnmarshalBSON(t *testing.T) {
	fset := token.NewFileSet()
	file := &ast.File{
		Name: &ast.Ident{Name: "main"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
							Type:  &ast.StarExpr{X: ast.NewIdent("TestA")},
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
									Type: &ast.StructType{
										Fields: &ast.FieldList{
											List: []*ast.Field{
												{
													Names: []*ast.Ident{ast.NewIdent("Id")},
													Type:  ast.NewIdent("uint64"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"_id\"`"},
												},
												{
													Names: []*ast.Ident{ast.NewIdent("A")},
													Type:  ast.NewIdent("int64"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"a\"`"},
												},
												{
													Names: []*ast.Ident{ast.NewIdent("B")},
													Type:  ast.NewIdent("int32"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"b\"`"},
												},
											},
										},
									},
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  ast.NewIdent("s.SetId"),
								Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("doc"), Sel: ast.NewIdent("Id")}},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  ast.NewIdent("s.SetA"),
								Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("doc"), Sel: ast.NewIdent("A")}},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  ast.NewIdent("s.SetB"),
								Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("doc"), Sel: ast.NewIdent("B")}},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
	}

	if err := printer.Fprint(os.Stdout, fset, file); err != nil {
		panic(err)
	}
}

func TestGenString(t *testing.T) {
	fset := token.NewFileSet()
	file := &ast.File{
		Name: &ast.Ident{Name: "main"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("s")},
							Type:  &ast.StarExpr{X: ast.NewIdent("TestA")},
						},
					},
				},
				Name: ast.NewIdent("String"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: ast.NewIdent("string")},
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
									Type: &ast.StructType{
										Fields: &ast.FieldList{
											List: []*ast.Field{
												{
													Names: []*ast.Ident{ast.NewIdent("Id")},
													Type:  ast.NewIdent("uint64"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"_id\"`"},
												},
												{
													Names: []*ast.Ident{ast.NewIdent("A")},
													Type:  ast.NewIdent("int64"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"a\"`"},
												},
												{
													Names: []*ast.Ident{ast.NewIdent("B")},
													Type:  ast.NewIdent("int32"),
													Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`bson:\"b\"`"},
												},
											},
										},
									},
									Elts: []ast.Expr{
										&ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent("id")},
										&ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent("a")},
										&ast.SelectorExpr{X: ast.NewIdent("s"), Sel: ast.NewIdent("b")},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun:  ast.NewIdent("fmt.Sprintf"),
									Args: []ast.Expr{ast.NewIdent("\"%v\""), &ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("doc")}},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := printer.Fprint(os.Stdout, fset, file); err != nil {
		panic(err)
	}
}
