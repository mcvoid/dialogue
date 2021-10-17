package parsetree

import (
	"fmt"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

func TestExpressions(t *testing.T) {
	for i, test := range []struct {
		a        Expression
		b        Expression
		expected bool
	}{
		{
			a:        Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			b:        Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			expected: true,
		},
		{
			a:        Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			b:        Literal{lexeme.Item{Type: lexeme.Symbol, Val: "b"}},
			expected: false,
		},
		{
			a:        Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			b:        Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			expected: false,
		},
		{
			a: Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			b: NestedExpression{
				Expr: Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
			},
			expected: false,
		},
		{
			a: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			b: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: true,
		},
		{
			a: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			b: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.TextLiteral, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: false,
		},
		{
			a: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			b: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.String, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			expected: false,
		},
		{
			a: NestedExpression{
				OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
				Expr:       Literal{lexeme.Item{Type: lexeme.Symbol, Val: "a"}},
				CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
			},
			b:        Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			expected: false,
		},
		{
			a: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			expected: true,
		},
		{
			a: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Inc},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			expected: false,
		},
		{
			a: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "b"}},
			},
			expected: false,
		},
		{
			a: UnaryExpression{
				Operator: lexeme.Item{Type: lexeme.Not},
				Operand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b:        Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			expected: false,
		},
		{
			a: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			expected: true,
		},
		{
			a: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "b"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			expected: false,
		},
		{
			a: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.Minus},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			expected: false,
		},
		{
			a: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "b"}},
			},
			expected: false,
		},
		{
			a: BinaryExpression{
				LeftOperand:  Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
				Operator:     lexeme.Item{Type: lexeme.And},
				RightOperand: Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
			},
			b:        Literal{lexeme.Item{Type: lexeme.Number, Val: "a"}},
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
