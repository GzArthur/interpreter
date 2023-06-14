package object

import (
	"bytes"
	"fmt"
	"github.com/GzArthur/interpreter/ast"
	"strings"
)

type Type string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"

	NULL_VALUE = "null"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return NULL_VALUE
}

type Return struct {
	Value Object
}

func (r *Return) Type() Type {
	return RETURN_OBJ
}
func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s" + e.Message)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var ps []string
	for _, p := range f.Parameters {
		ps = append(ps, p.PrintNode())
	}
	out.WriteString(fmt.Sprintf("fun(%s) {\n%s\n}", strings.Join(ps, ", "), f.Body.PrintNode()))
	return out.String()
}
