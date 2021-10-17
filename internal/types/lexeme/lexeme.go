package lexeme

type ItemType int

const (
	Eof ItemType = iota
	Error
	TextLiteral
	Hash
	ListItemPrefix
	OpenCodeFence
	CloseCodeFence
	OpenInlineCode
	CloseInlineCode
	LineBreak
	OpenCurlyBrace
	CloseCurlyBrace
	OpenSquareBrace
	CloseSquareBrace
	OpenParen
	CloseParen
	Comma
	Semicolon
	Symbol
	Boolean
	Number
	String
	Null
	IfLiteral
	ElseLiteral
	WhileLiteral
	GotoLiteral
	Plus
	Minus
	Star
	Slash
	Percent
	Gt
	Lt
	Gte
	Lte
	Eq
	Neq
	DoubleEq
	Inc
	Dec
	Dot
	And
	Or
	Not
	Type
	ExternKeyword
)

type Item struct {
	Type ItemType
	Val  string
}

func (i Item) CompareItem(b Item) bool {
	return i.Type == b.Type && i.Val == b.Val
}
