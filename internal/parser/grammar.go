package parser

import (
	"github.com/mcvoid/dialogue/internal/types/lexeme"
	"github.com/mcvoid/dialogue/internal/types/parsetree"
)

var Grammar = map[string]Parselet{
	"script": Seq(Nonterm("frontmatter"), Nonterm("nodes"), Term(lexeme.Eof))(func(m ...Val) Val {
		return Val{Script: parsetree.Script{
			FrontMatter: m[0].FrontMatter,
			Nodes:       m[1].Nodes,
			Eof:         m[2].Token,
		}}
	}),
	"frontmatter": Or(
		Seq(Nonterm("funcdecls"), Term(lexeme.CloseCodeFence), Term(lexeme.LineBreak))(func(m ...Val) Val {
			return Val{FrontMatter: parsetree.FrontMatter{
				FuncDecls: m[0].FuncDecls,
				Delimiter: m[1].Token,
				EndLine:   m[2].Token,
			}}
		}),
		Empty(func(m ...Val) Val {
			return Val{FrontMatter: parsetree.FrontMatter{
				FuncDecls: []parsetree.FuncDecl{},
			}}
		}),
	),
	"funcdecls": ZeroOrMore(Nonterm("funcdecl"))(func(m ...Val) Val {
		vals := []parsetree.FuncDecl{}
		for _, v := range m {
			vals = append(vals, v.FuncDecl)
		}
		return Val{FuncDecls: vals}
	}),
	"funcdecl": Seq(
		Term(lexeme.ExternKeyword),
		Term(lexeme.Symbol),
		Term(lexeme.OpenParen),
		Nonterm("params"),
		Term(lexeme.CloseParen),
		Term(lexeme.Semicolon),
		Term(lexeme.LineBreak),
	)(func(m ...Val) Val {
		return Val{FuncDecl: parsetree.FuncDecl{
			ExternKeyword: m[0].Token,
			Symbol:        m[1].Token,
			OpenParen:     m[2].Token,
			Params:        m[3].Params,
			CloseParen:    m[4].Token,
			Semicolon:     m[5].Token,
			EndLine:       m[6].Token,
		}}
	}),
	"params": Or(
		Seq(Term(lexeme.Type), Nonterm("restparams"))(func(m ...Val) Val {
			return Val{Params: append([]lexeme.Item{m[0].Token}, m[1].Params...)}
		}),
		Empty(func(m ...Val) Val {
			return Val{Params: []lexeme.Item{}}
		}),
	),
	"restparams": ZeroOrMore(Nonterm("restparam"))(func(m ...Val) Val {
		vals := []lexeme.Item{}
		for _, v := range m {
			vals = append(vals, v.Token)
		}
		return Val{Params: vals}
	}),
	"restparam": Seq(Term(lexeme.Comma), Term(lexeme.Type))(func(m ...Val) Val {
		return m[1]
	}),
	"nodes": OneOrMore(Nonterm("node"))(func(m ...Val) Val {
		vals := []parsetree.Node{}
		for _, v := range m {
			vals = append(vals, v.Node)
		}
		return Val{Nodes: vals}
	}),
	"node": Seq(Nonterm("header"), Term(lexeme.LineBreak), Nonterm("blocks"))(func(m ...Val) Val {
		return Val{Node: parsetree.Node{
			Header:  m[0].Header,
			EndLine: m[1].Token,
			Blocks:  m[2].Blocks,
		}}
	}),
	"header": Seq(Term(lexeme.Hash), Term(lexeme.Symbol), Term(lexeme.LineBreak))(func(m ...Val) Val {
		return Val{Header: parsetree.Header{
			Hash:    m[0].Token,
			Name:    m[1].Token,
			EndLine: m[2].Token,
		}}
	}),
	"blocks": OneOrMore(Nonterm("block"))(func(m ...Val) Val {
		vals := []parsetree.Block{}
		for _, v := range m {
			vals = append(vals, v.Block)
		}
		return Val{Blocks: vals}
	}),
	"block": Or(
		Nonterm("paragraph"),
		Nonterm("codeBlock"),
		Nonterm("linkBlock"),
		Nonterm("list"),
	),
	"paragraph": Seq(Nonterm("lines"), Term(lexeme.LineBreak))(func(m ...Val) Val {
		return Val{Block: parsetree.Paragraph{
			Lines:   m[0].Lines,
			EndLine: m[1].Token,
		}}
	}),
	"lines": OneOrMore(Nonterm("line"))(func(m ...Val) Val {
		vals := []parsetree.Line{}
		for _, v := range m {
			vals = append(vals, v.Line)
		}
		return Val{Lines: vals}
	}),
	"line": Seq(Nonterm("inlines"), Term(lexeme.LineBreak))(func(m ...Val) Val {
		return Val{Line: parsetree.Line{
			Items:   m[0].Inlines,
			EndLine: m[1].Token,
		}}
	}),
	"inlines": OneOrMore(Nonterm("inline"))(func(m ...Val) Val {
		vals := []parsetree.Inline{}
		for _, v := range m {
			vals = append(vals, v.Inline)
		}
		return Val{Inlines: vals}
	}),
	"inline": Or(
		Nonterm("text"),
		Nonterm("inlineCode"),
	),
	"text": Seq(Term(lexeme.TextLiteral))(func(m ...Val) Val {
		return Val{Inline: parsetree.Text{
			Text: m[0].Token,
		}}
	}),
	"inlineCode": Seq(Term(lexeme.OpenInlineCode), Nonterm("expression"), Term(lexeme.CloseInlineCode))(func(m ...Val) Val {
		return Val{Inline: parsetree.InlineCode{
			CodeStart: m[0].Token,
			Code:      m[1].Expression,
			CodeEnd:   m[2].Token,
		}}
	}),
	"list": Seq(Nonterm("listItems"), Term(lexeme.LineBreak))(func(m ...Val) Val {
		return Val{Block: parsetree.List{
			Links:   m[0].ListItems,
			EndLine: m[1].Token,
		}}
	}),
	"listItems": OneOrMore(Nonterm("listItem"))(func(m ...Val) Val {
		vals := []parsetree.ListItem{}
		for _, v := range m {
			vals = append(vals, v.ListItem)
		}
		return Val{ListItems: vals}
	}),
	"listItem": Seq(Term(lexeme.ListItemPrefix), Nonterm("link"))(func(m ...Val) Val {
		return Val{ListItem: parsetree.ListItem{
			Prefix: m[0].Token,
			Link:   m[1].Link,
		}}
	}),
	"linkBlock": Seq(Nonterm("link"), Term(lexeme.LineBreak))(func(m ...Val) Val {
		return Val{Block: parsetree.LinkBlock{Link: m[0].Link, EndLine: m[1].Token}}
	}),
	"link": Seq(
		Term(lexeme.OpenSquareBrace),
		Term(lexeme.Symbol),
		Term(lexeme.CloseSquareBrace),
		Term(lexeme.OpenParen),
		Nonterm("text"),
		Term(lexeme.CloseParen),
		Term(lexeme.LineBreak),
	)(func(m ...Val) Val {
		return Val{Link: parsetree.Link{
			OpenBrace:  m[0].Token,
			Symbol:     m[1].Token,
			CloseBrace: m[2].Token,
			OpenParen:  m[3].Token,
			Text:       m[4].Inline,
			CloseParen: m[5].Token,
			EndLine:    m[6].Token,
		}}
	}),
	"codeBlock": Seq(
		Term(lexeme.OpenCodeFence),
		Term(lexeme.LineBreak),
		Nonterm("statements"),
		Term(lexeme.CloseCodeFence),
		Term(lexeme.LineBreak),
		Term(lexeme.LineBreak),
	)(func(m ...Val) Val {
		return Val{Block: parsetree.CodeBlock{
			StartFence:   m[0].Token,
			StartEndline: m[1].Token,
			Code:         m[2].Statements,
			EndFence:     m[3].Token,
			CloseEndline: m[4].Token,
			EndLine:      m[5].Token,
		}}
	}),
	"statements": OneOrMore(Nonterm("statement"))(func(m ...Val) Val {
		vals := []parsetree.Statement{}
		for _, v := range m {
			vals = append(vals, v.Statement)
		}
		return Val{Statements: vals}
	}),
	"statement": Or(
		Nonterm("conditionalWithElse"),
		Nonterm("conditional"),
		Nonterm("statementBlock"),
		Nonterm("functionCall"),
		Nonterm("goto"),
		Nonterm("loop"),
		Nonterm("assignment"),
	),
	"conditional": Seq(Term(lexeme.IfLiteral), Nonterm("expression"), Nonterm("statementBlock"))(func(m ...Val) Val {
		return Val{Statement: parsetree.Conditional{
			IfLiteral:  m[0].Token,
			Cond:       m[1].Expression,
			Consequent: m[2].Statement.(parsetree.StatementBlock),
		}}
	}),
	"conditionalWithElse": Seq(
		Term(lexeme.IfLiteral),
		Nonterm("expression"),
		Nonterm("statementBlock"),
		Term(lexeme.ElseLiteral),
		Nonterm("statementBlock"),
	)(func(m ...Val) Val {
		return Val{Statement: parsetree.ConditionalWithElse{
			IfLiteral:   m[0].Token,
			Cond:        m[1].Expression,
			Consequent:  m[2].Statement.(parsetree.StatementBlock),
			ElseLiteral: m[3].Token,
			Alternate:   m[4].Statement.(parsetree.StatementBlock),
		}}
	}),
	"statementBlock": Seq(
		Term(lexeme.OpenCurlyBrace),
		Nonterm("statements"),
		Term(lexeme.CloseCurlyBrace),
	)(func(m ...Val) Val {
		return Val{Statement: parsetree.StatementBlock{
			OpenBrace:  m[0].Token,
			Statements: m[1].Statements,
			CloseBrace: m[2].Token,
		}}
	}),
	"functionCall": Seq(
		Term(lexeme.Symbol),
		Term(lexeme.OpenParen),
		Nonterm("funcArgList"),
		Term(lexeme.CloseParen),
		Term(lexeme.Semicolon),
	)(func(m ...Val) Val {
		return Val{Statement: parsetree.FunctionCall{
			Symbol:     m[0].Token,
			OpenParen:  m[1].Token,
			Args:       m[2].FuncArgsList,
			CloseParen: m[3].Token,
			Semicolon:  m[4].Token,
		}}
	}),
	"funcArgList": Or(
		Seq(Nonterm("expression"), Nonterm("restArgs"))(func(m ...Val) Val {
			return Val{FuncArgsList: append([]parsetree.Expression{m[0].Expression}, m[1].FuncArgsList...)}
		}),
		Empty(func(m ...Val) Val {
			return Val{FuncArgsList: []parsetree.Expression{}}
		}),
	),
	"restArgs": ZeroOrMore(Nonterm("restArg"))(func(m ...Val) Val {
		vals := []parsetree.Expression{}
		for _, v := range m {
			vals = append(vals, v.Expression)
		}
		return Val{FuncArgsList: vals}
	}),
	"restArg": Seq(Term(lexeme.Comma), Nonterm("expression"))(func(m ...Val) Val {
		// dropping the commas from the parse tree
		// in exchange for a simpler / more readable tree
		return m[1]
	}),
	"goto": Seq(Term(lexeme.GotoLiteral), Term(lexeme.Symbol), Term(lexeme.Semicolon))(func(m ...Val) Val {
		return Val{Statement: parsetree.Goto{
			GotoLiteral: m[0].Token,
			Symbol:      m[1].Token,
			Semicolon:   m[2].Token,
		}}
	}),
	"assignment": Seq(
		Term(lexeme.Symbol),
		Term(lexeme.Eq),
		Nonterm("expression"),
		Term(lexeme.Semicolon),
	)(func(m ...Val) Val {
		return Val{Statement: parsetree.Assignment{
			Symbol:    m[0].Token,
			EqualSign: m[1].Token,
			Value:     m[2].Expression,
			Semicolon: m[3].Token,
		}}
	}),
	"loop": Seq(
		Term(lexeme.WhileLiteral),
		Nonterm("expression"),
		Nonterm("statementBlock"),
	)(func(m ...Val) Val {
		return Val{Statement: parsetree.Loop{
			WhileLiteral: m[0].Token,
			Cond:         m[1].Expression,
			Body:         m[2].Statement,
		}}
	}),
	"expression": Or(
		Seq(Nonterm("andComparator"), Term(lexeme.Or), Nonterm("expression"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("andComparator"),
	),
	"andComparator": Or(
		Seq(Nonterm("eqComparator"), Term(lexeme.And), Nonterm("andComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("eqComparator"),
	),
	"eqComparator": Or(
		Seq(Nonterm("ineqComparator"), Term(lexeme.DoubleEq), Nonterm("eqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("ineqComparator"), Term(lexeme.Neq), Nonterm("eqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("ineqComparator"),
	),
	"ineqComparator": Or(
		Seq(Nonterm("term"), Term(lexeme.Gt), Nonterm("ineqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("term"), Term(lexeme.Gte), Nonterm("ineqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("term"), Term(lexeme.Lt), Nonterm("ineqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("term"), Term(lexeme.Lte), Nonterm("ineqComparator"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("term"),
	),
	"term": Or(
		Seq(Nonterm("factor"), Term(lexeme.Plus), Nonterm("term"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("factor"), Term(lexeme.Minus), Nonterm("term"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("factor"), Term(lexeme.Dot), Nonterm("term"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("factor"),
	),
	"factor": Or(
		Seq(Nonterm("unary"), Term(lexeme.Star), Nonterm("factor"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("unary"), Term(lexeme.Slash), Nonterm("factor"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Seq(Nonterm("unary"), Term(lexeme.Percent), Nonterm("factor"))(func(m ...Val) Val {
			return Val{Expression: parsetree.BinaryExpression{
				LeftOperand:  m[0].Expression,
				Operator:     m[1].Token,
				RightOperand: m[2].Expression,
			}}
		}),
		Nonterm("unary"),
	),
	"unary": Or(
		Seq(Term(lexeme.Inc), Nonterm("value"))(func(m ...Val) Val {
			return Val{Expression: parsetree.UnaryExpression{
				Operator: m[0].Token,
				Operand:  m[1].Expression,
			}}
		}),
		Seq(Term(lexeme.Dec), Nonterm("value"))(func(m ...Val) Val {
			return Val{Expression: parsetree.UnaryExpression{
				Operator: m[0].Token,
				Operand:  m[1].Expression,
			}}
		}),
		Seq(Term(lexeme.Not), Nonterm("value"))(func(m ...Val) Val {
			return Val{Expression: parsetree.UnaryExpression{
				Operator: m[0].Token,
				Operand:  m[1].Expression,
			}}
		}),
		Seq(Term(lexeme.Minus), Nonterm("value"))(func(m ...Val) Val {
			return Val{Expression: parsetree.UnaryExpression{
				Operator: m[0].Token,
				Operand:  m[1].Expression,
			}}
		}),
		Nonterm("value"),
	),
	"value": Or(
		Nonterm("nested"),
		Nonterm("literal"),
	),
	"nested": Seq(Term(lexeme.OpenParen), Nonterm("expression"), Term(lexeme.CloseParen))(func(m ...Val) Val {
		return Val{Expression: parsetree.NestedExpression{
			OpenParen:  m[0].Token,
			Expr:       m[1].Expression,
			CloseParen: m[2].Token,
		}}
	}),
	"literal": Or(
		Seq(Term(lexeme.Number))(func(m ...Val) Val {
			return Val{Expression: parsetree.Literal{Value: m[0].Token}}
		}),
		Seq(Term(lexeme.Boolean))(func(m ...Val) Val {
			return Val{Expression: parsetree.Literal{Value: m[0].Token}}
		}),
		Seq(Term(lexeme.Symbol))(func(m ...Val) Val {
			return Val{Expression: parsetree.Literal{Value: m[0].Token}}
		}),
		Seq(Term(lexeme.String))(func(m ...Val) Val {
			return Val{Expression: parsetree.Literal{Value: m[0].Token}}
		}),
		Seq(Term(lexeme.Null))(func(m ...Val) Val {
			return Val{Expression: parsetree.Literal{Value: m[0].Token}}
		}),
	),
}
