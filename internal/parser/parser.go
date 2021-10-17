package parser

import (
	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

type (
	Parser struct {
		grammar map[string]Parselet
		start   string
	}
	Lexer interface {
		Lex() []lexeme.Item
	}
)

func New() *Parser {
	return &Parser{
		start:   "script",
		grammar: Grammar,
	}
}

func (p *Parser) Parse(l Lexer) (parsetree.Script, error) {
	ctx := NewContext(
		l.Lex(),
		Grammar,
	)
	r, err := p.grammar[p.start](ctx)
	if err != nil {
		return parsetree.Script{}, err
	}
	return r.Val.Script, nil
}
