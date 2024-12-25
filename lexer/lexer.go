package lexer

import (
	"github.com/ekediala/jian/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input [points to current char]
	readPosition int  // current reading position in input [after current char]
	ch           byte // current char under examination
}

func New(source string) *Lexer {
	l := Lexer{input: source}
	l.readChar()
	return &l
}

// The purpose of readChar is to give us the next character and advance our cursor in
// the source code
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	end := l.position
	return l.input[start:end]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	end := l.position
	return l.input[start:end]

}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case toByte(token.ASSIGN):
		{
			if next := l.peekChar(); next == toByte(token.ASSIGN) {
				tok.Type = token.EQ
				ch := l.ch
				l.readChar()
				tok.Literal = string(ch) + string(l.ch)
			} else {
				tok = newToken(token.ASSIGN, l.ch)
			}

		}
	case toByte(token.SEMICOLON):
		tok = newToken(token.SEMICOLON, l.ch)
	case toByte(token.LPAREN):
		tok = newToken(token.LPAREN, l.ch)
	case toByte(token.RPAREN):
		tok = newToken(token.RPAREN, l.ch)
	case toByte(token.RPAREN):
		tok = newToken(token.RPAREN, l.ch)
	case toByte(token.COMMA):
		tok = newToken(token.COMMA, l.ch)
	case toByte(token.PLUS):
		tok = newToken(token.PLUS, l.ch)
	case toByte(token.LBRACE):
		tok = newToken(token.LBRACE, l.ch)
	case toByte(token.RBRACE):
		tok = newToken(token.RBRACE, l.ch)
	case toByte(token.MINUS):
		tok = newToken(token.MINUS, l.ch)
	case toByte(token.BANG):
		{
			if next := l.peekChar(); next == toByte(token.ASSIGN) {
				tok.Type = token.EQ
				ch := l.ch
				l.readChar()
				tok.Literal = string(ch) + string(l.ch)
			} else {
				tok = newToken(token.BANG, l.ch)
			}
		}

	case toByte(token.ASTERISK):
		tok = newToken(token.ASTERISK, l.ch)
	case toByte(token.SLASH):
		tok = newToken(token.SLASH, l.ch)
	case toByte(token.LT):
		tok = newToken(token.LT, l.ch)
	case toByte(token.GT):
		tok = newToken(token.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// we return here because readIdentifier has already advanced the token
			return tok
		}

		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			// we return here because readNumber has already advanced the token
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}
	// advance cursor for next scan
	l.readChar()
	return tok
}
