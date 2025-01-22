package object

import (
	"strings"

	"github.com/ekediala/jian/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out strings.Builder
	params := make([]string, 0, len(f.Parameters))

	for _, param := range f.Parameters {
		params = append(params, param.Value)
	}

	out.WriteString("fn (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("){\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func NewFunction(params []*ast.Identifier, body *ast.BlockStatement, env *Environment) *Function {
	return &Function{params, body, env}
}
