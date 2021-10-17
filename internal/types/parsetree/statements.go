package parsetree

import "github.com/mcvoid/dialogue/internal/types/lexeme"

type Statement interface {
	CompareStatement(n2 Statement) bool
}

type StatementBlock struct {
	OpenBrace  lexeme.Item
	Statements []Statement
	CloseBrace lexeme.Item
}

func (n StatementBlock) CompareStatement(n2 Statement) bool {
	b, ok := n2.(StatementBlock)
	if !ok {
		return false
	}
	if !n.OpenBrace.CompareItem(b.OpenBrace) {
		return false
	}
	if len(n.Statements) != len(b.Statements) {
		return false
	}
	for i := range n.Statements {
		if !n.Statements[i].CompareStatement(b.Statements[i]) {
			return false
		}
	}
	return n.CloseBrace.CompareItem(b.CloseBrace)
}

type FunctionCall struct {
	Symbol     lexeme.Item
	OpenParen  lexeme.Item
	Args       []Expression
	CloseParen lexeme.Item
	Semicolon  lexeme.Item
}

func (n FunctionCall) CompareStatement(n2 Statement) bool {
	b, ok := n2.(FunctionCall)
	if !ok {
		return false
	}
	if !n.Symbol.CompareItem(b.Symbol) {
		return false
	}
	if !n.OpenParen.CompareItem(b.OpenParen) {
		return false
	}
	if len(n.Args) != len(b.Args) {
		return false
	}
	for i := range n.Args {
		if !n.Args[i].CompareExpression(b.Args[i]) {
			return false
		}
	}
	if !n.CloseParen.CompareItem(b.CloseParen) {
		return false
	}
	return n.Semicolon.CompareItem(b.Semicolon)
}

type Goto struct {
	GotoLiteral lexeme.Item
	Symbol      lexeme.Item
	Semicolon   lexeme.Item
}

func (n Goto) CompareStatement(n2 Statement) bool {
	b, ok := n2.(Goto)
	if !ok {
		return false
	}
	if !n.GotoLiteral.CompareItem(b.GotoLiteral) {
		return false
	}
	if !n.Symbol.CompareItem(b.Symbol) {
		return false
	}
	return n.Semicolon.CompareItem(b.Semicolon)
}

type Conditional struct {
	IfLiteral  lexeme.Item
	Cond       Expression
	Consequent StatementBlock
}

func (n Conditional) CompareStatement(n2 Statement) bool {
	b, ok := n2.(Conditional)
	if !ok {
		return false
	}
	if !n.IfLiteral.CompareItem(b.IfLiteral) {
		return false
	}
	if !n.Cond.CompareExpression(b.Cond) {
		return false
	}
	return n.Consequent.CompareStatement(b.Consequent)
}

type ConditionalWithElse struct {
	IfLiteral   lexeme.Item
	Cond        Expression
	Consequent  StatementBlock
	ElseLiteral lexeme.Item
	Alternate   StatementBlock
}

func (n ConditionalWithElse) CompareStatement(n2 Statement) bool {
	b, ok := n2.(ConditionalWithElse)
	if !ok {
		return false
	}
	if !n.IfLiteral.CompareItem(b.IfLiteral) {
		return false
	}
	if !n.Cond.CompareExpression(b.Cond) {
		return false
	}
	if !n.Consequent.CompareStatement(b.Consequent) {
		return false
	}
	if !n.ElseLiteral.CompareItem(b.ElseLiteral) {
		return false
	}
	return n.Alternate.CompareStatement(b.Alternate)
}

type Loop struct {
	WhileLiteral lexeme.Item
	Cond         Expression
	Body         Statement
}

func (n Loop) CompareStatement(n2 Statement) bool {
	b, ok := n2.(Loop)
	if !ok {
		return false
	}
	if !n.WhileLiteral.CompareItem(b.WhileLiteral) {
		return false
	}
	if !n.Cond.CompareExpression(b.Cond) {
		return false
	}
	return n.Body.CompareStatement(b.Body)
}

type Assignment struct {
	Symbol    lexeme.Item
	EqualSign lexeme.Item
	Value     Expression
	Semicolon lexeme.Item
}

func (n Assignment) CompareStatement(n2 Statement) bool {
	b, ok := n2.(Assignment)
	if !ok {
		return false
	}
	if !n.Symbol.CompareItem(b.Symbol) {
		return false
	}
	if !n.EqualSign.CompareItem(b.EqualSign) {
		return false
	}
	if !n.Value.CompareExpression(b.Value) {
		return false
	}
	return n.Semicolon.CompareItem(b.Semicolon)
}
