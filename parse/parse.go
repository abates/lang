package parse

import "github.com/abates/lang/lex"

func Parse(lexer *lex.Lexer) {
	parsers := []parser{}
	for token := lexer.Next(); token.TokenType != lex.EOF; token = lexer.Next() {

	}
}
