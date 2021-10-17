package parsetree

import "github.com/mcvoid/dialogue/internal/types/lexeme"

type Expression interface {
	CompareExpression(n2 Expression) bool
}

type BinaryExpression struct {
	LeftOperand  Expression
	Operator     lexeme.Item
	RightOperand Expression
}

func (n BinaryExpression) CompareExpression(n2 Expression) bool {
	b, ok := n2.(BinaryExpression)
	if !ok {
		return false
	}
	if !n.Operator.CompareItem(b.Operator) {
		return false
	}
	if !n.LeftOperand.CompareExpression(b.LeftOperand) {
		return false
	}
	return n.RightOperand.CompareExpression(b.RightOperand)
}

type UnaryExpression struct {
	Operator lexeme.Item
	Operand  Expression
}

func (n UnaryExpression) CompareExpression(n2 Expression) bool {
	b, ok := n2.(UnaryExpression)
	if !ok {
		return false
	}
	if !n.Operator.CompareItem(b.Operator) {
		return false
	}
	return n.Operand.CompareExpression(b.Operand)
}

type Literal struct {
	Value lexeme.Item
}

func (n Literal) CompareExpression(n2 Expression) bool {
	b, ok := n2.(Literal)
	if !ok {
		return false
	}

	return n.Value.CompareItem(b.Value)
}

type NestedExpression struct {
	OpenParen  lexeme.Item
	Expr       Expression
	CloseParen lexeme.Item
}

func (n NestedExpression) CompareExpression(n2 Expression) bool {
	b, ok := n2.(NestedExpression)
	if !ok {
		return false
	}
	if !n.OpenParen.CompareItem(b.OpenParen) {
		return false
	}
	if !n.Expr.CompareExpression(b.Expr) {
		return false
	}
	return n.CloseParen.CompareItem(b.CloseParen)
}
