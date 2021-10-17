package semantic_analysis

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/ast"
)

func TestDeadCodeElimination(t *testing.T) {
	for name, test := range map[string]struct {
		input    []ast.Node
		expected []ast.Node
	}{
		"trivial": {
			input:    []ast.Node{},
			expected: []ast.Node{},
		},
		"single trivial node": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{}},
			},
		},
		"prune isolated nodes": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{}},
			},
		},
		"keep linked nodes": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Link{Dest: "def", Text: ast.Text("")},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Link{Dest: "def", Text: ast.Text("")},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
			},
		},
		"keep optioned nodes": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Option{
						ast.Link{Dest: "def", Text: ast.Text("")},
						ast.Link{Dest: "ghi", Text: ast.Text("")},
						ast.Link{Dest: "jkl", Text: ast.Text("")},
					},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{}},
				{Name: "jkl", Body: []ast.BlockElement{}},
				{Name: "mno", Body: []ast.BlockElement{}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Option{
						ast.Link{Dest: "def", Text: ast.Text("")},
						ast.Link{Dest: "ghi", Text: ast.Text("")},
						ast.Link{Dest: "jkl", Text: ast.Text("")},
					},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{}},
				{Name: "jkl", Body: []ast.BlockElement{}},
			},
		},
		"transitive nodes": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Link{Dest: "ghi", Text: ast.Text("")},
				}},
				{Name: "def", Body: []ast.BlockElement{
					ast.Link{Dest: "mno", Text: ast.Text("")},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.Link{Dest: "def", Text: ast.Text("")},
				}},
				{Name: "jkl", Body: []ast.BlockElement{}},
				{Name: "mno", Body: []ast.BlockElement{}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Link{Dest: "ghi", Text: ast.Text("")},
				}},
				{Name: "def", Body: []ast.BlockElement{

					ast.Link{Dest: "mno", Text: ast.Text("")},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.Link{Dest: "def", Text: ast.Text("")},
				}},
				{Name: "mno", Body: []ast.BlockElement{}},
			},
		},
		"prune unreachable blocks": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Link{Dest: "def", Text: ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
				}},
				{Name: "def", Body: []ast.BlockElement{
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Link{Dest: "abc", Text: ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
				}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Link{Dest: "def", Text: ast.Text("abc")},
				}},
				{Name: "def", Body: []ast.BlockElement{
					ast.Paragraph{ast.Text("abc")},
					ast.Paragraph{ast.Text("abc")},
					ast.Link{Dest: "abc", Text: ast.Text("abc")},
				}},
			},
		},
		"keep nodes linked by goto": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.StatementBlock{ast.StatementBlock{ast.GotoNode{Name: "ghi"}}},
					}},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.StatementBlock{ast.StatementBlock{ast.GotoNode{Name: "ghi"}}},
					}},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
		},
		"keep nodes linked by reference in loop": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Loop{
							Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
							Consequent: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
						},
					}},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Loop{
							Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
							Consequent: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
						},
					}},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
		},
		"keep nodes linked by reference in conditional consequent": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Conditional{
							Cond: ast.Literal{Type: ast.BooleanType, Val: true},
							Consequent: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
							Alternate: ast.StatementBlock{
								ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
							},
						},
					}},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Conditional{
							Cond: ast.Literal{Type: ast.BooleanType, Val: true},
							Consequent: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
							Alternate: ast.StatementBlock{
								ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
							},
						},
					}},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
		},
		"keep nodes linked by reference in conditional alternate": {
			input: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Conditional{
							Cond: ast.Literal{Type: ast.BooleanType, Val: true},
							Consequent: ast.StatementBlock{
								ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
							},
							Alternate: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
						},
					}},
				}},
				{Name: "def", Body: []ast.BlockElement{}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
			expected: []ast.Node{
				{Name: "abc", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Conditional{
							Cond: ast.Literal{Type: ast.BooleanType, Val: true},
							Consequent: ast.StatementBlock{
								ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
							},
							Alternate: ast.StatementBlock{
								ast.GotoNode{Name: "ghi"},
							},
						},
					}},
				}},
				{Name: "ghi", Body: []ast.BlockElement{
					ast.CodeBlock{Code: []ast.Statement{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					}},
				}},
			},
		},
		"unreachable code in code block": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.GotoNode{Name: "abc"},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.GotoNode{Name: "abc"},
			}}}}},
		},
		"unreachable code in conditional": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
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
		},
		"unreachable code in loop": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
		},
		"unreachable code in infinite loop": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.InfiniteLoop{
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
			}}}}},
		},
		"reachable code in degenerate loop": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.BooleanType, Val: false},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.StatementBlock{},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
		},
		"reachable code in conditional": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Conditional{
					Cond: ast.Literal{Type: ast.BooleanType, Val: true},
					Consequent: ast.StatementBlock{
						ast.GotoNode{Name: "abc"},
					},
					Alternate: ast.StatementBlock{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
		},
		"reachable code in loop": {
			input: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
					Consequent: ast.StatementBlock{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
			expected: []ast.Node{{Name: "abc", Body: []ast.BlockElement{ast.CodeBlock{Code: []ast.Statement{
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Loop{
					Cond: ast.Literal{Type: ast.SymbolType, Val: "var1"},
					Consequent: ast.StatementBlock{
						ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
					},
				},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
				ast.Assignment{Name: "abc", Val: ast.Literal{Type: ast.BooleanType, Val: true}},
			}}}}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := PruneScript(ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes:     test.input,
			})
			expected := ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes:     test.expected,
			}

			if !expected.CompareScript(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
