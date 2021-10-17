package parsetree

import (
	"fmt"
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

func TestBlocks(t *testing.T) {
	for i, test := range []struct {
		a        Block
		b        Block
		expected bool
	}{
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: true,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "("},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "ac"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: ")"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "["},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "de"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: "]"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
			},
			expected: false,
		},
		{
			a: LinkBlock{
				Link: Link{
					OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
					Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
					OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
					Text:       Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
					CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
					EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b:        CodeBlock{},
			expected: false,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: true,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "--"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "ab"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
			},
			expected: false,
		},
		{
			a: List{
				Links: []ListItem{
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
					{
						Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
						Link:   Link{Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}}},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b:        LinkBlock{},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: true,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "``"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "ab"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "``"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
			},
			expected: false,
		},
		{
			a: CodeBlock{
				StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
				StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Code: []Statement{
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
					Goto{Symbol: lexeme.Item{Type: lexeme.Symbol, Val: "abc"}},
				},
				EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
				CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b:        LinkBlock{},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: true,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "ac"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "'"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "ef"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "'"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			expected: false,
		},
		{
			a: Paragraph{
				Lines: []Line{
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
					{
						Items: []Inline{
							Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
							InlineCode{
								CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
								Code:      Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
								CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
			},
			b:        LinkBlock{},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareBlock(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestNode(t *testing.T) {
	for i, test := range []struct {
		a        Node
		b        Node
		expected bool
	}{
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: true,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "##"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "bc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
		{
			a: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			b: Node{
				Header: Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []Block{
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "ef"}},
						},
					},
					LinkBlock{
						Link: Link{
							Text: Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
						},
					},
				},
			},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareNode(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func TestScript(t *testing.T) {
	for i, test := range []struct {
		a        Script
		b        Script
		expected bool
	}{
		{
			a: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: true,
		},
		{
			a: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "##"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				Nodes: []Node{
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
					{Header: Header{Hash: lexeme.Item{Type: lexeme.Hash, Val: "#"}}},
				},
				Eof: lexeme.Item{Type: lexeme.Hash, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: true,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							ExternKeyword: lexeme.Item{Type: lexeme.ExternKeyword, Val: "extern"},
							Symbol:        lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen:     lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							ExternKeyword: lexeme.Item{Type: lexeme.ExternKeyword, Val: "external"},
							Symbol:        lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen:     lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "def"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: ")"},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: "("},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ":"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "/n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
		{
			a: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			b: Script{
				FrontMatter: FrontMatter{
					FuncDecls: []FuncDecl{
						{
							Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							OpenParen: lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Params: []lexeme.Item{
								{Type: lexeme.Type, Val: "bool"},
								{Type: lexeme.Type, Val: "number"},
								{Type: lexeme.Type, Val: "string"},
								{Type: lexeme.Type, Val: "null"},
							},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
					},
					Delimiter: lexeme.Item{Type: lexeme.OpenCodeFence},
					EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				Nodes: []Node{},
				Eof:   lexeme.Item{Type: lexeme.Eof, Val: ""},
			},
			expected: false,
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := test.a.CompareScript(test.b)
			if test.expected != actual {
				t.Errorf("expected %v got %v", test.expected, actual)
			}
		})
	}
}
