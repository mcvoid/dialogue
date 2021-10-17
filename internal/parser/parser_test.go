package parser

import (
	"testing"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

type mockLexer []lexeme.Item

func (l mockLexer) Lex() []lexeme.Item {
	return l
}

func TestParser(t *testing.T) {
	parser := New()
	tokens := []lexeme.Item{
		{Type: lexeme.CloseCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.OpenInlineCode, Val: "`"},
		{Type: lexeme.Inc, Val: "++"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Star, Val: "*"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Plus, Val: "+"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.CloseInlineCode, Val: "`"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenInlineCode, Val: "`"},
		{Type: lexeme.Symbol, Val: "def"},
		{Type: lexeme.CloseInlineCode, Val: "`"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "def"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "ghi"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "jkl"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.Eq, Val: "="},
		{Type: lexeme.Inc, Val: "++"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Star, Val: "*"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.GotoLiteral, Val: "goto"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.CloseCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Eof, Val: ""},
	}
	actual, err := parser.Parse(mockLexer(tokens))
	if err != nil {
		t.Fatalf("no error expected, got %v", err)
	}
	expected := parsetree.Script{
		FrontMatter: parsetree.FrontMatter{
			FuncDecls: []parsetree.FuncDecl{},
			Delimiter: lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
			EndLine:   lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
		},
		Nodes: []parsetree.Node{
			{
				Header: parsetree.Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []parsetree.Block{
					parsetree.Paragraph{
						Lines: []parsetree.Line{
							{
								Items: []parsetree.Inline{
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
							{
								Items: []parsetree.Inline{
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
									parsetree.InlineCode{
										CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
										Code: parsetree.BinaryExpression{
											LeftOperand: parsetree.BinaryExpression{
												LeftOperand: parsetree.UnaryExpression{
													Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
													Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
												},
												Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
												RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
											},
											Operator:     lexeme.Item{Type: lexeme.Plus, Val: "+"},
											RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
										},
										CodeEnd: lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
									},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
							{
								Items: []parsetree.Inline{
									parsetree.InlineCode{
										CodeStart: lexeme.Item{Type: lexeme.OpenInlineCode, Val: "`"},
										Code:      parsetree.Literal{Value: lexeme.Item{Type: lexeme.Symbol, Val: "def"}},
										CodeEnd:   lexeme.Item{Type: lexeme.CloseInlineCode, Val: "`"},
									},
									parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "abc"}},
								},
								EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
			},
			{
				Header: parsetree.Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "def"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []parsetree.Block{
					parsetree.List{
						Links: []parsetree.ListItem{
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
							{
								Prefix: lexeme.Item{Type: lexeme.ListItemPrefix, Val: "-"},
								Link: parsetree.Link{
									OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
									Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
									CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
									OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
									Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
									CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
									EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
								},
							},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
			},
			{
				Header: parsetree.Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "ghi"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []parsetree.Block{
					parsetree.LinkBlock{
						Link: parsetree.Link{
							OpenBrace:  lexeme.Item{Type: lexeme.OpenSquareBrace, Val: "["},
							Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
							CloseBrace: lexeme.Item{Type: lexeme.CloseSquareBrace, Val: "]"},
							OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
							Text:       parsetree.Text{Text: lexeme.Item{Type: lexeme.TextLiteral, Val: "def"}},
							CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
							EndLine:    lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						},
						EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
			},
			{
				Header: parsetree.Header{
					Hash:    lexeme.Item{Type: lexeme.Hash, Val: "#"},
					Name:    lexeme.Item{Type: lexeme.Symbol, Val: "jkl"},
					EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				},
				EndLine: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
				Blocks: []parsetree.Block{
					parsetree.CodeBlock{
						StartFence:   lexeme.Item{Type: lexeme.OpenCodeFence, Val: "```"},
						StartEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						Code: []parsetree.Statement{
							parsetree.FunctionCall{
								Symbol:     lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								OpenParen:  lexeme.Item{Type: lexeme.OpenParen, Val: "("},
								Args:       []parsetree.Expression{},
								CloseParen: lexeme.Item{Type: lexeme.CloseParen, Val: ")"},
								Semicolon:  lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
							parsetree.Assignment{
								Symbol:    lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								EqualSign: lexeme.Item{Type: lexeme.Eq, Val: "="},
								Value: parsetree.BinaryExpression{
									LeftOperand: parsetree.UnaryExpression{
										Operator: lexeme.Item{Type: lexeme.Inc, Val: "++"},
										Operand:  parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
									},
									Operator:     lexeme.Item{Type: lexeme.Star, Val: "*"},
									RightOperand: parsetree.Literal{Value: lexeme.Item{Type: lexeme.Number, Val: "5"}},
								},
								Semicolon: lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
							parsetree.Goto{
								GotoLiteral: lexeme.Item{Type: lexeme.GotoLiteral, Val: "goto"},
								Symbol:      lexeme.Item{Type: lexeme.Symbol, Val: "abc"},
								Semicolon:   lexeme.Item{Type: lexeme.Semicolon, Val: ";"},
							},
						},
						EndFence:     lexeme.Item{Type: lexeme.CloseCodeFence, Val: "```"},
						CloseEndline: lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
						EndLine:      lexeme.Item{Type: lexeme.LineBreak, Val: "\n"},
					},
				},
			},
		},
	}
	if !expected.CompareScript(actual) {
		t.Errorf("expected %v got %v", expected, actual)
	}

	tokens = []lexeme.Item{
		{Type: lexeme.CloseCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.OpenInlineCode, Val: "`"},
		{Type: lexeme.Inc, Val: "++"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Star, Val: "*"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Plus, Val: "+"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.CloseInlineCode, Val: "`"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenInlineCode, Val: "`"},
		{Type: lexeme.Symbol, Val: "def"},
		{Type: lexeme.CloseInlineCode, Val: "`"},
		{Type: lexeme.TextLiteral, Val: "abc"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "def"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.ListItemPrefix, Val: "-"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "ghi"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenSquareBrace, Val: "["},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.CloseSquareBrace, Val: "]"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.TextLiteral, Val: "def"},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Hash, Val: "#"},
		{Type: lexeme.Symbol, Val: "jkl"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.OpenCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.OpenParen, Val: "("},
		{Type: lexeme.CloseParen, Val: ")"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.Eq, Val: "="},
		{Type: lexeme.Inc, Val: "++"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Star, Val: "*"},
		{Type: lexeme.Number, Val: "5"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.GotoLiteral, Val: "goto"},
		{Type: lexeme.Symbol, Val: "abc"},
		{Type: lexeme.Semicolon, Val: ";"},
		{Type: lexeme.CloseCodeFence, Val: "```"},
		{Type: lexeme.LineBreak, Val: "\n"},
		{Type: lexeme.LineBreak, Val: "\n"},
	}
	_, err = parser.Parse(mockLexer(tokens))
	if err == nil {
		t.Errorf("error expected")
	}
}
