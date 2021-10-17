package parser

import (
	"errors"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

func TestExpression(t *testing.T) {
	for name, test := range map[string]struct {
		input    []lexeme.Item
		expected parsetree.Expression
		start    string
		consumed int
		err      error
	}{
		"number": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			start:    "value",
			consumed: 1,
			err:      nil,
		},
		"boolean": {
			input: []lexeme.Item{
				{Type: lexeme.Boolean, Val: "true"},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
			start:    "value",
			consumed: 1,
			err:      nil,
		},
		"symbol": {
			input: []lexeme.Item{
				{Type: lexeme.Symbol, Val: "abc"},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
			start:    "value",
			consumed: 1,
			err:      nil,
		},
		"string": {
			input: []lexeme.Item{
				{Type: lexeme.String, Val: "\"abc\""},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"abc\""}},
			start:    "value",
			consumed: 1,
			err:      nil,
		},
		"null": {
			input: []lexeme.Item{
				{Type: lexeme.Null, Val: "null"},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Null, Val: "null"}},
			start:    "value",
			consumed: 1,
			err:      nil,
		},
		"nested expression": {
			input: []lexeme.Item{
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Null, Val: "null"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Null, Val: "null"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			start:    "value",
			consumed: 3,
			err:      nil,
		},
		"neg": {
			input: []lexeme.Item{
				{Type: lexeme.Minus, Val: "-"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Minus, Val: "-"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "unary",
			consumed: 2,
			err:      nil,
		},
		"inc": {
			input: []lexeme.Item{
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "unary",
			consumed: 2,
			err:      nil,
		},
		"dec": {
			input: []lexeme.Item{
				{Type: lexeme.Dec, Val: "--"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Dec, Val: "--"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "unary",
			consumed: 2,
			err:      nil,
		},
		"not": {
			input: []lexeme.Item{
				{Type: lexeme.Not, Val: "!"},
				{Type: lexeme.Boolean, Val: "true"},
			},
			expected: parsetree.UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not, Val: "!"},
				Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
			},
			start:    "unary",
			consumed: 2,
			err:      nil,
		},
		"unary fallthrough": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			start:    "unary",
			consumed: 1,
			err:      nil,
		},
		"mul": {
			input: []lexeme.Item{
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.UnaryExpression{
					Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
					Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "factor",
			consumed: 4,
			err:      nil,
		},
		"div": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Slash, Val: "/"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Slash, Val: "/"},
				RightOperand: parsetree.UnaryExpression{
					Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
					Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
			},
			start:    "factor",
			consumed: 4,
			err:      nil,
		},
		"mod": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Percent, Val: "%"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Percent, Val: "%"},
				RightOperand: parsetree.UnaryExpression{
					Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
					Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
			},
			start:    "factor",
			consumed: 4,
			err:      nil,
		},
		"add": {
			input: []lexeme.Item{
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "term",
			consumed: 6,
			err:      nil,
		},
		"sub": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Minus, Val: "-"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Minus, Val: "-"},
				RightOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
			},
			start:    "term",
			consumed: 6,
			err:      nil,
		},
		"concat": {
			input: []lexeme.Item{
				{Type: lexeme.String, Val: "\"abc\""},
				{Type: lexeme.Dot, Val: "."},
				{Type: lexeme.String, Val: "\"abc\""},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"abc\""}},
				Operator:     lexeme.Item{Type: lexeme.Dot, Val: "."},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.String, Val: "\"abc\""}},
			},
			start:    "term",
			consumed: 3,
			err:      nil,
		},
		"gt": {
			input: []lexeme.Item{
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Gt, Val: ">"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
				Operator:     lexeme.Item{Type: lexeme.Gt, Val: ">"},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "ineqComparator",
			consumed: 5,
			err:      nil,
		},
		"lt": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Lt, Val: "<"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Lt, Val: "<"},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "ineqComparator",
			consumed: 5,
			err:      nil,
		},
		"gte": {
			input: []lexeme.Item{
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Gte, Val: ">="},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
				Operator:     lexeme.Item{Type: lexeme.Gte, Val: ">="},
				RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
			},
			start:    "ineqComparator",
			consumed: 5,
			err:      nil,
		},
		"lte": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Lte, Val: "<="},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Lte, Val: "<="},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "ineqComparator",
			consumed: 5,
			err:      nil,
		},
		"eq": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.DoubleEq, Val: "=="},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.DoubleEq, Val: "=="},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "eqComparator",
			consumed: 5,
			err:      nil,
		},
		"neq": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Neq, Val: "!="},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Neq, Val: "!="},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "eqComparator",
			consumed: 5,
			err:      nil,
		},
		"and": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.And, Val: "&&"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.And, Val: "&&"},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "andComparator",
			consumed: 5,
			err:      nil,
		},
		"or": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Or, Val: "||"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				Operator:    lexeme.Item{Type: lexeme.Or, Val: "||"},
				RightOperand: parsetree.NestedExpression{
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Expr:       parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				},
			},
			start:    "expression",
			consumed: 5,
			err:      nil,
		},
		"and-or-not precedence": {
			input: []lexeme.Item{
				{Type: lexeme.Boolean, Val: "true"},
				{Type: lexeme.And, Val: "&&"},
				{Type: lexeme.Not, Val: "!"},
				{Type: lexeme.Boolean, Val: "true"},
				{Type: lexeme.Or, Val: "||"},
				{Type: lexeme.Not, Val: "!"},
				{Type: lexeme.Boolean, Val: "true"},
				{Type: lexeme.And, Val: "&&"},
				{Type: lexeme.Boolean, Val: "true"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
					Operator:    lexeme.Item{Type: lexeme.And, Val: "&&"},
					RightOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Not, Val: "!"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
					},
				},
				Operator: lexeme.Item{Type: lexeme.Or, Val: "||"},
				RightOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Not, Val: "!"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
					},
					Operator:     lexeme.Item{Type: lexeme.And, Val: "&&"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Boolean, Val: "true"}},
				},
			},
			start:    "expression",
			consumed: 9,
			err:      nil,
		},
		"arithmetic precedence": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Dec, Val: "--"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Slash, Val: "/"},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					Operator:    lexeme.Item{Type: lexeme.Star, Val: "*"},
					RightOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
				},
				Operator: lexeme.Item{Type: lexeme.Plus, Val: "+"},
				RightOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Dec, Val: "--"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					Operator:     lexeme.Item{Type: lexeme.Slash, Val: "/"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
			},
			start:    "expression",
			consumed: 9,
			err:      nil,
		},
		"gt-lt-eq precedence": {
			input: []lexeme.Item{
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Gt, Val: ">"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.DoubleEq, Val: "=="},
				{Type: lexeme.Dec, Val: "--"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Lte, Val: "<="},
				{Type: lexeme.Number, Val: "5"},
			},
			expected: parsetree.BinaryExpression{
				LeftOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					Operator:    lexeme.Item{Type: lexeme.Gt, Val: ">"},
					RightOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
				},
				Operator: lexeme.Item{Type: lexeme.DoubleEq, Val: "=="},
				RightOperand: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Dec, Val: "--"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					Operator:     lexeme.Item{Type: lexeme.Lte, Val: "<="},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
			},
			start:    "expression",
			consumed: 9,
			err:      nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := Context{
				Tokens:  test.input,
				Pos:     0,
				Grammar: Grammar,
				Memo:    map[memoKey]memoVal{},
				Visited: map[memoKey]bool{},
			}
			parse := Grammar[test.start]
			r, err := parse(ctx)
			if err != nil {
				if test.err == nil {
					t.Errorf("no error expected got %v", err)
				}
				if errors.Is(err, test.err) {
					t.Errorf("expected %v got %v", test.err, err)
				}
				return
			}
			if r.Consumed != test.consumed {
				t.Errorf("expected %d consumed got %d", test.consumed, r.Consumed)
			}
			if !r.Val.Expression.CompareExpression(test.expected) {
				t.Errorf("expected value %v got %v", test.expected, r.Val.Expression)
			}
		})
	}
}

func TestStatement(t *testing.T) {
	for name, test := range map[string]struct {
		input    []lexeme.Item
		expected parsetree.Statement
		start    string
		consumed int
		err      error
	}{
		"goto": {
			input: []lexeme.Item{
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: parsetree.Goto{
				GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
				Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			start:    "goto",
			consumed: 3,
			err:      nil,
		},
		"assignment": {
			input: []lexeme.Item{
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: parsetree.Assignment{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
				Value: parsetree.BinaryExpression{
					LeftOperand: parsetree.UnaryExpression{
						Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
						Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			start:    "assignment",
			consumed: 7,
			err:      nil,
		},
		"function call": {
			input: []lexeme.Item{
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: parsetree.FunctionCall{
				Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args: []parsetree.Expression{
					parsetree.BinaryExpression{
						LeftOperand: parsetree.UnaryExpression{
							Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
							Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
						},
						Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
						RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					},
					parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			start:    "functionCall",
			consumed: 12,
			err:      nil,
		},
		"function call no params": {
			input: []lexeme.Item{
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
			},
			expected: parsetree.FunctionCall{
				Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Args:       []parsetree.Expression{},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
				Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
			},
			start:    "functionCall",
			consumed: 4,
			err:      nil,
		},
		"statement block": {
			input: []lexeme.Item{
				{Type: lexeme.OpenCurlyBrace, Val: "{"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: parsetree.StatementBlock{
				OpenBrace: lexeme.Item{Type: lexeme.OpenCurlyBrace, Val: "{"},
				Statements: []parsetree.Statement{
					parsetree.FunctionCall{
						Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
						Args:       []parsetree.Expression{},
						CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
						Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					parsetree.Assignment{
						Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
						Value: parsetree.BinaryExpression{
							LeftOperand: parsetree.UnaryExpression{
								Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
								Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
							},
							Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
							RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
						},
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
			start:    "statementBlock",
			consumed: 16,
			err:      nil,
		},
		"loop": {
			input: []lexeme.Item{
				{Type: lexeme.WhileLiteral, Val: "while"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Gt, Val: ">"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.OpenCurlyBrace, Val: "{"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: parsetree.Loop{
				WhileLiteral: lexeme.Item{Type: lexeme.WhileLiteral, Val: "while"},
				Cond: parsetree.BinaryExpression{
					LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					Operator:     lexeme.Item{Type: lexeme.Gt, Val: ">"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Body: parsetree.StatementBlock{
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
			start:    "loop",
			consumed: 9,
			err:      nil,
		},
		"cond": {
			input: []lexeme.Item{
				{Type: lexeme.IfLiteral, Val: "if"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Gt, Val: ">"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.OpenCurlyBrace, Val: "{"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: parsetree.Conditional{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond: parsetree.BinaryExpression{
					LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					Operator:     lexeme.Item{Type: lexeme.Gt, Val: ">"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Consequent: parsetree.StatementBlock{
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
			start:    "conditional",
			consumed: 9,
			err:      nil,
		},
		"cond with else": {
			input: []lexeme.Item{
				{Type: lexeme.IfLiteral, Val: "if"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Gt, Val: ">"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.OpenCurlyBrace, Val: "{"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCurlyBrace, Val: "}"},
				{Type: lexeme.ElseLiteral, Val: "else"},
				{Type: lexeme.OpenCurlyBrace, Val: "{"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCurlyBrace, Val: "}"},
			},
			expected: parsetree.ConditionalWithElse{
				IfLiteral: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"},
				Cond: parsetree.BinaryExpression{
					LeftOperand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
					Operator:     lexeme.Item{Type: lexeme.Gt, Val: ">"},
					RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
				},
				Consequent: parsetree.StatementBlock{
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
			start:    "conditionalWithElse",
			consumed: 15,
			err:      nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := Context{
				Tokens:  test.input,
				Pos:     0,
				Grammar: Grammar,
				Memo:    map[memoKey]memoVal{},
				Visited: map[memoKey]bool{},
			}
			parse := Grammar[test.start]
			r, err := parse(ctx)
			if err != nil {
				if test.err == nil {
					t.Errorf("no error expected got %v", err)
				}
				if errors.Is(err, test.err) {
					t.Errorf("expected %v got %v", test.err, err)
				}
				return
			}
			if r.Consumed != test.consumed {
				t.Errorf("expected %d consumed got %d", test.consumed, r.Consumed)
			}
			if !r.Val.Statement.CompareStatement(test.expected) {
				t.Errorf("expected value %v got %v", test.expected, r.Val.Expression)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	for name, test := range map[string]struct {
		input    []lexeme.Item
		expected parsetree.Block
		start    string
		consumed int
		err      error
	}{
		"code block": {
			input: []lexeme.Item{
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: parsetree.CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []parsetree.Statement{
					parsetree.FunctionCall{
						Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
						Args:       []parsetree.Expression{},
						CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
						Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					parsetree.Assignment{
						Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
						Value: parsetree.BinaryExpression{
							LeftOperand: parsetree.UnaryExpression{
								Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
								Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
							},
							Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
							RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
						},
						Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
					parsetree.Goto{
						GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
						Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
						Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
					},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			start:    "codeBlock",
			consumed: 19,
			err:      nil,
		},
		"linkBlock": {
			input: []lexeme.Item{
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: parsetree.LinkBlock{
				Link: parsetree.Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			start:    "linkBlock",
			consumed: 8,
			err:      nil,
		},
		"list": {
			input: []lexeme.Item{
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: parsetree.List{
				Links: []parsetree.ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			start:    "list",
			consumed: 25,
			err:      nil,
		},
		"paragraph": {
			input: []lexeme.Item{
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: parsetree.Paragraph{
				Lines: []parsetree.Line{
					{
						Items: []parsetree.Inline{
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []parsetree.Inline{
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							parsetree.InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code: parsetree.BinaryExpression{
									LeftOperand: parsetree.BinaryExpression{
										LeftOperand: parsetree.UnaryExpression{
											Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
											Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
										},
										Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
										RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
									},
									Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
									RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
								},
								CodeEnd: lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []parsetree.Inline{
							parsetree.InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
							parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			start:    "paragraph",
			consumed: 18,
			err:      nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := Context{
				Tokens:  test.input,
				Pos:     0,
				Grammar: Grammar,
				Memo:    map[memoKey]memoVal{},
				Visited: map[memoKey]bool{},
			}
			parse := Grammar[test.start]
			r, err := parse(ctx)
			if err != nil {
				if test.err == nil {
					t.Errorf("no error expected got %v", err)
				}
				if errors.Is(err, test.err) {
					t.Errorf("expected %v got %v", test.err, err)
				}
				return
			}
			if r.Consumed != test.consumed {
				t.Errorf("expected %d consumed got %d", test.consumed, r.Consumed)
			}
			if !r.Val.Block.CompareBlock(test.expected) {
				t.Errorf("expected value %v got %v", test.expected, r.Val.Expression)
			}
		})
	}
}

func TestNode(t *testing.T) {
	for name, test := range map[string]struct {
		input    []lexeme.Item
		expected parsetree.Node
		consumed int
		err      error
	}{
		"node": {
			input: []lexeme.Item{
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: parsetree.Node{
				Header: parsetree.Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []parsetree.Block{
					parsetree.Paragraph{
						Lines: []parsetree.Line{
							{
								Items: []parsetree.Inline{
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
							{
								Items: []parsetree.Inline{
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
									parsetree.InlineCode{
										CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
										Code: parsetree.BinaryExpression{
											LeftOperand: parsetree.BinaryExpression{
												LeftOperand: parsetree.UnaryExpression{
													Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
													Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
												},
												Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
												RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
											},
											Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
											RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
										},
										CodeEnd: lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
							{
								Items: []parsetree.Inline{
									parsetree.InlineCode{
										CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
										Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
										CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
									},
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					parsetree.List{
						Links: []parsetree.ListItem{
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					parsetree.LinkBlock{
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					parsetree.CodeBlock{
						StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
						StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Code: []parsetree.Statement{
							parsetree.FunctionCall{
								Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
								Args:       []parsetree.Expression{},
								CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
								Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
							parsetree.Assignment{
								Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
								Value: parsetree.BinaryExpression{
									LeftOperand: parsetree.UnaryExpression{
										Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
										Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
									},
									Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
									RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
								},
								Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
							parsetree.Goto{
								GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
								Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
						},
						EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
						CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
			},
			consumed: 74,
			err:      nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := Context{
				Tokens:  test.input,
				Pos:     0,
				Grammar: Grammar,
				Memo:    map[memoKey]memoVal{},
				Visited: map[memoKey]bool{},
			}
			parse := Grammar["node"]
			r, err := parse(ctx)
			if err != nil {
				if test.err == nil {
					t.Errorf("no error expected got %v", err)
				}
				if errors.Is(err, test.err) {
					t.Errorf("expected %v got %v", test.err, err)
				}
				return
			}
			if r.Consumed != test.consumed {
				t.Errorf("expected %d consumed got %d", test.consumed, r.Consumed)
			}
			if !r.Val.Node.CompareNode(test.expected) {
				t.Errorf("expected value %v got %v", test.expected, r.Val.Expression)
			}
		})
	}
}

func TestScript(t *testing.T) {
	for name, test := range map[string]struct {
		input    []lexeme.Item
		expected parsetree.Script
		consumed int
		err      error
	}{
		"node": {
			input: []lexeme.Item{
				{Type: lexeme.ExternKeyword, Val: "extern"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ExternKeyword, Val: "extern"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Type, Val: "bool"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Type, Val: "number"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "ghi"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "jkl"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Eof, Val: ""},
			},
			expected: parsetree.Script{
				FrontMatter: parsetree.FrontMatter{
					FuncDecls: []parsetree.FuncDecl{
						{
							ExternKeyword: lexeme.Item{Type: lexeme.ExternKeyword, Val: "extern"},
							Symbol:        lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen:     lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params:        []lexeme.Item{},
							CloseParen:    lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:     lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:       lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						{
							ExternKeyword: lexeme.Item{Type: lexeme.ExternKeyword, Val: "extern"},
							Symbol:        lexeme.Item{Type: lexeme.Symbol, Val: "def"},
							OpenParen:     lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []parsetree.Node{
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.Paragraph{
								Lines: []parsetree.Line{
									{
										Items: []parsetree.Inline{
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
									{
										Items: []parsetree.Inline{
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
											parsetree.InlineCode{
												CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
												Code: parsetree.BinaryExpression{
													LeftOperand: parsetree.BinaryExpression{
														LeftOperand: parsetree.UnaryExpression{
															Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
															Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
														},
														Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
														RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
													},
													Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
													RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
												},
												CodeEnd: lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
											},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
									{
										Items: []parsetree.Inline{
											parsetree.InlineCode{
												CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
												Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
												CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
											},
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "def"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.List{
								Links: []parsetree.ListItem{
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "ghi"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.LinkBlock{
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "jkl"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.CodeBlock{
								StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
								StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								Code: []parsetree.Statement{
									parsetree.FunctionCall{
										Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
										Args:       []parsetree.Expression{},
										CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
										Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
									parsetree.Assignment{
										Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
										Value: parsetree.BinaryExpression{
											LeftOperand: parsetree.UnaryExpression{
												Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
												Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
											},
											Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
											RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
										},
										Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
									parsetree.Goto{
										GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
										Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
								},
								EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
								CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
				},
			},
			consumed: 104,
			err:      nil,
		},
		"no frontmatter": {
			input: []lexeme.Item{
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Plus, Val: "+"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "def"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "ghi"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "def"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "jkl"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Eq, Val: "="},
				{Type: lexeme.Inc, Val: "++"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Star, Val: "*"},
				{Type: lexeme.Number, Val: "5"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.GotoLiteral, Val: "goto"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Eof, Val: ""},
			},
			expected: parsetree.Script{
				FrontMatter: parsetree.FrontMatter{
					FuncDecls: []parsetree.FuncDecl{},
				},
				Nodes: []parsetree.Node{
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.Paragraph{
								Lines: []parsetree.Line{
									{
										Items: []parsetree.Inline{
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
									{
										Items: []parsetree.Inline{
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
											parsetree.InlineCode{
												CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
												Code: parsetree.BinaryExpression{
													LeftOperand: parsetree.BinaryExpression{
														LeftOperand: parsetree.UnaryExpression{
															Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
															Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
														},
														Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
														RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
													},
													Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
													RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
												},
												CodeEnd: lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
											},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
									{
										Items: []parsetree.Inline{
											parsetree.InlineCode{
												CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
												Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
												CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
											},
											parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
										},
										EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "def"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.List{
								Links: []parsetree.ListItem{
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
									{
										Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
										Link: parsetree.Link{
											OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
											Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
											CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
											OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
											Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
											CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
											EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
										},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "ghi"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.LinkBlock{
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
					{
						Header: parsetree.Header{
							Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
							Name:    lexeme.Item{Type: lexeme.Symbol, Val: "jkl"},
							EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Blocks: []parsetree.Block{
							parsetree.CodeBlock{
								StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
								StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								Code: []parsetree.Statement{
									parsetree.FunctionCall{
										Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
										Args:       []parsetree.Expression{},
										CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
										Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
									parsetree.Assignment{
										Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
										Value: parsetree.BinaryExpression{
											LeftOperand: parsetree.UnaryExpression{
												Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
												Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
											},
											Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
											RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
										},
										Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
									parsetree.Goto{
										GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
										Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
										Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
									},
								},
								EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
								CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
					},
				},
			},
			consumed: 87,
			err:      nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := Context{
				Tokens:  test.input,
				Pos:     0,
				Grammar: Grammar,
				Memo:    map[memoKey]memoVal{},
				Visited: map[memoKey]bool{},
			}
			parse := Grammar["script"]
			r, err := parse(ctx)
			if err != nil {
				if test.err == nil {
					t.Errorf("no error expected got %v", err)
				}
				if errors.Is(err, test.err) {
					t.Errorf("expected %v got %v", test.err, err)
				}
				return
			}
			if r.Consumed != test.consumed {
				t.Errorf("expected %d consumed got %d", test.consumed, r.Consumed)
			}
			if !r.Val.Script.CompareScript(test.expected) {
				t.Errorf("expected value %v got %v", test.expected, r.Val.Expression)
			}
		})
	}
}
