package parser

import (
	"errors"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

func times(n int, v Val, err error) Parselet {
	return func(ctx Context) (*Result, error) {
		if n > 0 {
			n--
			return &Result{
				Consumed: 1,
				Val:      v,
			}, nil
		}
		return nil, err
	}
}

func val(v Val) Parselet {
	return func(ctx Context) (*Result, error) {
		return &Result{
			Consumed: 1,
			Val:      v,
		}, nil
	}
}

func fail(err error) Parselet {
	return func(ctx Context) (*Result, error) {
		return nil, err
	}
}

func TestEmpty(t *testing.T) {
	expected := Val{Token: lexeme.Item{Type: lexeme.And, Val: "&&"}}
	parse := Empty(func(m ...Val) Val {
		return expected
	})

	ctx := Context{}
	actual, err := parse(ctx)
	if actual.Val.Token != expected.Token {
		t.Errorf("expected %v got %v", expected.Token, actual.Val.Token)
	}
	if actual.Consumed != 0 {
		t.Errorf("expected 0 tokens consumed got %d", actual.Consumed)
	}
	if err != nil {
		t.Errorf("no error expected on ampty")
	}
}

func TestToken(t *testing.T) {
	for name, test := range map[string]struct {
		ctx      Context
		rule     string
		expected Result
		err      error
	}{
		"match": {
			ctx: Context{
				Tokens: []lexeme.Item{
					{Type: lexeme.Number, Val: "5"},
				},
				Grammar: map[string]Parselet{
					"tok": Term(lexeme.Number),
				},
			},
			rule:     "tok",
			expected: Result{Consumed: 1, Val: Val{Token: lexeme.Item{Type: lexeme.Number, Val: "5"}}},
			err:      nil,
		},
		"wrong type": {
			ctx: Context{
				Tokens: []lexeme.Item{
					{Type: lexeme.Number, Val: "5"},
				},
				Grammar: map[string]Parselet{
					"tok": Term(lexeme.Symbol),
				},
			},
			rule:     "tok",
			expected: Result{},
			err:      ErrorTokenMismatch,
		},
		"eof": {
			ctx: Context{
				Tokens: []lexeme.Item{},
				Grammar: map[string]Parselet{
					"tok": Term(lexeme.Number),
				},
			},
			rule:     "tok",
			expected: Result{},
			err:      ErrorUnexpectedEof,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parse := test.ctx.Grammar[test.rule]
			actual, err := parse(test.ctx)
			if test.err == nil && err != nil {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if err != nil {
				return
			}
			if test.expected.Consumed != actual.Consumed {
				t.Fatalf("expected %d consumed got %d", test.expected.Consumed, actual.Consumed)
			}
			if !test.expected.Val.Token.CompareItem(actual.Val.Token) {
				t.Errorf("expected %v got %v", test.expected.Val.Token, actual.Val.Token)
			}
		})
	}
}

func TestZeroOrMore(t *testing.T) {
	v := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}}
	a := func(m ...Val) Val {
		vals := []parsetree.Inline{}
		for _, v := range m {
			vals = append(vals, v.Inline)
		}
		return Val{Inlines: vals}
	}
	ctx := Context{
		Grammar: map[string]Parselet{
			"0 toks": ZeroOrMore(times(0, v, ErrorTokenMismatch))(a),
			"1 toks": ZeroOrMore(times(1, v, ErrorTokenMismatch))(a),
			"2 toks": ZeroOrMore(times(2, v, ErrorTokenMismatch))(a),
			"5 toks": ZeroOrMore(times(5, v, ErrorTokenMismatch))(a),
		},
		Pos:    0,
		Tokens: []lexeme.Item{},
	}
	for name, test := range map[string]struct {
		rule     string
		expected Result
		err      error
	}{
		"0 matches": {
			rule:     "0 toks",
			expected: Result{Consumed: 0, Val: Val{Inlines: []parsetree.Inline{}}},
			err:      nil,
		},
		"1 match": {
			rule:     "1 toks",
			expected: Result{Consumed: 1, Val: Val{Inlines: []parsetree.Inline{v.Inline}}},
			err:      nil,
		},
		"2 matches": {
			rule:     "2 toks",
			expected: Result{Consumed: 2, Val: Val{Inlines: []parsetree.Inline{v.Inline, v.Inline}}},
			err:      nil,
		},
		"5 matches": {
			rule: "5 toks",
			expected: Result{Consumed: 5, Val: Val{Inlines: []parsetree.Inline{
				v.Inline, v.Inline, v.Inline, v.Inline, v.Inline,
			}}},
			err: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parse := ctx.Grammar[test.rule]
			actual, err := parse(ctx)
			if test.err == nil && err != nil {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if err != nil {
				return
			}
			if test.expected.Consumed != actual.Consumed {
				t.Fatalf("expected %d consumed got %d", test.expected.Consumed, actual.Consumed)
			}
			if len(actual.Val.Inlines) != len(test.expected.Val.Inlines) {
				t.Errorf("expected %v got %v", test.expected.Val.Inlines, actual.Val.Inlines)
			}
		})
	}
}

func TestOneOrMore(t *testing.T) {
	v := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}}
	a := func(m ...Val) Val {
		vals := []parsetree.Inline{}
		for _, v := range m {
			vals = append(vals, v.Inline)
		}
		return Val{Inlines: vals}
	}
	ctx := Context{
		Grammar: map[string]Parselet{
			"0 toks": OneOrMore(times(0, v, ErrorTokenMismatch))(a),
			"1 toks": OneOrMore(times(1, v, ErrorTokenMismatch))(a),
			"2 toks": OneOrMore(times(2, v, ErrorTokenMismatch))(a),
			"5 toks": OneOrMore(times(5, v, ErrorTokenMismatch))(a),
		},
		Pos:    0,
		Tokens: []lexeme.Item{},
	}
	for name, test := range map[string]struct {
		rule     string
		expected Result
		err      error
	}{
		"0 matches": {
			rule:     "0 toks",
			expected: Result{},
			err:      ErrorTokenMismatch,
		},
		"1 match": {
			rule:     "1 toks",
			expected: Result{Consumed: 1, Val: Val{Inlines: []parsetree.Inline{v.Inline}}},
			err:      nil,
		},
		"2 matches": {
			rule:     "2 toks",
			expected: Result{Consumed: 2, Val: Val{Inlines: []parsetree.Inline{v.Inline, v.Inline}}},
			err:      nil,
		},
		"5 matches": {
			rule: "5 toks",
			expected: Result{Consumed: 5, Val: Val{Inlines: []parsetree.Inline{
				v.Inline, v.Inline, v.Inline, v.Inline, v.Inline,
			}}},
			err: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parse := ctx.Grammar[test.rule]
			actual, err := parse(ctx)
			if test.err == nil && err != nil {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if err != nil {
				return
			}
			if test.expected.Consumed != actual.Consumed {
				t.Fatalf("expected %d consumed got %d", test.expected.Consumed, actual.Consumed)
			}
			if len(actual.Val.Inlines) != len(test.expected.Val.Inlines) {
				t.Errorf("expected %v got %v", test.expected.Val.Inlines, actual.Val.Inlines)
			}
		})
	}
}

func TestSeq(t *testing.T) {
	v := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}}
	a := func(m ...Val) Val {
		vals := []parsetree.Inline{}
		for _, v := range m {
			vals = append(vals, v.Inline)
		}
		return Val{Inlines: vals}
	}
	ctx := Context{
		Grammar: map[string]Parselet{
			"matches":   Seq(val(v), val(v), val(v), val(v))(a),
			"early err": Seq(fail(ErrorTokenMismatch), val(v), val(v), val(v))(a),
			"late err":  Seq(val(v), val(v), val(v), fail(ErrorTokenMismatch))(a),
		},
		Pos:    0,
		Tokens: []lexeme.Item{},
	}
	for name, test := range map[string]struct {
		rule     string
		expected Result
		err      error
	}{
		"matches": {
			rule: "matches",
			expected: Result{
				Consumed: 4,
				Val:      Val{Inlines: []parsetree.Inline{v.Inline, v.Inline, v.Inline, v.Inline}},
			},
			err: nil,
		},
		"early err": {
			rule:     "early err",
			expected: Result{},
			err:      ErrorTokenMismatch,
		},
		"late err": {
			rule:     "late err",
			expected: Result{},
			err:      ErrorTokenMismatch,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parse := ctx.Grammar[test.rule]
			actual, err := parse(ctx)
			if test.err == nil && err != nil {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if err != nil {
				return
			}
			if test.expected.Consumed != actual.Consumed {
				t.Fatalf("expected %d consumed got %d", test.expected.Consumed, actual.Consumed)
			}
			if len(actual.Val.Inlines) != len(test.expected.Val.Inlines) {
				t.Errorf("expected %v got %v", test.expected.Val.Inlines, actual.Val.Inlines)
			}
		})
	}
}

func TestOr(t *testing.T) {
	v1 := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}}
	v2 := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}}}
	v3 := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}}}

	ctx := Context{
		Grammar: map[string]Parselet{
			"matches first":  Or(val(v1), val(v2), val(v3)),
			"matches second": Or(fail(ErrorTokenMismatch), val(v2), val(v3)),
			"matches third":  Or(fail(ErrorTokenMismatch), fail(ErrorTokenMismatch), val(v3)),
			"matches none":   Or(fail(ErrorTokenMismatch), fail(ErrorTokenMismatch), fail(ErrorTokenMismatch)),
		},
		Pos:    0,
		Tokens: []lexeme.Item{{Type: lexeme.TextLiteral, Val: "abc"}},
	}
	for name, test := range map[string]struct {
		rule     string
		expected Result
		err      error
	}{
		"matches first": {
			rule: "matches first",
			expected: Result{
				Consumed: 1,
				Val:      v1,
			},
			err: nil,
		},
		"matches second": {
			rule: "matches second",
			expected: Result{
				Consumed: 1,
				Val:      v2,
			},
			err: nil,
		},
		"matches third": {
			rule: "matches third",
			expected: Result{
				Consumed: 1,
				Val:      v3,
			},
			err: nil,
		},
		"matches none": {
			rule:     "matches none",
			expected: Result{},
			err:      ErrorNoMatches,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parse := ctx.Grammar[test.rule]
			actual, err := parse(ctx)
			if test.err == nil && err != nil {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if !errors.Is(err, test.err) {
				t.Fatalf("Expected %v got %v", test.err, err)
			}
			if err != nil {
				return
			}
			if test.expected.Consumed != actual.Consumed {
				t.Fatalf("expected %d consumed got %d", test.expected.Consumed, actual.Consumed)
			}
			if actual.Val.Inline != test.expected.Val.Inline {
				t.Errorf("expected %v got %v", test.expected.Val.Inline, actual.Val.Inline)
			}
		})
	}
}

func TestNonterm(t *testing.T) {
	v := Val{Inline: parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}}

	ctx := Context{
		Tokens: []lexeme.Item{},
		Pos:    0,
		Grammar: map[string]Parselet{
			"abc": val(v),
			"def": Nonterm("def"),
			"ghi": Nonterm("jkl"),
		},
		Memo: map[memoKey]memoVal{
			{name: "abc", pos: 0}: {result: &Result{Consumed: 1, Val: v}, err: nil},
		},
		Visited: map[memoKey]bool{},
	}

	parse := Nonterm("abc")
	r, err := parse(ctx)
	if err != nil {
		t.Fatalf("no error expected on memoize")
	}
	if r.Consumed != 1 {
		t.Errorf("expected 1 consumed, got %d", r.Consumed)
	}
	if !r.Val.Inline.CompareInline(v.Inline) {
		t.Errorf("expected %v got %v", v.Inline, r.Val.Inline)
	}

	parse = Nonterm("def")
	_, err = parse(ctx)
	if !errors.Is(err, ErrorLeftRecursion) {
		t.Errorf("expected error on left recursion")
	}

	parse = Nonterm("ghi")
	_, err = parse(ctx)
	if !errors.Is(err, ErrorBadRuleName) {
		t.Errorf("expected error on left recursion")
	}
}
