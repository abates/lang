package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) TokenLiteral() string {
	return t.Literal
}

type StateFn func(*Lexer) StateFn

type Lexer struct {
	state  StateFn
	input  string
	start  int
	pos    int
	width  int
	tokens chan Token
}

func New(input string, startState StateFn) *Lexer {
	l := &Lexer{
		state:  startState,
		input:  input,
		tokens: make(chan Token, 2),
	}

	return l
}

func (l *Lexer) NextToken() Token {
	for {
		select {
		case t := <-l.tokens:
			return t
		default:
			if l.state == nil {
				return Token{EOF, ""}
			}

			l.state = l.state(l)
		}
	}
	// unreachable
}

func (l *Lexer) Next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

func (l *Lexer) IgnoreWhitespace() {
	for IsSpace(l.Next()) {
		l.Ignore()
	}
	l.Backup()
}

func (l *Lexer) Ignore() {
	l.start = l.pos
}

func (l *Lexer) Backup() {
	l.pos -= l.width
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptAlpha() {
	for IsAlpha(l.Next()) {
	}
	l.Backup()
}

func (l *Lexer) AcceptDigits() {
	for IsDigit(l.Next()) {
	}
	l.Backup()
}

func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

func (l *Lexer) AcceptSequence(sequence string) bool {
	if l.input[l.pos:l.pos+len(sequence)] == sequence {
		l.pos += len(sequence)
		return true
	}
	return false
}

func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{
		Type:    ILLEGAL,
		Literal: fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *Lexer) Emit(tokenType TokenType) {
	l.tokens <- Token{
		Type:    tokenType,
		Literal: l.CurrentInput(),
	}
	l.start = l.pos
}

func (l *Lexer) CurrentInput() string {
	return l.input[l.start:l.pos]
}

func IsEndOfLine(r rune) bool {
	return r == '\n' || r == '\r'
}

func IsSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsAlpha(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}
