package semantic_analysis

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/ast"
)

func TestBlockElementFolding(t *testing.T) {
	for name, test := range map[string]struct {
		input    ast.Script
		expected ast.Script
	}{
		"paragraph folding": {
			input: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{
					{
						Name: ast.Symbol("abc"),
						Body: []ast.BlockElement{
							ast.Paragraph{
								ast.Text("abc"), ast.Text(" "), ast.Text("def"),
								ast.InlineCode{
									Expr: ast.BinaryOp{
										Operator: ast.SubOp,
										LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
										RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
									},
								},
								ast.Text("  \t\n"), ast.Text("ghi"),
							},
						},
					},
					{
						Name: ast.Symbol("def"),
						Body: []ast.BlockElement{
							ast.Paragraph{ast.Text("abc"), ast.InlineCode{Expr: ast.Literal{Type: ast.StringType, Val: "def"}}, ast.Text("ghi")},
						},
					},
				},
			},
			expected: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{
					{
						Name: ast.Symbol("abc"),
						Body: []ast.BlockElement{
							ast.Paragraph{
								ast.Text("abc def"),
								ast.InlineCode{
									Expr: ast.BinaryOp{
										Operator: ast.SubOp,
										LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
										RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
									},
								},
								ast.Text(" ghi"),
							},
						},
					},
					{
						Name: ast.Symbol("def"),
						Body: []ast.BlockElement{
							ast.Paragraph{ast.Text("abcdefghi")},
						},
					},
				},
			},
		},
		"Fold links and options": {
			input: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{{Name: ast.Symbol("abc"), Body: []ast.BlockElement{
					ast.Link{Dest: ast.Symbol("def"), Text: ast.Text("ghi")},
					ast.Option{
						ast.Link{Dest: ast.Symbol("abc"), Text: ast.Text("def")},
						ast.Link{Dest: ast.Symbol("def"), Text: ast.Text("ghi")},
						ast.Link{Dest: ast.Symbol("ghi"), Text: ast.Text("jkl")},
					},
				}}}},
			expected: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{{Name: ast.Symbol("abc"), Body: []ast.BlockElement{
					ast.Link{Dest: ast.Symbol("def"), Text: ast.Text("ghi")},
					ast.Option{
						ast.Link{Dest: ast.Symbol("abc"), Text: ast.Text("def")},
						ast.Link{Dest: ast.Symbol("def"), Text: ast.Text("ghi")},
						ast.Link{Dest: ast.Symbol("ghi"), Text: ast.Text("jkl")},
					},
				}}}},
		},
		"fold code blocks": {
			input: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{{Name: ast.Symbol("abc"), Body: []ast.BlockElement{
					ast.CodeBlock{
						Code: []ast.Statement{
							ast.GotoNode{Name: ast.Symbol("abc")},
							ast.Conditional{
								Cond: ast.BinaryOp{
									Operator: ast.AndOp,
									LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
									RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
								},
								Consequent: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
								},
								Alternate: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
								},
							},
							ast.Loop{
								Cond: ast.BinaryOp{
									Operator: ast.AndOp,
									LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
									RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
								},
								Consequent: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
								},
							},
						},
					},
				}}}},
			expected: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{{Name: ast.Symbol("abc"), Body: []ast.BlockElement{
					ast.CodeBlock{
						Code: []ast.Statement{
							ast.GotoNode{Name: ast.Symbol("abc")},
							ast.Conditional{
								Cond: ast.Literal{Type: ast.BooleanType, Val: false},
								Consequent: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.Literal{Type: ast.BooleanType, Val: false},
								},
								Alternate: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.Literal{Type: ast.BooleanType, Val: true},
								},
							},
							ast.Loop{
								Cond: ast.Literal{Type: ast.BooleanType, Val: false},
								Consequent: ast.Assignment{
									Name: ast.Symbol("abc"),
									Val:  ast.Literal{Type: ast.BooleanType, Val: false},
								},
							},
						},
					},
				}}}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := ConstantFoldScript(test.input)
			if !test.expected.CompareScript(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestFoldExpressions(t *testing.T) {
	for name, test := range map[string]struct {
		input    ast.Expression
		expected ast.Expression
	}{
		"add": {
			input: ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 8},
		},
		"x + 0 = x": {
			input: ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 0},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"0 + x = x": {
			input: ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 0},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"x + 1 = ++x": {
			input: ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 1},
			},
			expected: ast.UnaryOp{Operator: ast.IncOp, Arg: ast.Literal{Type: ast.SymbolType, Val: "abc"}},
		},
		"1 + x = ++x": {
			input: ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 1},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.UnaryOp{Operator: ast.IncOp, Arg: ast.Literal{Type: ast.SymbolType, Val: "abc"}},
		},
		"sub": {
			input: ast.BinaryOp{
				Operator: ast.SubOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 2},
		},
		"0 - x = -x": {
			input: ast.BinaryOp{
				Operator: ast.SubOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 0},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.UnaryOp{Operator: ast.NegOp, Arg: ast.Literal{Type: ast.SymbolType, Val: "abc"}},
		},
		"x - 0 = x": {
			input: ast.BinaryOp{
				Operator: ast.SubOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 0},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"x - 1 = --x": {
			input: ast.BinaryOp{
				Operator: ast.SubOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 1},
			},
			expected: ast.UnaryOp{Operator: ast.DecOp, Arg: ast.Literal{Type: ast.SymbolType, Val: "abc"}},
		},
		"mul": {
			input: ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 15},
		},
		"x * 0 = 0": {
			input: ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 0},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 0},
		},
		"0 * x = 0": {
			input: ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 0},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 0},
		},
		"x * 1 = x": {
			input: ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 1},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"1 * x = x": {
			input: ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 1},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"div": {
			input: ast.BinaryOp{
				Operator: ast.DivOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 1},
		},
		"0 / x = 0": {
			input: ast.BinaryOp{
				Operator: ast.DivOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 0},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 0},
		},
		"x / 1 = x": {
			input: ast.BinaryOp{
				Operator: ast.DivOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 1},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"mod": {
			input: ast.BinaryOp{
				Operator: ast.ModOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 5 % 3},
		},
		"gt": {
			input: ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"gte": {
			input: ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"lt": {
			input: ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"lte": {
			input: ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"eq": {
			input: ast.BinaryOp{
				Operator: ast.EqOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"neq": {
			input: ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"concat": {
			input: ast.BinaryOp{
				Operator: ast.ConcatOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 5},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
			},
			expected: ast.Literal{Type: ast.StringType, Val: "53"},
		},
		"and": {
			input: ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"x && T": {
			input: ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"T && x": {
			input: ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"x && F": {
			input: ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"F && x": {
			input: ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: false},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"or": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"x || T": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"T || x": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"x || F": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"F || x": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: false},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"nofold binary": {
			input: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
		},
		"neg": {
			input: ast.UnaryOp{
				Operator: ast.NegOp,
				Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: -5},
		},
		"double neg": {
			input: ast.UnaryOp{
				Operator: ast.NegOp,
				Arg: ast.UnaryOp{
					Operator: ast.NegOp,
					Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 5},
		},
		"double neg non-const": {
			input: ast.UnaryOp{
				Operator: ast.NegOp,
				Arg: ast.UnaryOp{
					Operator: ast.NegOp,
					Arg:      ast.Literal{Type: ast.SymbolType, Val: "abc"},
				},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"inc": {
			input: ast.UnaryOp{
				Operator: ast.IncOp,
				Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 6},
		},
		"dec": {
			input: ast.UnaryOp{
				Operator: ast.DecOp,
				Arg:      ast.Literal{Type: ast.NumberType, Val: 5},
			},
			expected: ast.Literal{Type: ast.NumberType, Val: 4},
		},
		"not": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg:      ast.Literal{Type: ast.BooleanType, Val: true},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"double not": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.UnaryOp{
					Operator: ast.NotOp,
					Arg:      ast.Literal{Type: ast.SymbolType, Val: "abc"},
				},
			},
			expected: ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"not eq => neq": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.BinaryOp{
					Operator: ast.EqOp,
					LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
				},
			},
			expected: ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
			},
		},
		"not gt => lte": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.BinaryOp{
					Operator: ast.GtOp,
					LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
			expected: ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"not lt => gte": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.BinaryOp{
					Operator: ast.LtOp,
					LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
			expected: ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"not gte => lt": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.BinaryOp{
					Operator: ast.GteOp,
					LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
			expected: ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"not lte => gt": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg: ast.BinaryOp{
					Operator: ast.LteOp,
					LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
					RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
			expected: ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"unary nofold": {
			input: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg:      ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.UnaryOp{
				Operator: ast.NotOp,
				Arg:      ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
		},
		"nil": {
			input:    nil,
			expected: nil,
		},
		"x == x => true": {
			input: ast.BinaryOp{
				Operator: ast.EqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"x == y => x == y": {
			input: ast.BinaryOp{
				Operator: ast.EqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.EqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
		"x != x => false": {
			input: ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"x != y => x != y": {
			input: ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
		"x > x => false": {
			input: ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"x > y => x > y": {
			input: ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
		"x < x => false": {
			input: ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"x < y => x < y": {
			input: ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
		"x >= x => true": {
			input: ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"x >= y => x >= y": {
			input: ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
		"x <= x => true": {
			input: ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "abc"},
			},
			expected: ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"x <= y => x <= y": {
			input: ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
			expected: ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.SymbolType, Val: "abc"},

				RightArg: ast.Literal{Type: ast.SymbolType, Val: "def"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual, _ := ConstantFoldExpression(test.input)
			if test.expected != actual {
				t.Errorf("excpected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestConstantFoldStatement(t *testing.T) {
	for name, test := range map[string]struct {
		input    ast.Statement
		expected ast.Statement
	}{
		"trivial": {
			input:    ast.GotoNode{Name: ast.Symbol("abc")},
			expected: ast.GotoNode{Name: ast.Symbol("abc")},
		},
		"assignment": {
			input: ast.Assignment{
				Name: ast.Symbol("abc"),
				Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
			},
			expected: ast.Assignment{
				Name: ast.Symbol("abc"),
				Val:  ast.Literal{Type: ast.BooleanType, Val: false},
			},
		},
		"function call": {
			input: ast.FunctionCall{
				Name: ast.Symbol("abc"),
				Params: []ast.Expression{
					ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
					ast.BinaryOp{
						Operator: ast.AddOp,
						LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
						RightArg: ast.Literal{Type: ast.NumberType, Val: 4},
					},
				},
			},
			expected: ast.FunctionCall{
				Name: ast.Symbol("abc"),
				Params: []ast.Expression{
					ast.Literal{Type: ast.BooleanType, Val: false},
					ast.Literal{Type: ast.NumberType, Val: 7},
				},
			},
		},
		"conditional": {
			input: ast.Conditional{
				Cond: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
				},
				Consequent: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
				},
				Alternate: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
				},
			},
			expected: ast.Conditional{
				Cond: ast.Literal{Type: ast.BooleanType, Val: false},
				Consequent: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.Literal{Type: ast.BooleanType, Val: false},
				},
				Alternate: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.Literal{Type: ast.BooleanType, Val: true},
				},
			},
		},
		"Empty consequent conditional": {
			input: ast.Conditional{
				Cond: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
				},
				Consequent: ast.StatementBlock{},
				Alternate: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
				},
			},
			expected: ast.Conditional{
				Cond: ast.Literal{Type: ast.BooleanType, Val: true},
				Consequent: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.Literal{Type: ast.BooleanType, Val: true},
				},
				Alternate: ast.StatementBlock{},
			},
		},
		"Loop": {
			input: ast.Loop{
				Cond: ast.BinaryOp{
					Operator: ast.AndOp,
					LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
					RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
				},
				Consequent: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
				},
			},
			expected: ast.Loop{
				Cond: ast.Literal{Type: ast.BooleanType, Val: false},
				Consequent: ast.Assignment{
					Name: ast.Symbol("abc"),
					Val:  ast.Literal{Type: ast.BooleanType, Val: false},
				},
			},
		},
		"Statement block": {
			input: ast.StatementBlock{
				ast.GotoNode{Name: ast.Symbol("abc")},
				ast.Conditional{
					Cond: ast.BinaryOp{
						Operator: ast.AndOp,
						LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
						RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
					},
					Consequent: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
					Alternate: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
					},
				},
				ast.Loop{
					Cond: ast.BinaryOp{
						Operator: ast.AndOp,
						LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
						RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
					},
					Consequent: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
				},
			},
			expected: ast.StatementBlock{
				ast.GotoNode{Name: ast.Symbol("abc")},
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: false},
					Consequent: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.Literal{Type: ast.BooleanType, Val: false},
					},
					Alternate: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.Literal{Type: ast.BooleanType, Val: true},
					},
				},
				ast.Loop{
					Cond: ast.Literal{Type: ast.BooleanType, Val: false},
					Consequent: ast.Assignment{
						Name: ast.Symbol("abc"),
						Val:  ast.Literal{Type: ast.BooleanType, Val: false},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := ConstantFoldStatement(test.input)
			if !test.expected.CompareStatement(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
