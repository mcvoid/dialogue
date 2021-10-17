package ast

import (
	"fmt"
	"testing"
)

func TestExpression(t *testing.T) {
	for i, test := range []struct {
		a        Expression
		b        Expression
		expected bool
	}{
		{
			a:        Literal{Type: NumberType, Val: 5},
			b:        Literal{Type: NumberType, Val: 5},
			expected: true,
		},
		{
			a:        Literal{Type: NumberType, Val: 5},
			b:        Literal{Type: SymbolType, Val: 5},
			expected: false,
		},
		{
			a:        Literal{Type: NumberType, Val: 5},
			b:        Literal{Type: NumberType, Val: "5"},
			expected: false,
		},
		{
			a:        Literal{Type: NumberType, Val: 5},
			b:        UnaryOp{},
			expected: false,
		},
		{
			a: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			b: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			expected: true,
		},
		{
			a: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			b: UnaryOp{
				Operator: IncOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			expected: false,
		},
		{
			a: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			b: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 6},
			},
			expected: false,
		},
		{
			a: UnaryOp{
				Operator: NotOp,
				Arg:      Literal{Type: NumberType, Val: 5},
			},
			b:        Literal{},
			expected: false,
		},
		{
			a: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			b: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			expected: true,
		},
		{
			a: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			b: BinaryOp{
				Operator: MulOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			expected: false,
		},
		{
			a: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			b: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 6},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			expected: false,
		},
		{
			a: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			b: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 6},
			},
			expected: false,
		},
		{
			a: BinaryOp{
				Operator: AddOp,
				LeftArg:  Literal{Type: NumberType, Val: 5},
				RightArg: Literal{Type: NumberType, Val: 5},
			},
			b:        Literal{Type: NumberType, Val: 5},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareExpression(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestStatement(t *testing.T) {
	for i, test := range []struct {
		a        Statement
		b        Statement
		expected bool
	}{
		{
			a: GotoNode{
				Name: "abc",
			},
			b: GotoNode{
				Name: "abc",
			},
			expected: true,
		},
		{
			a: GotoNode{
				Name: "abc",
			},
			b: GotoNode{
				Name: "ac",
			},
			expected: false,
		},
		{
			a: GotoNode{
				Name: "abc",
			},
			b:        StatementBlock{},
			expected: false,
		},
		{
			a: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "abc"},
			},
			b: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "abc"},
			},
			expected: true,
		},
		{
			a: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "abc"},
			},
			b: StatementBlock{
				GotoNode{Name: "abc"},
			},
			expected: false,
		},
		{
			a: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "abc"},
			},
			b: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "def"},
			},
			expected: false,
		},
		{
			a: StatementBlock{
				GotoNode{Name: "abc"},
				GotoNode{Name: "abc"},
			},
			b:        GotoNode{Name: "abc"},
			expected: false,
		},
		{
			a: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			b: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			expected: true,
		},
		{
			a: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			b: Assignment{
				Name: "ab",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			expected: false,
		},
		{
			a: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			b: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 6},
			},
			expected: false,
		},
		{
			a: Assignment{
				Name: "abc",
				Val:  Literal{Type: NumberType, Val: 5},
			},
			b:        GotoNode{},
			expected: false,
		},
		{
			a: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			b: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			expected: true,
		},
		{
			a: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			b: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
				},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			b: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 6},
				},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			b: FunctionCall{
				Name: "ac",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Name: "abc",
				Params: []Expression{
					Literal{Type: NumberType, Val: 5},
					Literal{Type: NumberType, Val: 5},
				},
			},
			b:        GotoNode{},
			expected: false,
		},
		{
			a: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			expected: true,
		},
		{
			a: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b: Loop{
				Cond:       Literal{Type: SymbolType, Val: "ac"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			expected: false,
		},
		{
			a: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{},
			},
			expected: false,
		},
		{
			a: Loop{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b:        GotoNode{Name: "def"},
			expected: false,
		},
		{
			a: InfiniteLoop{
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b: InfiniteLoop{
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			expected: true,
		},
		{
			a: InfiniteLoop{
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b: InfiniteLoop{
				Consequent: StatementBlock{GotoNode{Name: "df"}},
			},
			expected: false,
		},
		{
			a: InfiniteLoop{
				Consequent: StatementBlock{GotoNode{Name: "def"}},
			},
			b:        GotoNode{Name: "def"},
			expected: false,
		},
		{
			a: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			b: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			expected: true,
		},
		{
			a: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			b: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "ab"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			expected: false,
		},
		{
			a: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			b: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			expected: false,
		},
		{
			a: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			b: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{},
			},
			expected: false,
		},
		{
			a: Conditional{
				Cond:       Literal{Type: SymbolType, Val: "abc"},
				Consequent: StatementBlock{GotoNode{Name: "def"}},
				Alternate:  StatementBlock{GotoNode{Name: "def"}},
			},
			b:        GotoNode{Name: "def"},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareStatement(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestInlines(t *testing.T) {
	for i, test := range []struct {
		a        Inline
		b        Inline
		expected bool
	}{
		{
			a:        Text("abc"),
			b:        Text("abc"),
			expected: true,
		},
		{
			a:        Text("abc"),
			b:        Text("bc"),
			expected: false,
		},
		{
			a:        Text("abc"),
			b:        InlineCode{},
			expected: false,
		},
		{
			a: InlineCode{
				Expr: Literal{Type: SymbolType, Val: "abc"},
			},
			b: InlineCode{
				Expr: Literal{Type: SymbolType, Val: "abc"},
			},
			expected: true,
		},
		{
			a: InlineCode{
				Expr: Literal{Type: SymbolType, Val: "abc"},
			},
			b: InlineCode{
				Expr: Literal{Type: SymbolType, Val: "bc"},
			},
			expected: false,
		},
		{
			a: InlineCode{
				Expr: Literal{Type: SymbolType, Val: "abc"},
			},
			b:        Text("abc"),
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareInline(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestBlocks(t *testing.T) {
	for i, test := range []struct {
		a        BlockElement
		b        BlockElement
		expected bool
	}{
		{
			a: Link{
				Dest: "abc",
				Text: Text("def"),
			},
			b: Link{
				Dest: "abc",
				Text: Text("def"),
			},
			expected: true,
		},
		{
			a: Link{
				Dest: "abc",
				Text: Text("def"),
			},
			b: Link{
				Dest: "ab",
				Text: Text("def"),
			},
			expected: false,
		},
		{
			a: Link{
				Dest: "abc",
				Text: Text("def"),
			},
			b: Link{
				Dest: "abc",
				Text: Text("df"),
			},
			expected: false,
		},
		{
			a: Link{
				Dest: "abc",
				Text: Text("def"),
			},
			b:        CodeBlock{},
			expected: false,
		},
		{
			a: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "def"},
				},
			},
			b: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "def"},
				},
			},
			expected: true,
		},
		{
			a: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "def"},
				},
			},
			b: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
				},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "def"},
				},
			},
			b: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "d"},
				},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				Code: []Statement{
					GotoNode{Name: "abc"},
					GotoNode{Name: "def"},
				},
			},
			b:        Link{},
			expected: false,
		},
		{
			a: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("ghi")},
			},
			b: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("ghi")},
			},
			expected: true,
		},
		{
			a: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("ghi")},
			},
			b: Option{
				Link{Dest: "abc", Text: Text("def")},
			},
			expected: false,
		},
		{
			a: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("ghi")},
			},
			b: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("jjj")},
			},
			expected: false,
		},
		{
			a: Option{
				Link{Dest: "abc", Text: Text("def")},
				Link{Dest: "def", Text: Text("ghi")},
			},
			b:        Link{Dest: "def", Text: Text("ghi")},
			expected: false,
		},
		{
			a: Paragraph{
				Text("abc"),
				Text("def"),
			},
			b: Paragraph{
				Text("abc"),
				Text("def"),
			},
			expected: true,
		},
		{
			a: Paragraph{
				Text("abc"),
				Text("def"),
			},
			b: Paragraph{
				Text("abc"),
			},
			expected: false,
		},
		{
			a: Paragraph{
				Text("abc"),
				Text("def"),
			},
			b: Paragraph{
				Text("abc"),
				Text("fff"),
			},
			expected: false,
		},
		{
			a: Paragraph{
				Text("abc"),
				Text("def"),
			},
			b:        Link{},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareBlock(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestNode(t *testing.T) {
	for i, test := range []struct {
		a        Node
		b        Node
		expected bool
	}{
		{
			a: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			b: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			expected: true,
		},
		{
			a: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			b: Node{
				Name: "bc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			b: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("ghi")},
				},
			},
			b: Node{
				Name: "abc",
				Body: []BlockElement{
					Link{Dest: "abc", Text: Text("def")},
					Link{Dest: "def", Text: Text("gi")},
				},
			},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareNode(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestScript(t *testing.T) {
	for i, test := range []struct {
		a        Script
		b        Script
		expected bool
	}{
		{
			a: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: true,
		},
		{
			a: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "ef", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: true,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"b": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"a": {"a", "b"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "b"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
					"b": {"a"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
		{
			a: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			b: Script{
				Functions: map[string][]Type{
					"a": {"a", "b", "c"},
					"b": {"a"},
				},
				Nodes: []Node{
					{Name: "abc", Body: []BlockElement{}},
					{Name: "def", Body: []BlockElement{}},
				},
			},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareScript(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
