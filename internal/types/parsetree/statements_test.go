package parsetree

import (
	"fmt"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

func TestStatements(t *testing.T) {
	for i, test := range []struct {
		a        Statement
		b        Statement
		expected bool
	}{
		{
			a: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: true,
		},
		{
			a: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "gato"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "ac"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ":"},
			},
			expected: false,
		},
		{
			a: Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b:        StatementBlock{},
			expected: false,
		},
		{
			a: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: true,
		},
		{
			a: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b:        Goto{},
			expected: false,
		},
		{
			a: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "}"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: false,
		},
		{
			a: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "{"},
			},
			expected: false,
		},
		{
			a: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: true,
		},
		{
			a: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: false,
		},
		{
			a: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			b: StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []Statement{
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "def"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: true,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abb"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: ")"},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: "("},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ":"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "ddd"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []Expression{
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
					Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b:        Goto{},
			expected: false,
		},
		{
			a: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: true,
		},
		{
			a: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "iff"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "aaa"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "}"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b:        Goto{},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: true,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "iff"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "aaa"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "}"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "elle"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "{"},
				},
			},
			expected: false,
		},
		{
			a: ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Consequent: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
				ElseLiteral: lexeme.Item{Type: lexeme.ElseLiteral, Val: "else"},
				Alternate: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b:        Goto{},
			expected: false,
		},
		{
			a: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: true,
		},
		{
			a: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "whale"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "ac"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			expected: false,
		},
		{
			a: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "{"},
				},
			},
			expected: false,
		},
		{
			a: Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond:         Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				Body: StatementBlock{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
					Statements: []Statement{},
					CloseBrace: lexeme.Item{Type: lexeme.CloseCurlyBrace, Val: "}"},
				},
			},
			b:        Goto{},
			expected: false,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: true,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "ab"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "=="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "df"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: false,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ":"},
			},
			expected: false,
		},
		{
			a: Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value:     Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			b:        Goto{},
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
