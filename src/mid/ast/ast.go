package ast

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/midlang/mid/src/mid/lexer"
)

type Node interface {
	Begin() lexer.Pos
}

// Node
// - Field,FieldList,Method,MethodList,Comment,CommentGroup,File,Package
// - Expr
//   - BadExpr,Ident,BasicLit
// - Type
//   - BasicType,ArrayType,MapType,VectorType,StructType
// - Decl
//   - GenDecl,BeanDecl
// - Spec
//   - ImportSpec,ConstSpec

//-----------
// expr node
//-----------
type Expr interface {
	Node
	exprNode()
}

func (*BadExpr) exprNode()    {}
func (*Ident) exprNode()      {}
func (*BasicLit) exprNode()   {}
func (*BasicType) exprNode()  {}
func (*StructType) exprNode() {}
func (*MapType) exprNode()    {}
func (*ArrayType) exprNode()  {}
func (*VectorType) exprNode() {}
func (*FuncType) exprNode()   {}

type BadExpr struct {
	From lexer.Pos
	To   lexer.Pos
}

func (be *BadExpr) Begin() lexer.Pos { return be.From }

// Ident node
type Ident struct {
	Pos  lexer.Pos
	Name string
	Obj  *Object
}

func (ident *Ident) Begin() lexer.Pos { return ident.Pos }

// basic literal node
type BasicLit struct {
	TokPos lexer.Pos
	Tok    lexer.Token
	Value  string
}

func (bl *BasicLit) Begin() lexer.Pos { return bl.TokPos }

func (bl *BasicLit) IsString() (bool, string) {
	if bl.Tok != lexer.STRING {
		return false, ""
	}
	s, _ := strconv.Unquote(bl.Value)
	return true, s
}

//-----------
// Type node
//-----------
type Type interface {
	Expr
	typeNode()
	Ident() *Ident
}

func (*BasicType) typeNode()  {}
func (*ArrayType) typeNode()  {}
func (*MapType) typeNode()    {}
func (*VectorType) typeNode() {}
func (*StructType) typeNode() {}
func (*FuncType) typeNode()   {}

func (t *BasicType) Ident() *Ident  { return t.Name }
func (t *ArrayType) Ident() *Ident  { return nil }
func (t *MapType) Ident() *Ident    { return nil }
func (t *VectorType) Ident() *Ident { return nil }
func (t *StructType) Ident() *Ident { return t.Name }
func (t *FuncType) Ident() *Ident   { return nil }

// builtin basic types(int,uint,string,bool,...)
type BasicType struct {
	Name *Ident
}

func (bt *BasicType) Begin() lexer.Pos { return bt.Name.Begin() }

// array<T,Size>
type ArrayType struct {
	Pos     lexer.Pos
	Less    lexer.Pos // <
	T       Type
	Size    Expr
	Greater lexer.Pos // >
}

func (at *ArrayType) Begin() lexer.Pos { return at.Pos }

// map<K,V>
type MapType struct {
	Pos     lexer.Pos
	Less    lexer.Pos // <
	K       Type
	V       Type
	Greater lexer.Pos // >
}

func (mt *MapType) Begin() lexer.Pos { return mt.Pos }

// vector<T>
type VectorType struct {
	Pos     lexer.Pos
	Less    lexer.Pos // <
	T       Type
	Greater lexer.Pos // >
}

func (vt *VectorType) Begin() lexer.Pos { return vt.Pos }

// struct/protocol
type StructType struct {
	Package *Ident // package or nil
	Name    *Ident
}

func (st *StructType) Begin() lexer.Pos {
	if st.Package != nil {
		return st.Package.Begin()
	}
	return st.Name.Begin()
}

type FuncType struct {
	Func   lexer.Pos
	Params *FieldList // arguments
	Result Type       // return type or nil
}

func (ft *FuncType) Begin() lexer.Pos {
	if ft.Func.IsValid() {
		return ft.Func
	}
	return ft.Params.Begin()
}

//---------------------
// Field and FieldList
//---------------------

// Field node
type Field struct {
	Doc     *CommentGroup // doc comment or nil
	Options []*Ident      // required/optional etc, maybe nil
	Type    Type          // Field type
	Names   []*Ident      // Field is a placeholder if len(Names) == 0
	Default Expr          // default node or nil
	Tag     *BasicLit     // tag or nil
	Comment *CommentGroup // line comment or nil
}

func (f *Field) Begin() lexer.Pos {
	if len(f.Options) > 0 {
		return f.Options[0].Begin()
	}
	return f.Type.Begin()
}

// FieldList node
type FieldList struct {
	Opening lexer.Pos // {
	List    []*Field
	Closing lexer.Pos // }
}

func (fl *FieldList) Begin() lexer.Pos { return fl.Opening }

//-----------
// Decl node
//-----------
type Decl interface {
	Node
	declNode()
}

func (*BadDecl) declNode()  {}
func (*GenDecl) declNode()  {}
func (*BeanDecl) declNode() {}

type BadDecl struct {
	From lexer.Pos
	To   lexer.Pos
}

func (bd *BadDecl) Begin() lexer.Pos { return bd.From }

// generic declaration node
type GenDecl struct {
	Doc    *CommentGroup // doc or nil
	TokPos lexer.Pos
	Tok    lexer.Token // import or const
	Lparen lexer.Pos   // (
	Specs  []Spec
	Rparen lexer.Pos // )
}

func (gd *GenDecl) Begin() lexer.Pos { return gd.TokPos }

// bean declaration node: struct or protocol
type BeanDecl struct {
	Kind    string // struct or protocol
	Pos     lexer.Pos
	Doc     *CommentGroup
	Name    *Ident
	Extends []Type
	Tag     *BasicLit
	Fields  *FieldList
}

func (bd *BeanDecl) Begin() lexer.Pos { return bd.Pos }

//-----------
// spec node
//-----------
type Spec interface {
	Node
	specNode()
}

func (*ImportSpec) specNode() {}
func (*ConstSpec) specNode()  {}

type ImportSpec struct {
	Doc     *CommentGroup // doc or nil
	Name    *Ident        // local name or nil
	Package *BasicLit     // package path
	Comment *CommentGroup // line comments or nil
}

func (is *ImportSpec) Begin() lexer.Pos {
	if is.Name != nil {
		return is.Name.Begin()
	}
	return is.Package.Begin()
}

type ConstSpec struct {
	Doc     *CommentGroup // doc or nil
	Name    *Ident        // name
	Value   Expr          // value node or nil
	Comment *CommentGroup // line comments or nil
}

func (cs *ConstSpec) Begin() lexer.Pos { return cs.Name.Begin() }

//--------------
// Generic Node
//--------------

type Comment struct {
	Slash lexer.Pos
	Text  string
	Values
}

func (c *Comment) Begin() lexer.Pos { return c.Slash }

type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Begin() lexer.Pos { return g.List[0].Begin() }

func (g *CommentGroup) Text() string {
	if g == nil || len(g.List) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range g.List {
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(c.Text)
	}
	return buf.String()
}

type Values map[string][]string

func (v Values) Get(key string) []string {
	values, ok := v[key]
	if !ok {
		return nil
	}
	return values
}

func (v Values) Exist(key string) bool {
	_, ok := v[key]
	return ok
}

// File node
type File struct {
	Filename   string
	Doc        *CommentGroup   // associated documentation; or nil
	Package    lexer.Pos       // position of "package" keyword
	Name       *Ident          // package name
	Decls      []Decl          // top-level declarations; or nil
	Scope      *Scope          // package scope (this file only)
	Imports    []*ImportSpec   // imports in this file
	Unresolved []*Ident        // unresolved identifiers in this file
	Comments   []*CommentGroup // list of all comments in the source file
}

func (f *File) Begin() lexer.Pos { return f.Package }

// Package node
type Package struct {
	Name    string             // package name
	Scope   *Scope             // package scope across all files
	Imports map[string]*Object // map of package id -> package object
	Files   map[string]*File   // source files by filename
}

func (p *Package) Begin() lexer.Pos { return lexer.NoPos }

type Visitor interface {
	Visit(Node) Visitor
	In()
	Out()
}

func walkNodes(visitor Visitor, nodes ...Node) Visitor {
	for _, node := range nodes {
		if visitor == nil {
			return visitor
		}
		Walk(node, visitor)
	}
	return visitor
}

func walkIdents(visitor Visitor, nodes []*Ident) Visitor {
	for _, node := range nodes {
		if visitor == nil {
			return visitor
		}
		Walk(node, visitor)
	}
	return visitor
}

func walkFields(visitor Visitor, nodes []*Field) Visitor {
	for _, node := range nodes {
		if visitor == nil {
			return visitor
		}
		Walk(node, visitor)
	}
	return visitor
}

func walkDecls(visitor Visitor, nodes []Decl) Visitor {
	for _, node := range nodes {
		if visitor == nil {
			return visitor
		}
		Walk(node, visitor)
	}
	return visitor
}

func walkSpecs(visitor Visitor, nodes []Spec) Visitor {
	for _, node := range nodes {
		if visitor == nil {
			return visitor
		}
		Walk(node, visitor)
	}
	return visitor
}

func Walk(node Node, visitor Visitor) Visitor {
	if node == nil || visitor == nil {
		return visitor
	}
	visitor = visitor.Visit(node)
	if visitor != nil {
		visitor.In()
	}
	defer func() {
		if visitor != nil {
			visitor.Out()
		}
	}()
	switch n := node.(type) {
	case *Field:
		visitor = walkNodes(visitor, n.Doc)
		visitor = walkIdents(visitor, n.Options)
		visitor = walkNodes(visitor, n.Type)
		visitor = walkIdents(visitor, n.Names)
		visitor = walkNodes(visitor, n.Default, n.Tag, n.Comment)
	case *FieldList:
		visitor = walkFields(visitor, n.List)
	case *Comment:
		return visitor
	case *CommentGroup:
		return visitor
	case *File:
		visitor = walkNodes(visitor, n.Doc, n.Name)
		visitor = walkDecls(visitor, n.Decls)
	case *Package:
		for _, file := range n.Files {
			if visitor != nil {
				visitor = Walk(file, visitor)
			}
		}
	case *Ident:
		return visitor
	case *BasicLit:
		return visitor
	case *BasicType:
		visitor = walkNodes(visitor, n.Name)
	case *ArrayType:
		visitor = walkNodes(visitor, n.T, n.Size)
	case *MapType:
		visitor = walkNodes(visitor, n.K, n.V)
	case *VectorType:
		visitor = walkNodes(visitor, n.T)
	case *StructType:
		visitor = walkNodes(visitor, n.Package, n.Name)
	case *FuncType:
		visitor = walkNodes(visitor, n.Params, n.Result)
	case *GenDecl:
		visitor = walkNodes(visitor, n.Doc)
		visitor = walkSpecs(visitor, n.Specs)
	case *BeanDecl:
		visitor = walkNodes(visitor, n.Doc, n.Name, n.Fields)
	case *ImportSpec:
		visitor = walkNodes(visitor, n.Doc, n.Name, n.Package, n.Comment)
	case *ConstSpec:
		visitor = walkNodes(visitor, n.Doc, n.Name, n.Value, n.Comment)
	case *BadDecl:
		return visitor
	case *BadExpr:
		return visitor
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
	return visitor
}
