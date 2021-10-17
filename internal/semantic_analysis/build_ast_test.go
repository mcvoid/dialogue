package semantic_analysis

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/ast"
	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

type bogusParseTree struct{}

func (b bogusParseTree) CompareStatement(n2 parsetree.Statement) bool {
	return false
}

func TestBuildExpressionAst(t *testing.T) {
	for name, test := range map[string]struct {
		input    parsetree.Expression
		expected ast.Expression
	}{
		"null": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.Null, Val: "null"}},
			ast.Literal{Type: ast.NullType, Val: nil},
		},
		"symbol": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
			ast.Literal{Type: ast.SymbolType, Val: "abc"},
		},
		"string": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"abc\""}},
			ast.Literal{Type: ast.StringType, Val: "abc"},
		},
		"number": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			ast.Literal{Type: ast.NumberType, Val: 5},
		},
		"boolean true": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
			ast.Literal{Type: ast.BooleanType, Val: true},
		},
		"boolean false": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
			ast.Literal{Type: ast.BooleanType, Val: false},
		},
		"bogus literal": {
			parsetree.Literal{Value: lexeme.Item{Type: lexeme.And, Val: "false"}},
			ast.Literal{},
		},
		"nested literal": {
			parsetree.NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			ast.Literal{Type: ast.NumberType, Val: 5},
		},
		"bogus nested expression": {
			parsetree.NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       nil,
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			nil,
		},
		"not": {
			parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not, Val: "!"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
			},
			ast.UnaryOp{Operator: ast.NotOp, Arg: ast.Literal{Type: ast.BooleanType, Val: false}},
		},
		"neg": {
			parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Minus, Val: "-"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.UnaryOp{Operator: ast.NegOp, Arg: ast.Literal{Type: ast.NumberType, Val: 5}},
		},
		"inc": {
			parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.UnaryOp{Operator: ast.IncOp, Arg: ast.Literal{Type: ast.NumberType, Val: 5}},
		},
		"dec": {
			parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Dec, Val: "++"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.UnaryOp{Operator: ast.DecOp, Arg: ast.Literal{Type: ast.NumberType, Val: 5}},
		},
		"add": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.AddOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"sub": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Minus, Val: "-"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.SubOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"mul": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.MulOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"div": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Slash, Val: "/"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.DivOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"mod": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Percent, Val: "%"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.ModOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"equal": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.DoubleEq, Val: "=="},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.EqOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"not equal": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Neq, Val: "!="},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.NeqOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"gt": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Gt, Val: ">"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.GtOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"lt": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Lt, Val: "<"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.LtOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"gte": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Gte, Val: ">="},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.GteOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"lte": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Lte, Val: "<="},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "3"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			ast.BinaryOp{
				Operator: ast.LteOp,
				LeftArg:  ast.Literal{Type: ast.NumberType, Val: 3},
				RightArg: ast.Literal{Type: ast.NumberType, Val: 5},
			},
		},
		"and": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.And, Val: "&&"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
			},
			ast.BinaryOp{
				Operator: ast.AndOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
		},
		"or": {
			parsetree.BinaryExpression{
				Operator:     lexeme.Item{Type: lexeme.Or, Val: "||"},
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
			},
			ast.BinaryOp{
				Operator: ast.OrOp,
				LeftArg:  ast.Literal{Type: ast.BooleanType, Val: true},
				RightArg: ast.Literal{Type: ast.BooleanType, Val: false},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := BuildExpressionAst(test.input)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestBuildStatement(t *testing.T) {
	for name, test := range map[string]struct {
		input    parsetree.Statement
		expected ast.Statement
	}{
		"goto": {
			parsetree.Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			ast.GotoNode{Name: ast.Symbol("abc")},
		},
		"assignment": {
			parsetree.Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			ast.Assignment{Name: ast.Symbol("abc"), Val: ast.Literal{Type: ast.BooleanType, Val: false}},
		},
		"conditional": {
			parsetree.Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
				Consequent: parsetree.StatementBlock{
					OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []parsetree.Statement{
						parsetree.Assignment{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
							Value:     parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
							Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
						parsetree.Goto{
							GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
							Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
					},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			ast.Conditional{
				Cond: ast.Literal{Type: ast.BooleanType, Val: false},
				Consequent: ast.StatementBlock{
					ast.Assignment{Name: ast.Symbol("abc"), Val: ast.Literal{Type: ast.BooleanType, Val: false}},
					ast.GotoNode{Name: ast.Symbol("abc")},
				},
				Alternate: ast.StatementBlock{},
			},
		},
		"conditional with else": {
			parsetree.ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
				Consequent: parsetree.StatementBlock{
					OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []parsetree.Statement{
						parsetree.Assignment{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
							Value:     parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
							Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
					},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: parsetree.StatementBlock{
					OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []parsetree.Statement{
						parsetree.Goto{
							GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
							Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
					},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			ast.Conditional{
				Cond: ast.Literal{Type: ast.BooleanType, Val: false},
				Consequent: ast.StatementBlock{
					ast.Assignment{Name: ast.Symbol("abc"), Val: ast.Literal{Type: ast.BooleanType, Val: false}},
				},
				Alternate: ast.StatementBlock{
					ast.GotoNode{Name: ast.Symbol("abc")},
				},
			},
		},
		"loop": {
			parsetree.Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
				Body: parsetree.StatementBlock{
					OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []parsetree.Statement{
						parsetree.Assignment{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
							Value:     parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
							Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
						parsetree.Goto{
							GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
							Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
						},
					},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			ast.Loop{
				Cond: ast.Literal{Type: ast.BooleanType, Val: false},
				Consequent: ast.StatementBlock{
					ast.Assignment{Name: ast.Symbol("abc"), Val: ast.Literal{Type: ast.BooleanType, Val: false}},
					ast.GotoNode{Name: ast.Symbol("abc")},
				},
			},
		},
		"function call": {
			parsetree.FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []parsetree.Expression{
					parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"abc\""}},
					parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			ast.FunctionCall{
				Name: ast.Symbol("abc"),
				Params: []ast.Expression{
					ast.Literal{Type: ast.StringType, Val: "abc"},
					ast.Literal{Type: ast.NumberType, Val: 5},
				},
			},
		},
		"bogus node": {
			bogusParseTree{},
			nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := BuildStatementAst(test.input)
			if test.expected == nil {
				if test.expected != actual {
					t.Errorf("expected %v got %v", test.expected, actual)
				}
				return
			}
			if !test.expected.CompareStatement(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestBlockElements(t *testing.T) {
	for name, test := range map[string]struct {
		input    parsetree.Block
		expected ast.BlockElement
	}{
		"bogus block": {
			input:    nil,
			expected: nil,
		},
		"list": {
			input: parsetree.List{
				Links: []parsetree.ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
						},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "def"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
						},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: ast.Option{
				ast.Link{Dest: "abc", Text: ast.Text("abc")},
				ast.Link{Dest: "def", Text: ast.Text("def")},
			},
		},
		"link": {
			input: parsetree.LinkBlock{
				Link: parsetree.Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: ast.Link{Dest: "abc", Text: ast.Text("abc")},
		},
		"Code block": {
			input: parsetree.CodeBlock{
				StartFence: lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				Code: []parsetree.Statement{
					parsetree.Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					parsetree.Assignment{
						Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
						Value:     parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "false"}},
						Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				EndFence: lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				EndLine:  lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: ast.CodeBlock{
				Code: []ast.Statement{
					ast.GotoNode{Name: ast.Symbol("abc")},
					ast.Assignment{Name: ast.Symbol("abc"), Val: ast.Literal{Type: ast.BooleanType, Val: false}},
				},
			},
		},
		"Paragraph": {
			input: parsetree.Paragraph{
				Lines: []parsetree.Line{
					{
						Items: []parsetree.Inline{
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []parsetree.Inline{
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							parsetree.InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode},
								Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"ghi\""}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: ast.Paragraph{
				ast.Text("abc"),
				ast.Text("def"),
				ast.Text("\n"),
				ast.Text("abc"),
				ast.Text("def"),
				ast.InlineCode{Expr: ast.Literal{Type: ast.StringType, Val: "ghi"}},
				ast.Text("\n"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := BuildBlockAst(test.input)
			if test.expected == nil {
				if actual != test.expected {
					t.Errorf("expected %v got %v", test.expected, actual)
				}
				return
			}
			if !test.expected.CompareBlock(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestBuildScript(t *testing.T) {
	for name, test := range map[string]struct {
		input    parsetree.Script
		expected ast.Script
	}{
		"trivial": {
			input:    parsetree.Script{},
			expected: ast.Script{},
		},
		"one node one block": {
			input: parsetree.Script{
				Nodes: []parsetree.Node{
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						Blocks: []parsetree.Block{
							parsetree.LinkBlock{
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: ast.Script{
				Functions: map[string][]ast.Type{},
				Nodes: []ast.Node{
					{
						Name: ast.Symbol("abc"),
						Body: []ast.BlockElement{
							ast.Link{Dest: ast.Symbol("abc"), Text: ast.Text("abc")},
						},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := BuildScriptAst(test.input)
			if !test.expected.CompareScript(actual) {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
