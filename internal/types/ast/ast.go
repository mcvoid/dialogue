package ast

// top level elements
type (
	Script struct {
		Nodes     []Node
		Functions map[string][]Type
	}
	Symbol string
	Node   struct {
		Name Symbol
		Body []BlockElement
	}
)

// block elements
type (
	BlockElement interface {
		CompareBlock(b BlockElement) bool
	}
	Paragraph []Inline
	Link      struct {
		Dest Symbol
		Text Inline
	}
	Option    []Link
	CodeBlock struct {
		Code []Statement
	}
)

// inlines
type (
	Inline interface {
		CompareInline(b Inline) bool
	}
	Text       string
	InlineCode struct {
		Expr Expression
	}
)

// statements
type (
	Statement interface {
		CompareStatement(b Statement) bool
	}
	StatementBlock []Statement
	Assignment     struct {
		Name Symbol
		Val  Expression
	}
	FunctionCall struct {
		Name   Symbol
		Params []Expression
	}
	GotoNode struct {
		Name Symbol
	}
	Conditional struct {
		Cond       Expression
		Consequent Statement
		Alternate  Statement
	}
	Loop struct {
		Cond       Expression
		Consequent Statement
	}
	InfiniteLoop struct {
		Consequent Statement
	}
)

// expressions
type (
	Expression interface {
		CompareExpression(b Expression) bool
	}
	BinaryOperator string
	BinaryOp       struct {
		Operator BinaryOperator
		LeftArg  Expression
		RightArg Expression
	}
	UnaryOperator string
	UnaryOp       struct {
		Operator UnaryOperator
		Arg      Expression
	}
	Type    string
	Literal struct {
		Type Type
		Val  interface{}
	}
)

// type tags
const (
	StringType  = "string"
	BooleanType = "bool"
	NumberType  = "number"
	SymbolType  = "symbol"
	NullType    = "null"
)

// unary operators
const (
	IncOp UnaryOperator = "inc"
	DecOp UnaryOperator = "dec"
	NotOp UnaryOperator = "not"
	NegOp UnaryOperator = "neg"
)

// binary operators
const (
	AddOp    BinaryOperator = "add"
	SubOp    BinaryOperator = "sub"
	MulOp    BinaryOperator = "mul"
	DivOp    BinaryOperator = "div"
	ModOp    BinaryOperator = "mod"
	GtOp     BinaryOperator = "gt"
	GteOp    BinaryOperator = "gte"
	LtOp     BinaryOperator = "lt"
	LteOp    BinaryOperator = "lte"
	EqOp     BinaryOperator = "eq"
	NeqOp    BinaryOperator = "neq"
	AndOp    BinaryOperator = "and"
	OrOp     BinaryOperator = "or"
	ConcatOp BinaryOperator = "concat"
)

func (s Script) CompareScript(s2 Script) bool {
	if len(s.Nodes) != len(s2.Nodes) {
		return false
	}
	for i := range s.Nodes {
		if !s.Nodes[i].CompareNode(s2.Nodes[i]) {
			return false
		}
	}
	if len(s.Functions) != len(s2.Functions) {
		return false
	}
	for funcName := range s.Functions {
		protoA, protoB := s.Functions[funcName], s2.Functions[funcName]
		if len(protoA) != len(protoB) {
			return false
		}
		for i := range protoA {
			if protoA[i] != protoB[i] {
				return false
			}
		}
	}
	return true
}

func (n Node) CompareNode(n2 Node) bool {
	if len(n.Body) != len(n2.Body) {
		return false
	}
	for i := range n.Body {
		if !n.Body[i].CompareBlock(n2.Body[i]) {
			return false
		}
	}
	return n.Name == n2.Name
}

func (n Paragraph) CompareBlock(b BlockElement) bool {
	s, ok := b.(Paragraph)
	if !ok {
		return false
	}
	if len(n) != len(s) {
		return false
	}
	for i := range n {
		if !n[i].CompareInline(s[i]) {
			return false
		}
	}
	return true
}

func (n Link) CompareBlock(b BlockElement) bool {
	s, ok := b.(Link)
	if !ok {
		return false
	}
	if n.Dest != s.Dest {
		return false
	}
	return n.Text.CompareInline(s.Text)
}

func (n Option) CompareBlock(b BlockElement) bool {
	s, ok := b.(Option)
	if !ok {
		return false
	}
	if len(n) != len(s) {
		return false
	}
	for i := range n {
		if !n[i].CompareBlock(s[i]) {
			return false
		}
	}
	return true
}

func (n CodeBlock) CompareBlock(b BlockElement) bool {
	s, ok := b.(CodeBlock)
	if !ok {
		return false
	}
	if len(n.Code) != len(s.Code) {
		return false
	}
	for i := range n.Code {
		if !n.Code[i].CompareStatement(s.Code[i]) {
			return false
		}
	}
	return true
}

func (n Text) CompareInline(b Inline) bool {
	return n == b
}

func (n InlineCode) CompareInline(b Inline) bool {
	s, ok := b.(InlineCode)
	if !ok {
		return false
	}
	return n.Expr.CompareExpression(s.Expr)
}

func (n StatementBlock) CompareStatement(b Statement) bool {
	s, ok := b.(StatementBlock)
	if !ok {
		return false
	}
	if len(n) != len(s) {
		return false
	}
	for i := range n {
		if !n[i].CompareStatement(s[i]) {
			return false
		}
	}
	return true
}

func (n Assignment) CompareStatement(b Statement) bool {
	s, ok := b.(Assignment)
	if !ok {
		return false
	}
	if !n.Val.CompareExpression(s.Val) {
		return false
	}
	return n.Name == s.Name
}

func (n FunctionCall) CompareStatement(b Statement) bool {
	s, ok := b.(FunctionCall)
	if !ok {
		return false
	}
	if len(n.Params) != len(s.Params) {
		return false
	}
	for i := range n.Params {
		if !n.Params[i].CompareExpression(s.Params[i]) {
			return false
		}
	}
	return n.Name == s.Name
}

func (n GotoNode) CompareStatement(b Statement) bool {
	return n == b
}

func (n Conditional) CompareStatement(b Statement) bool {
	s, ok := b.(Conditional)
	if !ok {
		return false
	}
	if !n.Cond.CompareExpression(s.Cond) {
		return false
	}
	if !n.Consequent.CompareStatement(s.Consequent) {
		return false
	}
	return n.Alternate.CompareStatement(s.Alternate)
}

func (n Loop) CompareStatement(b Statement) bool {
	s, ok := b.(Loop)
	if !ok {
		return false
	}
	if !n.Cond.CompareExpression(s.Cond) {
		return false
	}
	return n.Consequent.CompareStatement(s.Consequent)
}

func (n InfiniteLoop) CompareStatement(b Statement) bool {
	s, ok := b.(InfiniteLoop)
	if !ok {
		return false
	}

	return n.Consequent.CompareStatement(s.Consequent)
}

func (a BinaryOp) CompareExpression(b Expression) bool {
	return a == b
}

func (a UnaryOp) CompareExpression(b Expression) bool {
	return a == b
}

func (a Literal) CompareExpression(b Expression) bool {
	return a == b
}
