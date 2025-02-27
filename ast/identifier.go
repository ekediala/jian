package ast

import (
	"github.com/ekediala/jian/token"
)

type Identifier struct {
	Token token.Token // token.IDENT
	Value  string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String () string {
	return i.Value
}
