package semantic_analysis

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/ast"
)

type bogusAstNode struct{}

func (b bogusAstNode) CompareExpression(n ast.Expression) bool { return false }
func (b bogusAstNode) CompareStatement(n ast.Statement) bool   { return false }
func (b bogusAstNode) CompareBlock(n ast.BlockElement) bool    { return false }

func TestTypeChecking(t *testing.T) {
	for name, test := range map[string]struct {
		input    []ast.Node
		expected EffectiveType
	}{
		"trivial": {
			input:    []ast.Node{},
			expected: Void,
		},
		"valid inline": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.AddOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.SubOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.MulOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.DivOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.ModOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.GtOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.GteOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.LtOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.LteOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.EqOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.NeqOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.NeqOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.ConcatOp,
					LeftArg:  ast.Literal{Type: ast.StringType, Val: "5"},
					RightArg: ast.Literal{Type: ast.NullType, Val: nil},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
				}},
				ast.InlineCode{Expr: ast.BinaryOp{
					Operator: ast.OrOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
				}},
				ast.InlineCode{Expr: ast.UnaryOp{
					Operator: ast.NegOp,
					Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
				}},
				ast.InlineCode{Expr: ast.UnaryOp{
					Operator: ast.IncOp,
					Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
				}},
				ast.InlineCode{Expr: ast.UnaryOp{
					Operator: ast.DecOp,
					Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
				}},
				ast.InlineCode{Expr: ast.UnaryOp{
					Operator: ast.NotOp,
					Arg:      ast.Literal{Type: ast.BooleanType, Val: true},
				}},
				ast.InlineCode{Expr: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.InlineCode{Expr: ast.Literal{Type: ast.NumberType, Val: 5}},
				ast.InlineCode{Expr: ast.Literal{Type: ast.StringType, Val: "5"}},
				ast.InlineCode{Expr: ast.Literal{Type: ast.NullType, Val: nil}},
				ast.InlineCode{Expr: ast.Literal{Type: ast.SymbolType, Val: "nil"}},
			}}}},
			expected: Void,
		},
		"Symbol literal err": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.SymbolType, Val: nil},
			}}}}},
			expected: Error,
		},
		"num literal err": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.NumberType, Val: "5"},
			}}}}},
			expected: Error,
		},
		"bool literal err": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.BooleanType, Val: "5"},
			}}}}},
			expected: Error,
		},
		"string literal err": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.StringType, Val: 5},
			}}}}},
			expected: Error,
		},
		"null literal err": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.NullType, Val: 5},
			}}}}},
			expected: Error,
		},
		"bad literal type": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.Literal{Type: ast.Type("bogus"), Val: 5},
			}}}}},
			expected: Error,
		},
		"not error": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.UnaryOp{
					Operator: ast.NotOp,
					Arg:      ast.Literal{Type: ast.NullType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"neg error": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.UnaryOp{
					Operator: ast.NegOp,
					Arg:      ast.Literal{Type: ast.StringType, Val: "5"},
				},
			}}}}},
			expected: Error,
		},
		"inc error": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.UnaryOp{
					Operator: ast.IncOp,
					Arg:      ast.Literal{Type: ast.StringType, Val: "5"},
				},
			}}}}},
			expected: Error,
		},
		"dec error": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.UnaryOp{
					Operator: ast.DecOp,
					Arg:      ast.Literal{Type: ast.StringType, Val: "5"},
				},
			}}}}},
			expected: Error,
		},
		"bad unary operator": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.UnaryOp{
					Operator: ast.UnaryOperator("bogus"),
					Arg:      ast.Literal{Type: ast.StringType, Val: "5"},
				},
			}}}}},
			expected: Error,
		},
		"bad binary operator": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.BinaryOperator("bogus"),
					LeftArg:  ast.Literal{Type: ast.StringType, Val: "5"},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"add left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.AddOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"sub left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.SubOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"mul left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.MulOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"div left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.DivOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"mod left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.ModOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"gt left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.GtOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"gte left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.GteOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"lt left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.LtOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"lte left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.LteOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"add right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.AddOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"sub right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.SubOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"mul right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.MulOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"div right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.DivOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"mod right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.ModOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"gt right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.GtOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"gte right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.GteOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"lt right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.LtOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"lte right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.LteOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"and left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"or left bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.OrOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			}}}}},
			expected: Error,
		},
		"and right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"or right bool": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: ast.BinaryOp{
					Operator: ast.OrOp,
					LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			}}}}},
			expected: Error,
		},
		"bogus expression type": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Paragraph{ast.InlineCode{
				Expr: bogusAstNode{},
			}}}}},
			expected: Error,
		},
		"bogus statement type": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				bogusAstNode{},
			}}}}},
			expected: Error,
		},
		"bogus block element": {
			input:    []ast.Node{{Name: "abc", Body: []ast.BlockElement{bogusAstNode{}}}},
			expected: Error,
		},
		"link": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Link{
				Dest: "abc",
				Text: ast.Text("def"),
			}}}},
			expected: Void,
		},
		"option": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.Option{
				ast.Link{Dest: "abc", Text: ast.Text("def")},
				ast.Link{Dest: "def", Text: ast.Text("ghi")},
				ast.Link{Dest: "ghi", Text: ast.Text("jkl")},
			}}}},
			expected: Void,
		},
		"good assignment": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: Void,
		},
		"bad assignment": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: bogusAstNode{}},
			}}}}},
			expected: Error,
		},
		"goto": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.GotoNode{Name: "abc"},
			}}}}},
			expected: Void,
		},
		"good function call": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.FunctionCall{
					Name:   "abc",
					Params: []ast.Expression{ast.Literal{Type: ast.BooleanType, Val: true}}},
			}}}}},
			expected: Void,
		},
		"bad function call": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.FunctionCall{
					Name:   "abc",
					Params: []ast.Expression{bogusAstNode{}}},
			}}}}},
			expected: Error,
		},
		"bad function param numbers": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.FunctionCall{
					Name:   "abc",
					Params: []ast.Expression{},
				},
			}}}}},
			expected: Error,
		},
		"bad function param types": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.FunctionCall{
					Name:   "abc",
					Params: []ast.Expression{ast.Literal{Type: ast.StringType, Val: "true"}},
				},
			}}}}},
			expected: Error,
		},
		"good statement Block": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.StatementBlock{
					ast.GotoNode{Name: "abc"},
					ast.GotoNode{Name: "abc"},
					ast.GotoNode{Name: "abc"},
				},
			}}}}},
			expected: Void,
		},
		"bad statement Block": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.StatementBlock{
					ast.GotoNode{Name: "abc"},
					ast.GotoNode{Name: "abc"},
					ast.GotoNode{Name: "abc"},
					bogusAstNode{},
				},
			}}}}},
			expected: Error,
		},
		"good loop": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Loop{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
			expected: Void,
		},
		"loop bad condition": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Loop{
					Cond: ast.Literal{Type: ast.StringType, Val: "true"},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
			expected: Error,
		},
		"loop bad body": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Loop{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						bogusAstNode{},
					},
				},
			}}}}},
			expected: Error,
		},
		"good conditional": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
			expected: Void,
		},
		"conditional bad condition": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Conditional{
					Cond: ast.Literal{Type: ast.StringType, Val: "true"},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
			expected: Error,
		},
		"conditional bad consequent": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						bogusAstNode{},
					},
					Alternate: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
			expected: Error,
		},
		"conditional bad alternate": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						bogusAstNode{},
					},
				},
			}}}}},
			expected: Error,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := TypeCheckScript(ast.Script{
				Functions: map[string][]ast.Type{
					"abc": {ast.BooleanType},
				},
				Nodes: test.input,
			})
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
