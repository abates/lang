package lexer

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
		{IsNumber, '0', true},
		{IsNumber, '1', true},
		{IsNumber, '2', true},
		{IsNumber, '3', true},
		{IsNumber, '4', true},
		{IsNumber, '5', true},
		{IsNumber, '6', true},
		{IsNumber, '7', true},
		{IsNumber, '8', true},
		{IsNumber, '9', true},
		{IsNumber, 'a', false},
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

		if IsNumber(lexer.Peek()) {
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
		lexer.AcceptRun("0123456789")
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
	}

	token := lexer.NextToken()
	e := Token{EOF, ""}
	if token != e {
		t.Errorf("Expected %v got %v", e, token)
	}
}
