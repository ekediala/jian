package lexer

import (
	"unicode"

	"github.com/ekediala/jian/token"
)

func newToken(t token.TokenType, literal byte) token.Token {
	return token.Token{Type: t, Literal: string(literal)}
}

func toByte(t token.TokenType) byte {
	return byte(t[0])
}

// checks if the current byte is a valid identifier character
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
