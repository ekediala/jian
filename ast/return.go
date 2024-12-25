package ast

import (
	"strings"

	"github.com/ekediala/jian/token"
)

type ReturnStatement struct {
	Token       token.Token //r eturn token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(rs.ReturnValue.String())
	out.WriteString(";")

	return out.String()
}
