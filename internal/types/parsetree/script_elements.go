package parsetree

import "github.com/mcvoid/dialogue/internal/types/lexeme"

type (
	Script struct {
		FrontMatter FrontMatter
		Nodes       []Node
		Eof         lexeme.Item
	}
	FrontMatter struct {
		FuncDecls []FuncDecl
		Delimiter lexeme.Item
		EndLine   lexeme.Item
	}
	FuncDecl struct {
		ExternKeyword lexeme.Item
		Symbol        lexeme.Item
		OpenParen     lexeme.Item
		Params        []lexeme.Item
		CloseParen    lexeme.Item
		Semicolon     lexeme.Item
		EndLine       lexeme.Item
	}
	Node struct {
		Header  Header
		EndLine lexeme.Item
		Blocks  []Block
	}
	Header struct {
		Hash    lexeme.Item
		Name    lexeme.Item
		EndLine lexeme.Item
	}
)

type (
	Block interface {
		CompareBlock(n2 Block) bool
	}
	Paragraph struct {
		Lines   []Line
		EndLine lexeme.Item
	}
	LinkBlock struct {
		Link    Link
		EndLine lexeme.Item
	}
	Link struct {
		OpenBrace  lexeme.Item
		Symbol     lexeme.Item
		CloseBrace lexeme.Item
		OpenParen  lexeme.Item
		Text       Inline
		CloseParen lexeme.Item
		EndLine    lexeme.Item
	}
	CodeBlock struct {
		StartFence   lexeme.Item
		StartEndline lexeme.Item
		Code         []Statement
		EndFence     lexeme.Item
		CloseEndline lexeme.Item
		EndLine      lexeme.Item
	}
	List struct {
		Links   []ListItem
		EndLine lexeme.Item
	}
	ListItem struct {
		Prefix lexeme.Item
		Link   Link
	}
	Line struct {
		Items   []Inline
		EndLine lexeme.Item
	}
)

type (
	Inline interface {
		CompareInline(n2 Inline) bool
	}
	Text struct {
		Text lexeme.Item
	}
	InlineCode struct {
		CodeStart lexeme.Item
		Code      Expression
		CodeEnd   lexeme.Item
	}
)

func (n Script) CompareScript(n2 Script) bool {
	if len(n.Nodes) != len(n2.Nodes) {
		return false
	}
	for i := range n.Nodes {
		if !n.Nodes[i].CompareNode(n2.Nodes[i]) {
			return false
		}
	}
	if !n.FrontMatter.CompareFrontMatter(n2.FrontMatter) {
		return false
	}
	return n.Eof.CompareItem(n2.Eof)
}

func (n FrontMatter) CompareFrontMatter(n2 FrontMatter) bool {
	if len(n.FuncDecls) != len(n2.FuncDecls) {
		return false
	}
	for i := range n.FuncDecls {
		if !n.FuncDecls[i].CompareFuncDecl(n2.FuncDecls[i]) {
			return false
		}
	}
	if !n.Delimiter.CompareItem(n2.Delimiter) {
		return false
	}
	return n.EndLine.CompareItem(n2.EndLine)
}

func (n FuncDecl) CompareFuncDecl(n2 FuncDecl) bool {
	if !n.ExternKeyword.CompareItem(n2.ExternKeyword) {
		return false
	}
	if !n.Symbol.CompareItem(n2.Symbol) {
		return false
	}
	if !n.OpenParen.CompareItem(n2.OpenParen) {
		return false
	}
	if len(n.Params) != len(n2.Params) {
		return false
	}
	for i := range n.Params {
		if !n.Params[i].CompareItem(n2.Params[i]) {
			return false
		}
	}
	if !n.CloseParen.CompareItem(n2.CloseParen) {
		return false
	}
	if !n.Semicolon.CompareItem(n2.Semicolon) {
		return false
	}
	return n.EndLine.CompareItem(n2.EndLine)
}

func (n Node) CompareNode(n2 Node) bool {
	if !n.Header.CompareHeader(n2.Header) {
		return false
	}
	if !n.EndLine.CompareItem(n2.EndLine) {
		return false
	}
	if len(n.Blocks) != len(n2.Blocks) {
		return false
	}
	for i := range n.Blocks {
		if !n.Blocks[i].CompareBlock(n2.Blocks[i]) {
			return false
		}
	}
	return true
}

func (n Header) CompareHeader(n2 Header) bool {
	if !n.Hash.CompareItem(n2.Hash) {
		return false
	}
	if !n.Name.CompareItem(n2.Name) {
		return false
	}
	return n.EndLine.CompareItem(n2.EndLine)
}

func (n Paragraph) CompareBlock(n2 Block) bool {
	b, ok := n2.(Paragraph)
	if !ok {
		return false
	}
	if len(n.Lines) != len(b.Lines) {
		return false
	}
	for i := range n.Lines {
		if !n.Lines[i].CompareLine(b.Lines[i]) {
			return false
		}
	}
	return n.EndLine.CompareItem(b.EndLine)
}

func (n List) CompareBlock(n2 Block) bool {
	b, ok := n2.(List)
	if !ok {
		return false
	}
	if len(n.Links) != len(b.Links) {
		return false
	}
	for i := range n.Links {
		if !n.Links[i].CompareListItem(b.Links[i]) {
			return false
		}
	}
	return n.EndLine.CompareItem(b.EndLine)
}

func (n LinkBlock) CompareBlock(n2 Block) bool {
	b, ok := n2.(LinkBlock)
	if !ok {
		return false
	}
	if !n.Link.CompareLink(b.Link) {
		return false
	}
	return n.EndLine.CompareItem(b.EndLine)
}

func (n Link) CompareLink(n2 Link) bool {
	if !n.OpenBrace.CompareItem(n2.OpenBrace) {
		return false
	}
	if !n.Symbol.CompareItem(n2.Symbol) {
		return false
	}
	if !n.CloseBrace.CompareItem(n2.CloseBrace) {
		return false
	}
	if !n.OpenParen.CompareItem(n2.OpenParen) {
		return false
	}
	if !n.Text.CompareInline(n2.Text) {
		return false
	}
	if !n.CloseParen.CompareItem(n2.CloseParen) {
		return false
	}
	return n.EndLine.CompareItem(n2.EndLine)
}

func (n CodeBlock) CompareBlock(n2 Block) bool {
	b, ok := n2.(CodeBlock)
	if !ok {
		return false
	}
	if !n.StartFence.CompareItem(b.StartFence) {
		return false
	}
	if !n.StartEndline.CompareItem(b.StartEndline) {
		return false
	}
	if len(n.Code) != len(b.Code) {
		return false
	}
	for i := range n.Code {
		if !n.Code[i].CompareStatement(b.Code[i]) {
			return false
		}
	}
	if !n.EndFence.CompareItem(b.EndFence) {
		return false
	}
	if !n.CloseEndline.CompareItem(b.CloseEndline) {
		return false
	}
	return n.EndLine.CompareItem(b.EndLine)
}

func (n ListItem) CompareListItem(n2 ListItem) bool {
	if !n.Prefix.CompareItem(n2.Prefix) {
		return false
	}
	return n.Link.CompareLink(n2.Link)
}

func (n Line) CompareLine(n2 Line) bool {
	if len(n.Items) != len(n2.Items) {
		return false
	}
	for i := range n.Items {
		if !n.Items[i].CompareInline(n2.Items[i]) {
			return false
		}
	}
	return n.EndLine.CompareItem(n2.EndLine)
}

func (n Text) CompareInline(n2 Inline) bool {
	b, ok := n2.(Text)
	if !ok {
		return false
	}
	return n.Text.CompareItem(b.Text)
}

func (n InlineCode) CompareInline(n2 Inline) bool {
	b, ok := n2.(InlineCode)
	if !ok {
		return false
	}
	if !n.CodeStart.CompareItem(b.CodeStart) {
		return false
	}
	if !n.Code.CompareExpression(b.Code) {
		return false
	}
	return n.CodeEnd.CompareItem(b.CodeEnd)
}
