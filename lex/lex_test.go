package lex

import "testing"

func TestRuneTests(t *testing.T) {
	tests := []struct {
		test  func(rune) bool
		input rune
		want  bool
	}{
		{IsSpace, ' ', true},
		{IsSpace, '\t', true},
		{IsSpace, '\r', true},
		{IsSpace, '\n', true},
		{IsSpace, 'a', false},
		{IsDigit, '0', true},
		{IsDigit, '1', true},
		{IsDigit, '2', true},
		{IsDigit, '3', true},
		{IsDigit, '4', true},
		{IsDigit, '5', true},
		{IsDigit, '6', true},
		{IsDigit, '7', true},
		{IsDigit, '8', true},
		{IsDigit, '9', true},
		{IsDigit, 'a', false},
		{IsAlpha, 'a', true},
		{IsAlpha, 'z', true},
		{IsAlpha, 'A', true},
		{IsAlpha, 'Z', true},
		{IsAlpha, '_', true},
		{IsAlpha, '0', false},
	}

	for i, test := range tests {
		got := test.test(test.input)
		if test.want != got {
			t.Errorf("tests[%d] rune %c wanted %v got %v", i, test.input, test.want, got)
		}
	}
}

func TestLexing(t *testing.T) {
	word := TokenType("WORD")
	number := TokenType("NUMBER")
	punctuation := TokenType("PUNCTUATION")

	var lex StateFn
	var lexWord StateFn
	var lexNumber StateFn
	var lexPunctuation StateFn

	lex = func(lexer *Lexer) StateFn {
		if lexer.Peek() == -1 {
			return nil
		}

		lexer.IgnoreWhitespace()
		if IsAlpha(lexer.Peek()) {
			return lexWord
		}

		if IsDigit(lexer.Peek()) {
			return lexNumber
		}

		return lexPunctuation
	}

	lexWord = func(lexer *Lexer) StateFn {
		lexer.AcceptAlpha()
		lexer.Emit(word)
		return lex
	}

	lexNumber = func(lexer *Lexer) StateFn {
		lexer.AcceptDigits()
		lexer.Emit(number)
		return lex
	}

	lexPunctuation = func(lexer *Lexer) StateFn {
		if !lexer.Accept(":") {
			lexer.Next()
			lexer.Errorf("bad punctuation")
		} else {
			lexer.Emit(punctuation)
		}
		return lex
	}

	input := "this is a string of words and numbers: 1234567890abc."
	expected := []Token{
		{word, "this"},
		{word, "is"},
		{word, "a"},
		{word, "string"},
		{word, "of"},
		{word, "words"},
		{word, "and"},
		{word, "numbers"},
		{punctuation, ":"},
		{number, "1234567890"},
		{word, "abc"},
		{ILLEGAL, "bad punctuation"},
	}

	lexer := New(input, lex)
	for _, e := range expected {
		token := lexer.NextToken()
		if token != e {
			t.Errorf("Expected %v got %v", e, token)
		}

		if token.TokenLiteral() != e.Literal {
			t.Errorf("Expected %v got %v", e.Literal, token.TokenLiteral())
		}
	}

	token := lexer.NextToken()
	e := Token{EOF, ""}
	if token != e {
		t.Errorf("Expected %v got %v", e, token)
	}
}

func TestAcceptRun(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"aabbcc", "cc"},
	}

	for i, test := range tests {
		lexer := New(test.input, nil)
		lexer.AcceptRun("ab")
		if lexer.input[lexer.pos:] != test.expected {
			t.Errorf("tests[%d] expected %q got %q", i, test.expected, lexer.input[lexer.pos:])
		}
	}
}

func TestAcceptSequence(t *testing.T) {
	tests := []struct {
		input    string
		sequence string
		expected string
	}{
		{"aabbcc", "abc", "aabbcc"},
		{"abcc", "abc", "c"},
		{"foo ©2017", "foo ©", "2017"},
	}

	for i, test := range tests {
		lexer := New(test.input, nil)
		lexer.AcceptSequence(test.sequence)
		if lexer.input[lexer.pos:] != test.expected {
			t.Errorf("tests[%d] expected %q got %q", i, test.expected, lexer.input[lexer.pos:])
		}
	}
}
