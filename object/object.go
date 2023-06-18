package object

import (
	"bytes"
	"fmt"
	"github.com/GzzyZm/interpreter/ast"
	"strings"
)

type Type string

const (
	IntegerObj  = "INTEGER"
	BooleanObj  = "BOOLEAN"
	NullObj     = "NULL"
	ReturnObj   = "RETURN"
	ErrorObj    = "ERROR"
	FunctionObj = "FUNCTION"

	NullValue = "null"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return IntegerObj
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanObj
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NullObj
}
func (n *Null) Inspect() string {
	return NullValue
}

type Return struct {
	Value Object
}

func (r *Return) Type() Type {
	return ReturnObj
}
func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ErrorObj
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
	return FunctionObj
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
