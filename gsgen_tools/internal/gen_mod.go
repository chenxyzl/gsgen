package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// generateGetters 生成getter
func generateGetters(file *ast.File, structTypeExpr *ast.Ident, field *ast.Field) {
	//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
	fieldName := field.Names[0].Name
	//getter
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Name: ast.NewIdent(fieldNameToGetter(fieldName)),
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
func generateSetters(file *ast.File, structTypeExpr *ast.Ident, field *ast.Field, idx int, exportSetter bool, needDirty bool) {
	//经过检测，要么是基本类型的值类型，要么是struct的指针类型，且名字一定为1
	fieldName := field.Names[0].Name
	isBaseType := isBasicType1(field.Type)
	//setter-body
	var setterBody []ast.Stmt
	//不是基本类型先设置value的Parent
	if needDirty && !isBaseType {
		setterBody = append(setterBody, []ast.Stmt{
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
		}...)
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
		})
	//更新当前dirty
	if needDirty {
		setterBody = append(setterBody, &ast.ExprStmt{ //更新当前的dirty
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("s"),
					Sel: ast.NewIdent("UpdateDirty"),
				},
				Args: []ast.Expr{
					//&ast.BasicLit{
					//	Kind:  token.INT,
					//	Value: strconv.Itoa(idx),
					//},
					&ast.BinaryExpr{
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
			},
		})
	}
	//setter方法体
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Name: ast.NewIdent(fieldNameToSetter(fieldName, exportSetter)),
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

// genString 生成string方法
func genString(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
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
						buildUnnamedStructWithValue(fields),
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
	})
}

// genJsonMarshal 生成json的marshal
func genJsonMarshal(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field) {
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
				},
			},
		},
		Name: ast.NewIdent("MarshalJSON"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("[]byte")},
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
						buildUnnamedStructWithValue(fields),
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent("json.Marshal"),
							Args: []ast.Expr{ast.NewIdent("doc")},
						},
					},
				},
			},
		},
	})
}

// genJsonUnmarshal 生成json的Unmarshal
func genJsonUnmarshal(file *ast.File, structTypeExpr *ast.Ident, fields []*ast.Field, exportSetter bool) {
	var setList []ast.Stmt
	for _, field := range fields {
		name := field.Names[0].Name //已提前检查
		setList = append(setList, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent("s." + fieldNameToSetter(name, exportSetter)),
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
		Name: ast.NewIdent("UnmarshalJSON"),
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
							Fun:  ast.NewIdent("json.Unmarshal"),
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

// genClone 生成Clone方法,Copy一个一样的返回
func genClone(file *ast.File, structTypeExpr *ast.Ident) {
	file.Decls = append(file.Decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: structTypeExpr},
				},
			},
		},
		Name: ast.NewIdent("Clone"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.StarExpr{X: structTypeExpr}},
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "data"}, &ast.Ident{Name: "err"}},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun:  ast.NewIdent("json.Marshal"),
							Args: []ast.Expr{ast.NewIdent("s")},
						},
					},
				},
				&ast.IfStmt{ //field设置自己的dirtyIdx
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "ret"}},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: structTypeExpr,
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("json.Unmarshal"),
							Args: []ast.Expr{
								ast.NewIdent("data"),
								&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("ret")},
							},
						},
					},
				},
				&ast.IfStmt{ //field设置自己的dirtyIdx
					Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "err"},
						Op: token.NEQ,
						Y:  &ast.Ident{Name: "nil"},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{Op: token.AND, X: ast.NewIdent("ret")},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	})
}
