package lexer

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/mcvoid/dialogue/internal/types/lexeme"
)

type Lexer struct {
	input string
	start int
	pos   int
	width int
	items []lexeme.Item
	state State
}

func New(input string) *Lexer {
	l := Lexer{
		input: input,
		items: []lexeme.Item{},
		state: LexFrontMatter,
	}
	return &l
}

func run(l *Lexer) {
	for l.state != nil {
		l.state = l.state(l)
	}
}

func emit(l *Lexer, t lexeme.ItemType) {
	substr := l.input[l.start:l.pos]
	item := lexeme.Item{
		Type: t,
		Val:  substr,
	}
	l.items = append(l.items, item)
	l.start = l.pos
}

func next(l *Lexer) (rune, error) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0, io.EOF
	}
	utf8Rune := rune(0)

	utf8Rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return utf8Rune, nil
}

func ignore(l *Lexer) {
	l.start = l.pos
}

func backup(l *Lexer) {
	l.pos -= l.width
}

func peek(l *Lexer) (rune, error) {
	r, err := next(l)
	backup(l)
	return r, err
}

func accept(l *Lexer, valid string) bool {
	r, _ := next(l)

	if strings.ContainsRune(valid, r) {
		return true
	}

	backup(l)
	return false
}

func acceptRun(l *Lexer, valid string) bool {
	start := l.pos
	r, err := next(l)

	for {
		if err == io.EOF {
			break
		}
		if !strings.ContainsRune(valid, r) {
			break
		}
		r, err = next(l)
	}

	backup(l)
	return l.pos > start
}

func errorf(l *Lexer, format string, args ...interface{}) State {
	l.items = append(l.items, lexeme.Item{
		Type: lexeme.Error,
		Val:  fmt.Sprintf(format, args...),
	})
	return nil
}

func (l *Lexer) Lex() []lexeme.Item {
	run(l)
	return l.items
}
