package lexer

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

func TestLexText(t *testing.T) {
	tests := map[string]struct {
		input  string
		tokens []lexeme.Item
	}{
		"empty": {
			input: ``,
			tokens: []lexeme.Item{
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"header": {
			input: `# abcABC123
`,
			tokens: []lexeme.Item{
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abcABC123"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"text, link, list": {
			input: `
abc def ghi
[ abc123 ](abc def ghi)
- [abc123](abc def ghi)
abc`,
			tokens: []lexeme.Item{
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc def ghi"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "abc def ghi"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.ListItemPrefix, Val: "-"},
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "abc def ghi"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"link error 1": {
			input: "[9abc](abc)\n",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Error, Val: ErrorBadLink},
			},
		},
		"link error 2": {
			input: "[abc](abc\n)",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Error, Val: ErrorBadLink},
			},
		},
		"link error 3": {
			input: "[abc](abc",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Error, Val: ErrorBadLink},
			},
		},
		"link error 4": {
			input: "[abc](abc)",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenSquareBrace, Val: "["},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.CloseSquareBrace, Val: "]"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.TextLiteral, Val: "abc"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Error, Val: ErrorBadLink},
			},
		},
		"header error 1": {
			input: "# abc",
			tokens: []lexeme.Item{
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.Error, Val: ErrorBadHeader},
			},
		},
		"header error 2": {
			input: "# 9abc\n",
			tokens: []lexeme.Item{
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Error, Val: ErrorBadHeader},
			},
		},
		"inline code": {
			input: "`abc123`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"inline code surrounded by text": {
			input: "text`abc123`text",
			tokens: []lexeme.Item{
				{Type: lexeme.TextLiteral, Val: "text"},
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.TextLiteral, Val: "text"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"code fence": {
			input: "```\nabc123```\n",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"bogus character": {
			input: "```\n?```\n",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Error, Val: ErrorBadCode},
			},
		},
		"inline number": {
			input: "`123 123e-15 0`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Number, Val: "123"},
				{Type: lexeme.Number, Val: "123e-15"},
				{Type: lexeme.Number, Val: "0"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"bad number 1": {
			input: "`0123`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadNumber},
			},
		},
		"bad number 2": {
			input: "`123e+a`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadNumber},
			},
		},
		"incomplete code": {
			input: "`s",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Symbol, Val: "s"},
				{Type: lexeme.Error, Val: ErrorBadCode},
			},
		},
		"operators": {
			input: "`if else while true false null goto ++ -- && || >= <= { } ( ) + - * / > < = . , ; ! % !=`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.IfLiteral, Val: IfLiteral},
				{Type: lexeme.ElseLiteral, Val: ElseLiteral},
				{Type: lexeme.WhileLiteral, Val: WhileLiteral},
				{Type: lexeme.Boolean, Val: TrueLiteral},
				{Type: lexeme.Boolean, Val: FalseLiteral},
				{Type: lexeme.Null, Val: NullLiteral},
				{Type: lexeme.GotoLiteral, Val: GotoLiteral},
				{Type: lexeme.Inc, Val: PlusPlus},
				{Type: lexeme.Dec, Val: MinusMinus},
				{Type: lexeme.And, Val: And},
				{Type: lexeme.Or, Val: Or},
				{Type: lexeme.Gte, Val: Gte},
				{Type: lexeme.Lte, Val: Lte},
				{Type: lexeme.OpenCurlyBrace, Val: OpenCurlyBrace},
				{Type: lexeme.CloseCurlyBrace, Val: CloseCurlyBrace},
				{Type: lexeme.OpenParen, Val: OpenParen},
				{Type: lexeme.CloseParen, Val: CloseParen},
				{Type: lexeme.Plus, Val: Plus},
				{Type: lexeme.Minus, Val: Minus},
				{Type: lexeme.Star, Val: Star},
				{Type: lexeme.Slash, Val: Slash},
				{Type: lexeme.Gt, Val: Gt},
				{Type: lexeme.Lt, Val: Lt},
				{Type: lexeme.Eq, Val: Eq},
				{Type: lexeme.Dot, Val: Dot},
				{Type: lexeme.Comma, Val: Comma},
				{Type: lexeme.Semicolon, Val: Semicolon},
				{Type: lexeme.Not, Val: Not},
				{Type: lexeme.Percent, Val: Percent},
				{Type: lexeme.Neq, Val: Neq},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"bad operator": {
			input: "`&`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadOperator},
			},
		},
		"comment": {
			input: "`// abc 123 +-- \n`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"bad comment": {
			input: "`// abc 123 +--`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadCode},
			},
		},
		"string unescaped": {
			input: "`\"// abc 123 +--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.String, Val: "\"// abc 123 +--\""},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"string escaped": {
			input: "`\"// abc 123 \\n+--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.String, Val: "\"// abc 123 \\n+--\""},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"string unicode escaped": {
			input: "`\"// abc 123 \\u0aF5 +--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.String, Val: "\"// abc 123 \\u0aF5 +--\""},
				{Type: lexeme.CloseInlineCode, Val: "`"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"bad string 1": {
			input: "`\"// abc 123 \n+--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadString},
			},
		},
		"bad string 2": {
			input: "`\"// abc 123 +--\\\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadString},
			},
		},
		"bad string 3": {
			input: "`\"// abc 123 \\j +--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadEscape},
			},
		},
		"bad string 4": {
			input: "`\"// abc 123 \\u123 +--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadEscape},
			},
		},
		"bad string 5": {
			input: "`\"// abc 123 \\u123j +--\"`",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenInlineCode, Val: "`"},
				{Type: lexeme.Error, Val: ErrorBadEscape},
			},
		},
		"bad code block 1": {
			input: "```abc123```\n",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.Error, Val: ErrorBadCode},
			},
		},
		"bad code block 2": {
			input: "```\nabc123```",
			tokens: []lexeme.Item{
				{Type: lexeme.OpenCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Symbol, Val: "abc123"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.Error, Val: ErrorBadCode},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := New(test.input)
			l.state = LexLine
			actual := l.Lex()
			if len(actual) != len(test.tokens) {
				t.Fatalf("expected\n%v\ngot\n%v\n", test.tokens, actual)
			}
			for i := range test.tokens {
				if !test.tokens[i].CompareItem(actual[i]) {
					t.Errorf(
						"token %d: expected [%v, %v] got [%v, %v]",
						i,
						test.tokens[i].Type,
						test.tokens[i].Val,
						actual[i].Type,
						actual[i].Val,
					)
				}
			}
		})
	}
}

func TestLexer(t *testing.T) {
	tests := map[string]struct {
		input  string
		tokens []lexeme.Item
	}{
		"with frontmatter": {
			input: "\n  extern abc(bool, number, null, string);\n//fdsghafgkafghjsaghjfsa\n```\n\n# abc\n\n",
			tokens: []lexeme.Item{
				{Type: lexeme.ExternKeyword, Val: "extern"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.OpenParen, Val: "("},
				{Type: lexeme.Type, Val: "bool"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Type, Val: "number"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Type, Val: "null"},
				{Type: lexeme.Comma, Val: ","},
				{Type: lexeme.Type, Val: "string"},
				{Type: lexeme.CloseParen, Val: ")"},
				{Type: lexeme.Semicolon, Val: ";"},
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Hash, Val: "#"},
				{Type: lexeme.Symbol, Val: "abc"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.LineBreak, Val: "\n"},
				{Type: lexeme.Eof, Val: ""},
			},
		},
		"no enline after frontmatter": {
			input: "```",
			tokens: []lexeme.Item{
				{Type: lexeme.CloseCodeFence, Val: "```"},
				{Type: lexeme.Error, Val: ErrorBadFrontMatterEnd},
			},
		},
		"unterminated comment": {
			input: "//fdjskl;fj;lzf",
			tokens: []lexeme.Item{
				{Type: lexeme.Error, Val: ErrorBadFrontMatterEnd},
			},
		},
		"no frontmatter close": {
			input: "# newnode",
			tokens: []lexeme.Item{
				{Type: lexeme.Error, Val: ErrorBadFrontmatter},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := New(test.input)
			actual := l.Lex()
			if len(actual) != len(test.tokens) {
				t.Fatalf("expected\n%v\ngot\n%v\n", test.tokens, actual)
			}
			for i := range test.tokens {
				if !test.tokens[i].CompareItem(actual[i]) {
					t.Errorf(
						"token %d: expected [%v, %v] got [%v, %v]",
						i,
						test.tokens[i].Type,
						test.tokens[i].Val,
						actual[i].Type,
						actual[i].Val,
					)
				}
			}
		})
	}
}
