package build

import (
	"encoding/gob"
	"strings"

	"github.com/midlang/mid/src/mid/ast"
	"github.com/midlang/mid/src/mid/lexer"
)

func init() {
	gob.Register(&Field{})
	gob.Register(&Bean{})
	gob.Register(&File{})
	gob.Register(&Package{})

	gob.Register(ExprBase{})
	gob.Register(TypeBase{})
	gob.Register(&BasicLit{})
	gob.Register(&BasicType{})
	gob.Register(&ArrayType{})
	gob.Register(&MapType{})
	gob.Register(&VectorType{})
	gob.Register(&GenDecl{})
	gob.Register(&ImportSpec{})
	gob.Register(&ConstSpec{})
	gob.Register(&StructType{})
	gob.Register(&FuncType{})
}

// helper functions
func BuildDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	text := doc.Text()
	if text != "" {
		text += "\n"
	}
	return text
}

func BuildComment(comment *ast.CommentGroup) string {
	if comment == nil {
		return ""
	}
	return comment.Text()
}

func BuildTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return tag.Value
}

func BuildIdent(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func BuildIdentList(idents []*ast.Ident) []string {
	if len(idents) == 0 {
		return []string{}
	}
	strs := make([]string, 0, len(idents))
	for _, ident := range idents {
		strs = append(strs, BuildIdent(ident))
	}
	return strs
}

type Field struct {
	Doc     string
	Options []string
	Type    Type
	Names   []string
	Default Expr
	Tag     string
	Comment string
}

func (field Field) NamesString() string {
	if len(field.Names) == 0 {
		return ""
	}
	return strings.Join(field.Names, ", ")
}

func BuildField(field *ast.Field) *Field {
	out := &Field{
		Doc:     BuildDoc(field.Doc),
		Options: BuildIdentList(field.Options),
		Type:    BuildType(field.Type),
		Names:   BuildIdentList(field.Names),
		Default: BuildExpr(field.Default),
		Tag:     BuildTag(field.Tag),
		Comment: BuildComment(field.Comment),
	}
	return out
}

func BuildFieldList(fields *ast.FieldList) []*Field {
	if fields == nil || len(fields.List) == 0 {
		return []*Field{}
	}
	list := make([]*Field, 0, len(fields.List))
	for _, field := range fields.List {
		list = append(list, BuildField(field))
	}
	return list
}

type Expr interface {
	ExprNode()
}

type ExprBase struct{}

func (ExprBase) ExprNode() {}
func (Ident) ExprNode()    {}
func (BasicLit) ExprNode() {}

type Ident string

type BasicLit struct {
	Kind  lexer.Token
	Value string
}

func BuildBasicLit(lit *ast.BasicLit) *BasicLit {
	return &BasicLit{
		Kind:  lit.Tok,
		Value: lit.Value,
	}
}

func BuildExpr(expr ast.Expr) Expr {
	if typ, ok := expr.(ast.Type); ok {
		return BuildType(typ)
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return Ident(BuildIdent(ident))
	}
	if lit, ok := expr.(*ast.BasicLit); ok {
		return BuildBasicLit(lit)
	}
	// TODO: alert error
	return &ExprBase{}
}

type Type interface {
	Expr
	TypeNode()
}

type TypeBase struct {
	ExprBase
}

func (TypeBase) TypeNode() {}

func BuildType(typ ast.Type) Type {
	switch t := typ.(type) {
	case *ast.BasicType:
		return &BasicType{Name: BuildIdent(t.Name)}
	case *ast.StructType:
		return BuildStruct(t)
	case *ast.MapType:
		return BuildMap(t)
	case *ast.ArrayType:
		return BuildArray(t)
	case *ast.VectorType:
		return BuildVector(t)
	case *ast.FuncType:
		return BuildFunc(t)
	default:
		//TODO: alert error
		return &TypeBase{}
	}
}

type BasicType struct {
	TypeBase
	Name string
}

type ArrayType struct {
	TypeBase
	T    Type
	Size Expr
}

func BuildArray(t *ast.ArrayType) *ArrayType {
	return &ArrayType{
		T:    BuildType(t.T),
		Size: BuildExpr(t.Size),
	}
}

type MapType struct {
	TypeBase
	K Type
	V Type
}

func BuildMap(t *ast.MapType) *MapType {
	return &MapType{
		K: BuildType(t.K),
		V: BuildType(t.V),
	}
}

type VectorType struct {
	TypeBase
	T Type
}

func BuildVector(t *ast.VectorType) *VectorType {
	return &VectorType{
		T: BuildType(t.T),
	}
}

type StructType struct {
	TypeBase
	Package string
	Name    string
}

func BuildStruct(t *ast.StructType) *StructType {
	return &StructType{
		Package: BuildIdent(t.Package),
		Name:    BuildIdent(t.Name),
	}
}

type FuncType struct {
	TypeBase
	Params []*Field
	Result Type
}

func BuildFunc(t *ast.FuncType) *FuncType {
	return &FuncType{
		Params: BuildFieldList(t.Params),
		Result: BuildType(t.Result),
	}
}

type Bean struct {
	Kind    string
	Doc     string
	Name    string
	Fields  []*Field
	Comment string
}

func BuildBean(bean *ast.BeanDecl) *Bean {
	return &Bean{
		Kind:   bean.Kind,
		Doc:    BuildDoc(bean.Doc),
		Name:   BuildIdent(bean.Name),
		Fields: BuildFieldList(bean.Fields),
	}
}

type ImportSpec struct {
	Doc     string
	Name    string
	Package string
	Comment string
}

func BuildImportSpec(spec *ast.ImportSpec) *ImportSpec {
	return &ImportSpec{
		Doc:     BuildDoc(spec.Doc),
		Name:    BuildIdent(spec.Name),
		Package: spec.Package.Value,
		Comment: BuildDoc(spec.Comment),
	}
}

type ConstSpec struct {
	Doc     string
	Name    string
	Value   Expr
	Comment string
}

func BuildConstSpec(spec *ast.ConstSpec) *ConstSpec {
	return &ConstSpec{
		Doc:     BuildDoc(spec.Doc),
		Name:    BuildIdent(spec.Name),
		Value:   BuildExpr(spec.Value),
		Comment: BuildComment(spec.Comment),
	}
}

type GenDecl struct {
	Doc     string
	Imports []*ImportSpec
	Consts  []*ConstSpec
}

func BuildGenDecl(decl *ast.GenDecl) *GenDecl {
	d := &GenDecl{
		Doc: BuildDoc(decl.Doc),
	}
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			d.Imports = append(d.Imports, BuildImportSpec(s))
		case *ast.ConstSpec:
			d.Consts = append(d.Consts, BuildConstSpec(s))
		default:
			//TODO: alert error
		}
	}
	return d
}

type File struct {
	Doc        string
	Package    string
	Beans      []*Bean
	Decls      []*GenDecl
	Unresolved []string
}

func BuildFile(file *ast.File) *File {
	f := &File{
		Doc:        BuildDoc(file.Doc),
		Package:    BuildIdent(file.Name),
		Unresolved: BuildIdentList(file.Unresolved),
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.BeanDecl:
			f.Beans = append(f.Beans, BuildBean(d))
		case *ast.GenDecl:
			f.Decls = append(f.Decls, BuildGenDecl(d))
		}
	}
	return f
}

type Package struct {
	Name    string
	Imports map[string]string
	Files   map[string]*File
}

func BuildPackage(pkg *ast.Package) *Package {
	p := &Package{
		Name:    pkg.Name,
		Imports: make(map[string]string),
		Files:   make(map[string]*File),
	}
	if pkg.Imports != nil {
		for id, name := range pkg.Imports {
			p.Imports[id] = name.Name
		}
	}
	if pkg.Files != nil {
		for name, file := range pkg.Files {
			p.Files[name] = BuildFile(file)
		}
	}
	return p
}
