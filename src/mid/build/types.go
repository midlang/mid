package build

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strconv"
	"strings"

	"github.com/midlang/mid/src/mid/ast"
	"github.com/midlang/mid/src/mid/lexer"
	"github.com/mkideal/log"
)

var (
	ErrAmbiguousNames = errors.New("ambiguous names")
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
	text := comment.Text()
	return text
}

func BuildTag(tag *ast.BasicLit) Tag {
	if tag == nil {
		return ""
	}
	s := tag.Value
	const q = "`"
	if strings.HasPrefix(s, q) {
		s = strings.TrimPrefix(s, q)
	}
	if strings.HasSuffix(s, q) {
		s = strings.TrimSuffix(s, q)
	}
	return Tag(s)
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

type Tag string

func (tag Tag) Clone() *Tag {
	x := Tag(string(tag))
	return &x
}
func (tag Tag) String() string { return string(tag) }

func (tag Tag) HasKey(key string) bool {
	_, ok := tag.Lookup(key)
	return ok
}

func (tag Tag) Get(key string) string {
	value, _ := tag.Lookup(key)
	return value
}

func (tag *Tag) Set(key, value string) string {
	pairs, _, index := tag.parse(key)
	if index >= 0 {
		pairs[index][1] = value
	} else {
		pairs = append(pairs, tagpair{key, value})
	}
	*tag = Tag(tag.format(pairs))
	return ""
}

func (tag *Tag) Del(key string) string {
	pairs, _, index := tag.parse(key)
	if index >= 0 {
		pairs = append(pairs[:index], pairs[index+1:]...)
	}
	*tag = Tag(tag.format(pairs))
	return ""
}

func (tag Tag) Lookup(key string) (value string, ok bool) {
	_, value, index := tag.parse(key)
	return value, index >= 0
}

type tagpair [2]string

func (tag Tag) format(pairs []tagpair) string {
	var buf bytes.Buffer
	for i, pair := range pairs {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(pair[0])
		buf.WriteByte(':')
		buf.WriteByte('"')
		buf.WriteString(pair[1])
		buf.WriteByte('"')
	}
	return buf.String()
}

func (tag Tag) parse(key string) (pairs []tagpair, value string, index int) {
	pairs = make([]tagpair, 0)
	index = -1
	count := 0
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		v, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
		}
		if key == name && key != "" && index == -1 {
			value = v
			index = len(pairs)
		}
		pairs = append(pairs, tagpair{name, v})
		count++
	}
	return
}

type Field struct {
	Doc     string
	Options []string
	Type    Type
	Names   []string
	Default Expr
	Tag     Tag
	Comment string
}

func (field Field) Name() (string, error) {
	if len(field.Names) == 0 {
		return "", nil
	}
	if len(field.Names) == 1 {
		return field.Names[0], nil
	}
	return "", ErrAmbiguousNames
}

func (field Field) Value() string {
	switch e := field.Default.(type) {
	case *BasicLit:
		return e.Value
	}
	panic("unsupported expr")
}

func (field Field) GetTag(key string) string {
	return field.Tag.Get(key)
}

func (field Field) HasTag(key string) bool {
	_, ok := field.Tag.Lookup(key)
	return ok
}

func (field *Field) AddTag(key, value string) {
	field.Tag.Set(key, value)
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
	log.Trace("BuildField: field=%v", out.Names)
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

func IntFromExpr(expr Expr) (string, bool) {
	switch e := expr.(type) {
	case *BasicLit:
		if e.Kind == lexer.INT {
			return e.Value, true
		}
	}
	return "", false
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
	IsArray() bool
	IsVector() bool
	IsMap() bool
	IsStruct() bool
	IsString() bool
	IsInt() bool
	IsBool() bool
}

type TypeBase struct {
	ExprBase
}

func (TypeBase) TypeNode()      {}
func (TypeBase) IsArray() bool  { return false }
func (TypeBase) IsVector() bool { return false }
func (TypeBase) IsMap() bool    { return false }
func (TypeBase) IsStruct() bool { return false }
func (TypeBase) IsString() bool { return false }
func (TypeBase) IsInt() bool    { return false }
func (TypeBase) IsBool() bool   { return false }

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
		return &TypeBase{}
	}
}

type BasicType struct {
	TypeBase
	Name string
}

func (t BasicType) IsVector() bool { return t.Name == lexer.Bytes.String() }
func (t BasicType) IsString() bool { return t.Name == lexer.String.String() }
func (t BasicType) IsBool() bool   { return t.Name == lexer.Bool.String() }
func (t BasicType) IsInt() bool {
	bt, ok := lexer.LookupType(t.Name)
	if !ok {
		return false
	}
	return bt.IsInt()
}

type ArrayType struct {
	TypeBase
	T    Type
	Size Expr
}

func (ArrayType) IsArray() bool { return true }

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

func (MapType) IsMap() bool { return true }

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

func (VectorType) IsVector() bool { return true }

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

func (StructType) IsStruct() bool { return true }

func BuildStruct(t *ast.StructType) *StructType {
	return &StructType{
		Package: BuildIdent(t.Package),
		Name:    BuildIdent(t.Name),
	}
}

func (t StructType) String(sep string) string {
	if t.Package == "" {
		return t.Name
	}
	return t.Package + sep + t.Name
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
	Extends []Type
	Tag     Tag
	Fields  []*Field
	Comment string
}

func (bean *Bean) IsNil() bool { return bean == nil }

func BuildBean(bean *ast.BeanDecl) *Bean {
	b := &Bean{
		Kind:   bean.Kind,
		Doc:    BuildDoc(bean.Doc),
		Name:   BuildIdent(bean.Name),
		Tag:    BuildTag(bean.Tag),
		Fields: BuildFieldList(bean.Fields),
	}
	if len(bean.Extends) > 0 {
		b.Extends = make([]Type, 0, len(bean.Extends))
		for _, e := range bean.Extends {
			b.Extends = append(b.Extends, BuildType(e))
		}
	}
	return b
}

func (bean Bean) Field(i int) *Field {
	if i >= len(bean.Fields) || i < 0 {
		return nil
	}
	return bean.Fields[i]
}

func (bean Bean) NumField() int { return len(bean.Fields) }

func (bean Bean) GetTag(key string) string {
	return bean.Tag.Get(key)
}

func (bean Bean) HasTag(key string) bool {
	_, ok := bean.Tag.Lookup(key)
	return ok
}

func (bean *Bean) AddTag(key, value string) {
	bean.Tag.Set(key, value)
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

func (c ConstSpec) ValueString() string {
	switch e := c.Value.(type) {
	case *BasicLit:
		return e.Value
	}
	panic("unsupported expr")
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
	Filename   string
	Doc        string
	Package    string
	Beans      []*Bean
	Decls      []*GenDecl
	Unresolved []string
}

func BuildFile(file *ast.File) *File {
	f := &File{
		Filename:   file.Filename,
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
