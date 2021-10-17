package parser

import "github.com/mcvoid/dialogue/internal/types/lexeme"

type (
	memoKey struct {
		name string
		pos  int
	}

	memoVal struct {
		result *Result
		err    error
	}

	Context struct {
		Tokens  []lexeme.Item
		Pos     int
		Grammar map[string]Parselet
		Memo    map[memoKey]memoVal
		Visited map[memoKey]bool
	}
)

func (ctx Context) Move(consumed int) Context {
	return Context{
		Tokens:  ctx.Tokens,
		Pos:     ctx.Pos + consumed,
		Grammar: ctx.Grammar,
		Memo:    ctx.Memo,
		Visited: ctx.Visited,
	}
}

func (ctx Context) Visit(name string) bool {
	key := memoKey{name: name, pos: ctx.Pos}
	visited := ctx.Visited[key]
	ctx.Visited[key] = true
	return visited
}

func (ctx Context) Result(name string, r *Result, err error) {
	key := memoKey{name: name, pos: ctx.Pos}
	val := memoVal{
		result: r,
		err:    err,
	}
	ctx.Memo[key] = val
}

func (ctx Context) GetMemoizedResult(name string) (*Result, bool, error) {
	key := memoKey{name: name, pos: ctx.Pos}
	val, ok := ctx.Memo[key]
	return val.result, ok, val.err
}

func NewContext(tokens []lexeme.Item, grammar map[string]Parselet) Context {
	return Context{
		Tokens:  tokens,
		Pos:     0,
		Grammar: grammar,
		Memo:    map[memoKey]memoVal{},
		Visited: map[memoKey]bool{},
	}
}
