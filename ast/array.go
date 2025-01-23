package ast

import (
	"strings"

	"github.com/ekediala/jian/token"
)

type ArrayLiteral struct {
	Token    token.Token // [ token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // the [ Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var s strings.Builder

	s.WriteByte('(')
	s.WriteString(ie.Left.String())
	s.WriteByte('[')
	s.WriteString(ie.Index.String())
	s.WriteByte(']')
	s.WriteByte(')')

	return s.String()
}
