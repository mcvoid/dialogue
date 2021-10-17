package codegen

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/program"
	"github.com/mcvoid/dialogue/internal/types/asm"
	"github.com/mcvoid/dialogue/internal/types/ast"
)

func compareProgram(a, b program.Program) bool {
	if a.Start != b.Start {
		return false
	}
	if len(a.Code) != len(b.Code) {
		return false
	}
	for i := range a.Code {
		i1, i2 := a.Code[i], b.Code[i]
		if i1 != i2 {
			return false
		}
	}

	return true
}

func TestCodegen(t *testing.T) {
	tests := []struct {
		name     string
		ast      []ast.Node
		expected program.Program
		hasError bool
	}{
		{
			name: "empty script",
			ast:  []ast.Node{},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "single empty node",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "two chained node",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.Link{Dest: "Node2", Text: ast.Text("Go to Node2")},
					},
				},
				{
					Name: "Node2",
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Go to Node2"}},
					{Opcode: asm.ShowLine},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "link error",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.Link{Dest: "Node3", Text: ast.Text("Go to Node3")},
					},
				},
			},
			expected: program.Program{},
			hasError: true,
		},
		{
			name: "paragraph",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.Paragraph{
							ast.Text("This is regular"),
							ast.Text("This is emphasized"),
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "This is regular"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "This is emphasized"}},
					{Opcode: asm.Concat},
					{Opcode: asm.ShowLine},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "inline",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.Paragraph{
							ast.Text("This is regular"),
							ast.InlineCode{
								Expr: ast.Literal{Type: ast.SymbolType, Val: ast.Symbol("abc")},
							},
							ast.Text("This is emphasized"),
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "This is regular"}},
					{Opcode: asm.LoadVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "abc"}},
					{Opcode: asm.Concat},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "This is emphasized"}},
					{Opcode: asm.Concat},
					{Opcode: asm.ShowLine},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "option",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.Option{
							{Dest: "Node1", Text: ast.Text("Go to Node1")},
							{Dest: "Node2", Text: ast.Text("Go to Node2")},
						},
					},
				},
				{
					Name: "Node2",
					Body: []ast.BlockElement{},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Go to Node1"}},
					{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 0}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "Go to Node2"}},
					{Opcode: asm.PushChoice, Arg: asm.Value{Type: asm.NumberType, Val: 9}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.ShowChoice},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "assignment",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: ast.StatementBlock{
								ast.Assignment{Name: ast.Symbol("val1"), Val: ast.Literal{Type: ast.NumberType, Val: 3}},
								ast.Assignment{Name: ast.Symbol("val2"), Val: ast.Literal{Type: ast.StringType, Val: "3"}},
								ast.Assignment{Name: ast.Symbol("val3"), Val: ast.Literal{Type: ast.BooleanType, Val: true}},
								ast.Assignment{Name: ast.Symbol("val4"), Val: ast.Literal{Type: ast.NullType, Val: nil}},
								ast.Assignment{Name: ast.Symbol("val5"), Val: ast.Literal{Type: ast.SymbolType, Val: ast.Symbol("val1")}},
							},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val2"}},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val3"}},
					{Opcode: asm.PushNull, Arg: asm.Null},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val4"}},
					{Opcode: asm.LoadVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val5"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "function asm.Call",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{
								ast.FunctionCall{Name: ast.Symbol("func1"), Params: []ast.Expression{
									ast.Literal{Type: ast.NumberType, Val: 3},
									ast.Literal{Type: ast.StringType, Val: "3"},
									ast.Literal{Type: ast.BooleanType, Val: true},
									ast.Literal{Type: ast.NullType, Val: nil},
									ast.Literal{Type: ast.SymbolType, Val: ast.Symbol("val1")},
								}}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.PushNull, Arg: asm.Null},
					{Opcode: asm.LoadVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.Call, Arg: asm.Value{Type: asm.SymbolType, Val: "func1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "goto",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.GotoNode{Name: ast.Symbol("Node2")}},
						},
					},
				},
				{
					Name: "Node2",
					Body: []ast.BlockElement{},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 5}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "multiple goto",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: ast.StatementBlock{
								ast.GotoNode{Name: ast.Symbol("Node2")},
								ast.GotoNode{Name: ast.Symbol("Node2")},
								ast.GotoNode{Name: ast.Symbol("Node2")},
								ast.GotoNode{Name: ast.Symbol("Node2")},
							},
						},
					},
				},
				{
					Name: "Node2",
					Body: []ast.BlockElement{},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "if",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Conditional{
								Cond:       ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
								Consequent: ast.Assignment{Name: ast.Symbol("val1"), Val: ast.Literal{Type: ast.NumberType, Val: 3}},
								Alternate:  ast.Assignment{Name: ast.Symbol("val2"), Val: ast.Literal{Type: ast.StringType, Val: "3"}},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushBool, Arg: asm.False},
					{Opcode: asm.Not, Arg: asm.Value{}},
					{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 9}},
					{Opcode: asm.PushString, Arg: asm.Value{Type: asm.StringType, Val: "3"}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{name: "if",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Conditional{
								Cond:       ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
								Consequent: ast.Assignment{Name: ast.Symbol("val1"), Val: ast.Literal{Type: ast.NumberType, Val: 3}},
								Alternate:  ast.StatementBlock{},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushBool, Arg: asm.False},
					{Opcode: asm.Not, Arg: asm.Value{}},
					{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 6}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "loop",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Loop{
								Cond:       ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
								Consequent: ast.Assignment{Name: ast.Symbol("val1"), Val: ast.Literal{Type: ast.NumberType, Val: 3}},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushBool, Arg: asm.False},
					{Opcode: asm.Not},
					{Opcode: asm.JumpIfFalse, Arg: asm.Value{Type: asm.NumberType, Val: 7}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 1}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "infinite loop",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.InfiniteLoop{
								Consequent: ast.Assignment{Name: ast.Symbol("val1"), Val: ast.Literal{Type: ast.NumberType, Val: 3}},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 1}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "add",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.AddOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Add},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "sub",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.SubOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Subtract},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "mul",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.MulOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Multiply},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "div",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.DivOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Divide},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "mod",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.ModOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Modulo},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "gt",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.GtOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.GreaterThan},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "gte",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.GteOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.GreaterThanOrEqual},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "lt",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.LtOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Lessthan},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "lte",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.LteOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.LessthanOrEqual},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "eq",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.EqOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Equal},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "neq",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.NeqOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.NotEqual},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "and",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.AndOp,
									LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
									RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.And},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "or",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.OrOp,
									LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
									RightArg: ast.Literal{Type: ast.BooleanType, Val: true},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.PushBool, Arg: asm.True},
					{Opcode: asm.Or},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "asm.Concat",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.BinaryOp{
									Operator: ast.ConcatOp,
									LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
									RightArg: ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Concat},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "inc",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.UnaryOp{
									Operator: ast.IncOp,
									Arg:      ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Increment},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "dec",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.UnaryOp{
									Operator: ast.DecOp,
									Arg:      ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Decrement},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "neg",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: []ast.Statement{ast.Assignment{
								Name: ast.Symbol("val1"),
								Val: ast.UnaryOp{
									Operator: ast.NegOp,
									Arg:      ast.Literal{Type: ast.NumberType, Val: 3},
								},
							}},
						},
					},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.PushNumber, Arg: asm.Value{Type: asm.NumberType, Val: 3}},
					{Opcode: asm.Negative},
					{Opcode: asm.StoreVariable, Arg: asm.Value{Type: asm.SymbolType, Val: "val1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
		{
			name: "nested statement block",
			ast: []ast.Node{
				{
					Name: "Node1",
					Body: []ast.BlockElement{
						ast.CodeBlock{
							Code: ast.StatementBlock{
								ast.StatementBlock{
									ast.GotoNode{Name: ast.Symbol("Node2")},
									ast.GotoNode{Name: ast.Symbol("Node2")},
									ast.GotoNode{Name: ast.Symbol("Node2")},
									ast.GotoNode{Name: ast.Symbol("Node2")},
								},
							},
						},
					},
				},
				{
					Name: "Node2",
					Body: []ast.BlockElement{},
				},
			},
			expected: program.Program{
				Start: 0,
				Code: []asm.Instruction{
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.Jump, Arg: asm.Value{Type: asm.NumberType, Val: 11}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node1"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
					{Opcode: asm.EnterNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.ExitNode, Arg: asm.Value{Type: asm.SymbolType, Val: "Node2"}},
					{Opcode: asm.EndDialogue, Arg: asm.Value{}},
				},
			},
			hasError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := Codegen(ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes:     test.ast,
			})
			hasError := err != nil
			if hasError != test.hasError {
				t.Fatalf("Error expected %v got %v", test.hasError, err)
			}
			if !compareProgram(p, test.expected) {
				t.Errorf("Expected %v got %v", test.expected, p)
			}
		})
	}
}
