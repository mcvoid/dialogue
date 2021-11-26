package lexer

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

type State func(*Lexer) State

const (
	Hash                     = "#"
	LineEnd                  = "\n"
	CodeDelimiter            = "`"
	FencedCodeBlockDelimiter = "```"
	OpenSquareBracket        = "["
	CloseSquareBracket       = "]"
	OpenParen                = "("
	CloseParen               = ")"
	OpenCurlyBrace           = "{"
	CloseCurlyBrace          = "}"
	Comma                    = ","
	UnorderedListPrefix      = "-"
	Whitespace               = " \t"
	SymbolStart              = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SymbolTail               = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberStart              = "0123456789"
	StringStart              = "\""
	PlusPlus                 = "++"
	MinusMinus               = "--"
	Not                      = "!"
	Plus                     = "+"
	Minus                    = "-"
	Star                     = "*"
	Slash                    = "/"
	Percent                  = "%"
	Gte                      = ">="
	Gt                       = ">"
	Lte                      = "<="
	Lt                       = "<"
	Eq                       = "="
	EqEq                     = "=="
	Neq                      = "!="
	TrueLiteral              = "true"
	FalseLiteral             = "false"
	NullLiteral              = "null"
	GotoLiteral              = "goto"
	Dot                      = "."
	Semicolon                = ";"
	And                      = "&&"
	Or                       = "||"
	Operators                = "!+-*/><=&|{}().,;%"
	IfLiteral                = "if"
	ElseLiteral              = "else"
	WhileLiteral             = "while"
	Comment                  = "//"
	BoolType                 = "bool"
	NumberType               = "number"
	StringType               = "string"
	NullType                 = "null"
	ExternKeyword            = "extern"
)

const (
	ErrorBadHeader         = "Header must only be of the form '# HeaderName\\n"
	ErrorBadLink           = "Link must only be of the form '[symbol] (text)\\n"
	ErrorBadCode           = "Unrecognized code element"
	ErrorBadNumber         = "Numbers must be in format -?0|([1-9][0-9]*(e[+-]?[0-9])?)"
	ErrorBadOperator       = "Not a valid operator"
	ErrorBadString         = "Not a valid string"
	ErrorBadEscape         = "The only valid escapes are \\\\, \\/, \\b, \\f, \\n, \\r, \\t, \\uxxxx"
	ErrorBadFrontMatterEnd = "Fromtmatter must end in mewline"
	ErrorBadFrontmatter    = "Unrecognized token in frontmatter"
)

func LexFrontMatter(l *Lexer) State {
	if strings.HasPrefix(l.input[l.pos:], FencedCodeBlockDelimiter) {
		l.pos += len(FencedCodeBlockDelimiter)
		emit(l, lexeme.CloseCodeFence)
		if !accept(l, LineEnd) {
			return errorf(l, ErrorBadFrontMatterEnd)
		}
		emit(l, lexeme.LineBreak)
		return LexLine
	}

	if acceptRun(l, LineEnd) {
		ignore(l)
		return LexFrontMatter
	}
	if acceptRun(l, Whitespace) {
		ignore(l)
		return LexFrontMatter
	}
	if accept(l, OpenParen) {
		emit(l, lexeme.OpenParen)
		return LexFrontMatter
	}
	if accept(l, CloseParen) {
		emit(l, lexeme.CloseParen)
		return LexFrontMatter
	}
	if accept(l, Comma) {
		emit(l, lexeme.Comma)
		return LexFrontMatter
	}
	if accept(l, Semicolon) {
		emit(l, lexeme.Semicolon)
		return LexFrontMatter
	}
	if accept(l, SymbolStart) {
		acceptRun(l, SymbolTail)

		switch l.input[l.start:l.pos] {
		case BoolType:
			fallthrough
		case NumberType:
			fallthrough
		case StringType:
			fallthrough
		case NullType:
			emit(l, lexeme.Type)
		case ExternKeyword:
			emit(l, lexeme.ExternKeyword)
		default:
			emit(l, lexeme.Symbol)
		}
		return LexFrontMatter
	}
	if strings.HasPrefix(l.input[l.pos:], Comment) {
		for l.pos < len(l.input) && !strings.HasPrefix(l.input[l.pos:], LineEnd) {
			l.pos++
		}
		ignore(l)
		if !accept(l, LineEnd) {
			return errorf(l, ErrorBadFrontMatterEnd)
		}
		ignore(l)
		return LexFrontMatter
	}
	return errorf(l, ErrorBadFrontmatter)
}

func LexLine(l *Lexer) State {
	if strings.HasPrefix(l.input[l.pos:], Hash) {
		return LexHeader
	}
	if strings.HasPrefix(l.input[l.pos:], FencedCodeBlockDelimiter) {
		return LexOpenCodeFence
	}
	if strings.HasPrefix(l.input[l.pos:], UnorderedListPrefix) {
		return LexUnorderedListItem
	}
	if strings.HasPrefix(l.input[l.pos:], OpenSquareBracket) {
		return LexLink
	}

	// ignore indents
	acceptRun(l, Whitespace)
	ignore(l)
	for {
		if strings.HasPrefix(l.input[l.pos:], CodeDelimiter) {
			if l.pos > l.start {
				emit(l, lexeme.TextLiteral)
			}
			return LexOpenInlineCode
		}
		if strings.HasPrefix(l.input[l.pos:], LineEnd) {
			if l.pos > l.start {
				emit(l, lexeme.TextLiteral)
			}
			accept(l, LineEnd)
			emit(l, lexeme.LineBreak)
			return LexLine
		}
		if _, err := next(l); err == io.EOF {
			break
		}
	}

	if l.pos > l.start {
		emit(l, lexeme.TextLiteral)
	}
	emit(l, lexeme.Eof)
	return nil
}

func LexHeader(l *Lexer) State {
	if acceptRun(l, Whitespace) {
		ignore(l)
		return LexHeader
	}
	if accept(l, Hash) {
		emit(l, lexeme.Hash)
		return LexHeader
	}
	if accept(l, SymbolStart) {
		acceptRun(l, SymbolTail)
		emit(l, lexeme.Symbol)
		return LexHeader
	}
	if accept(l, LineEnd) {
		emit(l, lexeme.LineBreak)
		return LexLine
	}
	return errorf(l, ErrorBadHeader)
}

func LexLink(l *Lexer) State {
	if acceptRun(l, Whitespace) {
		ignore(l)
		return LexLink
	}

	if accept(l, OpenSquareBracket) {
		emit(l, lexeme.OpenSquareBrace)
		return LexLink
	}
	if accept(l, SymbolStart) {
		acceptRun(l, SymbolTail)
		emit(l, lexeme.Symbol)
		return LexLink
	}
	if accept(l, CloseSquareBracket) {
		emit(l, lexeme.CloseSquareBrace)
		return LexLink
	}

	if accept(l, OpenParen) {
		emit(l, lexeme.OpenParen)
		acceptRun(l, Whitespace)
		ignore(l)

		for {
			r, err := peek(l)
			if err != nil {
				return errorf(l, ErrorBadLink)
			}
			if strings.ContainsRune(LineEnd, r) {
				return errorf(l, ErrorBadLink)
			}
			if strings.ContainsRune(CloseParen, r) {
				emit(l, lexeme.TextLiteral)
				accept(l, CloseParen)
				emit(l, lexeme.CloseParen)
				break
			}
			next(l)
		}
		return LexLink
	}

	if accept(l, LineEnd) {
		emit(l, lexeme.LineBreak)
		return LexLine
	}
	return errorf(l, ErrorBadLink)
}

func LexUnorderedListItem(l *Lexer) State {
	l.pos += len(UnorderedListPrefix)
	emit(l, lexeme.ListItemPrefix)

	acceptRun(l, Whitespace)
	ignore(l)

	return LexLink
}

func LexOpenInlineCode(l *Lexer) State {
	l.pos += len(CodeDelimiter)
	emit(l, lexeme.OpenInlineCode)
	return LexCode
}

func LexCloseInlineCode(l *Lexer) State {
	l.pos += len(CodeDelimiter)
	emit(l, lexeme.CloseInlineCode)
	return LexLine
}

func LexOpenCodeFence(l *Lexer) State {
	l.pos += len(FencedCodeBlockDelimiter)
	emit(l, lexeme.OpenCodeFence)
	if !accept(l, LineEnd) {
		return errorf(l, ErrorBadCode)
	}
	emit(l, lexeme.LineBreak)
	return LexCode
}

func LexCloseCodeFence(l *Lexer) State {
	l.pos += len(FencedCodeBlockDelimiter)
	emit(l, lexeme.CloseCodeFence)
	if !accept(l, LineEnd) {
		return errorf(l, ErrorBadCode)
	}
	emit(l, lexeme.LineBreak)
	return LexLine
}

func LexCode(l *Lexer) State {
	acceptRun(l, Whitespace+LineEnd)
	ignore(l)

	delimitersAndKeywords := []struct {
		prefix    string
		nextState State
	}{
		{LineEnd, LexLineEndInCode},
		{FencedCodeBlockDelimiter, LexCloseCodeFence},
		{CodeDelimiter, LexCloseInlineCode},
		{Comment, LexComment},
	}

	for _, d := range delimitersAndKeywords {
		if strings.HasPrefix(l.input[l.pos:], d.prefix) {
			return d.nextState
		}
	}

	codeTokens := []struct {
		tokenStart string
		nextState  State
	}{
		{SymbolStart, LexSymbol},
		{NumberStart, LexNumber},
		{StringStart, LexString},
	}

	for _, t := range codeTokens {
		nextRune, err := peek(l)
		if err != nil {
			continue
		}
		if strings.ContainsRune(t.tokenStart, nextRune) {
			return t.nextState
		}
	}

	{
		nextRune, err := peek(l)
		if err != nil {
			return errorf(l, ErrorBadCode)
		}
		if strings.ContainsRune(Operators, nextRune) {
			return LexOperator
		}
	}

	return errorf(l, ErrorBadCode)
}

func LexSymbol(l *Lexer) State {
	accept(l, SymbolStart)
	acceptRun(l, SymbolTail)
	switch l.input[l.start:l.pos] {
	case TrueLiteral:
		fallthrough
	case FalseLiteral:
		emit(l, lexeme.Boolean)
	case GotoLiteral:
		emit(l, lexeme.GotoLiteral)
	case NullLiteral:
		emit(l, lexeme.Null)
	case IfLiteral:
		emit(l, lexeme.IfLiteral)
	case ElseLiteral:
		emit(l, lexeme.ElseLiteral)
	case WhileLiteral:
		emit(l, lexeme.WhileLiteral)
	default:
		emit(l, lexeme.Symbol)
	}
	return LexCode
}

func LexNumber(l *Lexer) State {
	if accept(l, "0") {
		if accept(l, "0123456789eE") {
			return errorf(l, ErrorBadNumber)
		}
		emit(l, lexeme.Number)
		return LexCode
	}
	accept(l, "123456789")
	acceptRun(l, "0123456789")

	if accept(l, "eE") {
		accept(l, "+-")
		if !acceptRun(l, "0123456789") {
			return errorf(l, ErrorBadNumber)
		}
	}

	emit(l, lexeme.Number)
	return LexCode
}

func LexOperator(l *Lexer) State {
	operators := []struct {
		text      string
		tokenType lexeme.ItemType
	}{
		{PlusPlus, lexeme.Inc},
		{MinusMinus, lexeme.Dec},
		{And, lexeme.And},
		{Or, lexeme.Or},
		{EqEq, lexeme.DoubleEq},
		{Neq, lexeme.Neq},
		{Lte, lexeme.Lte},
		{Gte, lexeme.Gte},
		{OpenParen, lexeme.OpenParen},
		{CloseParen, lexeme.CloseParen},
		{OpenCurlyBrace, lexeme.OpenCurlyBrace},
		{CloseCurlyBrace, lexeme.CloseCurlyBrace},
		{Plus, lexeme.Plus},
		{Minus, lexeme.Minus},
		{Star, lexeme.Star},
		{Slash, lexeme.Slash},
		{Percent, lexeme.Percent},
		{Dot, lexeme.Dot},
		{Gt, lexeme.Gt},
		{Lt, lexeme.Lt},
		{Eq, lexeme.Eq},
		{Comma, lexeme.Comma},
		{Semicolon, lexeme.Semicolon},
		{Not, lexeme.Not},
	}

	for _, str := range operators {
		if strings.HasPrefix(l.input[l.pos:], str.text) {
			l.pos += len(str.text)
			emit(l, str.tokenType)
			return LexCode
		}
	}

	return errorf(l, ErrorBadOperator)
}

func LexComment(l *Lexer) State {
	for {
		if l.pos == len(l.input) {
			return errorf(l, ErrorBadCode)
		}
		if strings.HasPrefix(l.input[l.pos:], LineEnd) {
			break
		}
		l.pos++
	}
	ignore(l)
	return LexLineEndInCode
}

func LexLineEndInCode(l *Lexer) State {
	l.pos += len(LineEnd)
	ignore(l)
	return LexCode
}

func LexString(l *Lexer) State {
	accept(l, StringStart)
	for {
		if l.pos >= len(l.input) {
			return errorf(l, ErrorBadString)
		}
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if r < 0x20 {
			return errorf(l, ErrorBadString)
		}
		if r == '"' {
			l.pos += size
			break
		}
		if r == '\\' {
			l.pos += size
			if accept(l, "\\\"/bfnrt") {
				continue
			}
			if accept(l, "u") {
				isHex := accept(l, "0123456789abcdefABCDEF")
				isHex = isHex && accept(l, "0123456789abcdefABCDEF")
				isHex = isHex && accept(l, "0123456789abcdefABCDEF")
				isHex = isHex && accept(l, "0123456789abcdefABCDEF")
				if !isHex {
					return errorf(l, ErrorBadEscape)
				}
				continue
			}
			return errorf(l, ErrorBadEscape)
		}
		l.pos += size
	}
	emit(l, lexeme.String)
	return LexCode
}
