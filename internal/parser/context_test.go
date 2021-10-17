package parser

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

func TestContext(t *testing.T) {
	tokens := []lexeme.Item{}
	ctx := NewContext(tokens, Grammar)
	if ctx.Grammar == nil || ctx.Tokens == nil || ctx.Memo == nil || ctx.Visited == nil || ctx.Pos != 0 {
		t.Errorf("bad ctx init")
	}
	ctx = ctx.Move(5)
	if ctx.Pos != 5 {
		t.Errorf("ctx move - expected %d got %d", 5, ctx.Pos)
	}

	visited := ctx.Visit("abc")
	if visited != false {
		t.Errorf("ctx visit - should produce false on unvisited rule")
	}
	visited = ctx.Visit("abc")
	if visited != true {
		t.Errorf("ctx visit - should produce true on visited rule")
	}

	r, ok, err := ctx.GetMemoizedResult("abc")
	if ok || r != nil || err != nil {
		t.Errorf("not expected to have value or error when no result is memoized")
	}

	expectedResult := Result{Consumed: 1, Val: Val{Token: lexeme.Item{Type: lexeme.IfLiteral, Val: "if"}}}

	ctx.Result("abc", &expectedResult, nil)
	r, ok, err = ctx.GetMemoizedResult("abc")
	if !ok {
		t.Errorf("value expected after memoization")
	}
	if r.Val.Token != expectedResult.Val.Token {
		t.Errorf("Wrong value")
	}
	if err != nil {
		t.Errorf("expected error")
	}
}
