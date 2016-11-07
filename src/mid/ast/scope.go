package ast

import (
	"bytes"
	"fmt"

	"github.com/midlang/mid/src/mid/lexer"
)

// A Scope maintains the set of named language entities declared
// in the scope and a link to the immediately surrounding (outer)
// scope.
//
type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}

// NewScope creates a new scope nested in the outer scope.
func NewScope(outer *Scope) *Scope {
	const n = 4 // initial scope capacity
	return &Scope{outer, make(map[string]*Object, n)}
}

// Lookup returns the object with the given name if it is
// found in scope s, otherwise it returns nil. Outer scopes
// are ignored.
//
func (s *Scope) Lookup(name string) *Object {
	return s.Objects[name]
}

// Insert attempts to insert a named object obj into the scope s.
// If the scope already contains an object alt with the same name,
// Insert leaves the scope unchanged and returns alt. Otherwise
// it inserts obj and returns nil.
//
func (s *Scope) Insert(obj *Object) (alt *Object) {
	if alt = s.Objects[obj.Name]; alt == nil {
		s.Objects[obj.Name] = obj
	}
	return
}

// Debugging support
func (s *Scope) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "scope %p {", s)
	if s != nil && len(s.Objects) > 0 {
		fmt.Fprintln(&buf)
		for _, obj := range s.Objects {
			fmt.Fprintf(&buf, "\t%s %s\n", obj.Kind, obj.Name)
		}
	}
	fmt.Fprintf(&buf, "}\n")
	return buf.String()
}

type Object struct {
	Kind ObjKind
	Name string
	Decl interface{} // Field,Spec,Decl or Scope
	Data interface{}
}

// NewObj creates a new object of a given kind and name.
func NewObj(kind ObjKind, name string) *Object {
	return &Object{Kind: kind, Name: name}
}

func (obj *Object) Begin() lexer.Pos {
	name := obj.Name
	switch d := obj.Decl.(type) {
	case *Field:
		for _, n := range d.Names {
			if n.Name == name {
				return n.Begin()
			}
		}
	case *ImportSpec:
		if d.Name != nil && d.Name.Name == name {
			return d.Name.Begin()
		}
		return d.Package.Begin()
	case *ConstSpec:
		return d.Begin()
	case *GenDecl:
		return d.Begin()
	case *BeanDecl:
		return d.Begin()
	case *Scope:
		// nothing to do
	}
	return lexer.NoPos
}

type ObjKind int

const (
	Bad ObjKind = iota
	Pkg
	Const
	Var
	Bean
	Fun
)

var objKindStrings = [...]string{
	Bad:   "bad",
	Pkg:   "package",
	Const: "const",
	Var:   "var",
	Bean:  "bean",
	Fun:   "func",
}

func (kind ObjKind) String() string { return objKindStrings[kind] }
