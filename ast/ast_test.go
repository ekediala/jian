package ast_test

import (
	"testing"

	"github.com/ekediala/jian/ast"
	"github.com/ekediala/jian/token"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Literal: "let",
					Type:    token.LET,
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Literal: "myVar",
						Type:    token.IDENT,
					},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Literal: "anotherVar",
						Type:    token.IDENT,
					},
					Value: "anotherVar",
				},
			},
		},
	}

	if got, expected := program.String(), "let myVar = anotherVar;"; got != expected {
		t.Errorf("expected program.String() to produce %q, got %s", expected, got)
	}
}
