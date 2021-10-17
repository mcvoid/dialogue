package parser

import (
	"errors"
	"fmt"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

var (
	ErrorNoMatches     = errors.New("no matches found")
	ErrorBadRuleName   = errors.New("invalid production name")
	ErrorTokenMismatch = errors.New("invalid token type")
	ErrorLeftRecursion = errors.New("illegal left recursion")
	ErrorUnexpectedEof = errors.New("unexpected EOF")
)

type (
	Val struct {
		Token        lexeme.Item
		Script       parsetree.Script
		FrontMatter  parsetree.FrontMatter
		FuncDecls    []parsetree.FuncDecl
		FuncDecl     parsetree.FuncDecl
		Params       []lexeme.Item
		Nodes        []parsetree.Node
		Node         parsetree.Node
		Header       parsetree.Header
		Link         parsetree.Link
		ListItems    []parsetree.ListItem
		ListItem     parsetree.ListItem
		Lines        []parsetree.Line
		Line         parsetree.Line
		Blocks       []parsetree.Block
		Block        parsetree.Block
		Inlines      []parsetree.Inline
		Inline       parsetree.Inline
		Statements   []parsetree.Statement
		Statement    parsetree.Statement
		Expression   parsetree.Expression
		FuncArgsList []parsetree.Expression
	}
	Result struct {
		Consumed int
		Val      Val
	}
	Action   func(m ...Val) Val
	Rule     func(a Action) Parselet
	Parselet func(ctx Context) (*Result, error)
)

func Seq(p ...Parselet) Rule {
	return func(action Action) Parselet {
		return func(ctx Context) (*Result, error) {
			consumed := 0
			matches := []Val{}
			for _, parse := range p {
				r, err := parse(ctx.Move(consumed))
				if err != nil {
					return nil, err
				}
				matches = append(matches, r.Val)
				consumed += r.Consumed
			}
			return &Result{
				Consumed: consumed,
				Val:      action(matches...),
			}, nil
		}
	}
}

func OneOrMore(parse Parselet) Rule {
	return func(action Action) Parselet {
		return func(ctx Context) (*Result, error) {
			consumed := 0
			matches := []Val{}

			r, err := parse(ctx)
			if err != nil {
				return nil, err
			}
			matches = append(matches, r.Val)
			consumed += r.Consumed

			for {
				r, err = parse(ctx.Move(consumed))
				if err != nil {
					break
				}

				matches = append(matches, r.Val)
				consumed += r.Consumed
			}

			return &Result{
				Consumed: consumed,
				Val:      action(matches...),
			}, nil
		}
	}
}

func ZeroOrMore(parse Parselet) Rule {
	return func(action Action) Parselet {
		return func(ctx Context) (*Result, error) {
			consumed := 0
			matches := []Val{}
			for {
				r, err := parse(ctx.Move(consumed))
				if err != nil {
					break
				}

				matches = append(matches, r.Val)
				consumed += r.Consumed
			}

			return &Result{
				Consumed: consumed,
				Val:      action(matches...),
			}, nil
		}
	}
}

func Empty(action Action) Parselet {
	return func(ctx Context) (*Result, error) {
		return &Result{
			Consumed: 0,
			Val:      action(),
		}, nil
	}
}

func Or(p ...Parselet) Parselet {
	return func(ctx Context) (*Result, error) {
		for _, parse := range p {
			r, err := parse(ctx)
			if err != nil {
				continue
			}
			return r, nil
		}
		if ctx.Pos >= len(ctx.Tokens) {
			return nil, ErrorUnexpectedEof
		}
		return nil, fmt.Errorf("%w at token: %v", ErrorNoMatches, ctx.Tokens[ctx.Pos])
	}
}

func Nonterm(name string) Parselet {
	return func(ctx Context) (*Result, error) {
		if r, ok, err := ctx.GetMemoizedResult(name); ok {
			return r, err
		}
		if ctx.Visit(name) {
			return nil, fmt.Errorf("%w, rule: %s", ErrorLeftRecursion, name)
		}
		parse, ok := ctx.Grammar[name]
		if !ok {
			return nil, fmt.Errorf("%w name: %s", ErrorBadRuleName, name)
		}
		r, err := parse(ctx)
		ctx.Result(name, r, err)
		return r, err
	}
}

func Term(name lexeme.ItemType) Parselet {
	return func(ctx Context) (*Result, error) {
		if ctx.Pos >= len(ctx.Tokens) {
			return nil, ErrorUnexpectedEof
		}
		token := ctx.Tokens[ctx.Pos]
		if token.Type != name {
			return nil, fmt.Errorf("%w: token type %d expected %d", ErrorTokenMismatch, token.Type, name)
		}
		return &Result{
			Consumed: 1,
			Val:      Val{Token: token},
		}, nil
	}
}
