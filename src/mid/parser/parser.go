package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/midlang/mid/src/mid/ast"
	"github.com/midlang/mid/src/mid/lexer"
	"github.com/midlang/mid/src/mid/scanner"
	"github.com/mkideal/pkg/errors"
)

type parser struct {
	scanner *scanner.Scanner
	file    *lexer.File
	errors  *errors.ErrorList
	mode    uint

	pos lexer.Pos
	tok lexer.Token
	lit string

	syncPos lexer.Pos
	syncCnt int

	comments    []*ast.CommentGroup
	leadComment *ast.CommentGroup
	lineComment *ast.CommentGroup

	imports []*ast.ImportSpec

	unresolved []*ast.Ident

	topScope *ast.Scope
	pkgScope *ast.Scope

	inRhs bool
}

func (p *parser) parseFile() *ast.File {
	if p.errors.Len() > 0 {
		return nil
	}
	var (
		doc   = p.leadComment
		pos   = p.expect(lexer.PACKAGE)
		ident = p.parseIdent()
		decls []ast.Decl
	)
	p.expectSemi()

	if p.errors.Len() > 0 {
		return nil
	}

	p.openScope()
	p.pkgScope = p.topScope

	// parse imports
	for p.tok == lexer.IMPORT {
		decls = append(decls, p.parseGenDecl(lexer.IMPORT, p.parseImportSpec))
	}

	// parse body
	for p.tok != lexer.EOF {
		decls = append(decls, p.parseDecl(syncDecl))
	}

	p.closeScope()
	assert(p.topScope == nil, "unbalanced scopes")

	i := 0
	for _, ident := range p.unresolved {
		assert(ident.Obj == unresolved, "object already resolved")
		ident.Obj = p.pkgScope.Lookup(ident.Name)
		if ident.Obj == nil {
			p.unresolved[i] = ident
			i++
		}
	}

	return &ast.File{
		Doc:        doc,
		Package:    pos,
		Name:       ident,
		Decls:      decls,
		Scope:      p.pkgScope,
		Imports:    p.imports,
		Unresolved: p.unresolved[0:i],
		Comments:   p.comments,
	}
}

func syncDecl(p *parser) {
	for {
		switch p.tok {
		case lexer.CONST, lexer.PROTOCOL, lexer.STRUCT, lexer.SERVICE, lexer.ENUM:
			// see comments in syncStmt
			if p.pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
				return
			}
			if p.pos > p.syncPos {
				p.syncPos = p.pos
				p.syncCnt = 0
				return
			}
		case lexer.EOF:
			return
		}
		p.next()
	}
}

func (p *parser) parseGenDecl(keyword lexer.Token, specFunc parseSpecFunction) *ast.GenDecl {
	var (
		doc            = p.leadComment
		pos            = p.expect(keyword)
		lparen, rparen lexer.Pos
		list           []ast.Spec
	)
	if p.tok == lexer.LPAREN {
		lparen = p.pos
		p.next()
		for iota := 0; p.tok != lexer.RPAREN && p.tok != lexer.EOF; iota++ {
			list = append(list, specFunc(p.leadComment, keyword, iota))
		}
		rparen = p.expect(lexer.RPAREN)
		//p.expectSemi()
	} else {
		list = append(list, specFunc(nil, keyword, 0))
	}
	return &ast.GenDecl{
		Doc:    doc,
		TokPos: pos,
		Tok:    keyword,
		Lparen: lparen,
		Specs:  list,
		Rparen: rparen,
	}
}

func (p *parser) parseDecl(sync func(*parser)) ast.Decl {
	var specFunc parseSpecFunction
	switch p.tok {
	case lexer.CONST:
		specFunc = p.parseValueSpec
	case lexer.STRUCT, lexer.PROTOCOL, lexer.SERVICE, lexer.ENUM:
		return p.parseBeanDecl(p.topScope)
	default:
		pos := p.pos
		p.errorExpected(pos, "declaration")
		sync(p)
		return &ast.BadDecl{From: pos, To: p.pos}
	}
	return p.parseGenDecl(p.tok, specFunc)
}

func (p *parser) parseBeanDecl(parentScope *ast.Scope) ast.Decl {
	var (
		tag   *ast.BasicLit
		doc   = p.leadComment
		tok   = p.tok
		pos   = p.expectOneOf(lexer.PROTOCOL, lexer.STRUCT, lexer.SERVICE, lexer.ENUM)
		ident = p.parseIdent()
	)
	if p.tok == lexer.STRING {
		tag = &ast.BasicLit{TokPos: p.pos, Tok: p.tok, Value: p.lit}
		p.next()
	}
	scope := ast.NewScope(parentScope)
	lbrace := p.expect(lexer.LBRACE)
	var list []*ast.Field
	switch tok {
	case lexer.SERVICE:
		for p.tok == lexer.IDENT {
			list = append(list, p.parseMethodSpec(scope))
		}
	case lexer.ENUM:
		for p.tok == lexer.IDENT {
			list = append(list, p.parseEnumSpec(scope))
		}
	default:
		for p.tok == lexer.IDENT || p.tok == lexer.EXTEND || p.tok == lexer.LPAREN {
			list = append(list, p.parseFieldDecl(scope))
		}
	}
	rbrace := p.expect(lexer.RBRACE)
	spec := &ast.BeanDecl{
		Kind: tok.String(),
		Pos:  pos,
		Doc:  doc,
		Name: ident,
		Tag:  tag,
		Fields: &ast.FieldList{
			Opening: lbrace,
			List:    list,
			Closing: rbrace,
		},
	}
	p.declare(spec, nil, parentScope, ast.Bean, ident)
	return spec
}

func (p *parser) parseFieldDecl(scope *ast.Scope) *ast.Field {
	var (
		doc     = p.leadComment
		options = p.parseFieldOptions()
		typ     ast.Type
		idents  []*ast.Ident
		tag     *ast.BasicLit
	)
	if p.tok == lexer.EXTEND {
		p.next()
		typ = p.parseTypeName()
	} else {
		typ = p.parseTypeName()
		idents = p.parseIdentList()
	}
	if p.tok == lexer.STRING {
		tag = &ast.BasicLit{TokPos: p.pos, Tok: p.tok, Value: p.lit}
		p.next()
	}
	p.expectSemi()
	field := &ast.Field{
		Doc:     doc,
		Options: options,
		Type:    typ,
		Names:   idents,
		Tag:     tag,
		Comment: p.lineComment,
	}
	p.declare(field, nil, scope, ast.Var, idents...)
	p.resolve(typ)
	return field
}

func (p *parser) parseMethodSpec(scope *ast.Scope) *ast.Field {
	var (
		doc    = p.leadComment
		typ    ast.Type
		idents []*ast.Ident
	)
	x := p.parseTypeName()
	if ident := x.Ident(); ident != nil && p.tok == lexer.LPAREN {
		idents = append(idents, ident)
		scope := ast.NewScope(nil)
		params, result := p.parseSignature(scope)
		typ = &ast.FuncType{
			Func:   lexer.NoPos,
			Params: params,
			Result: result,
		}
	} else {
		typ = x
		p.resolve(typ)
	}
	spec := &ast.Field{
		Doc:     doc,
		Names:   idents,
		Type:    typ,
		Comment: p.lineComment,
	}
	p.declare(spec, nil, scope, ast.Fun, idents...)
	return spec
}

func (p *parser) parseEnumSpec(scope *ast.Scope) *ast.Field {
	doc := p.leadComment
	name := p.parseIdent()
	p.expect(lexer.ASSIGN)
	var value ast.Expr
	if p.tok == lexer.INT {
		value = &ast.BasicLit{
			TokPos: p.pos,
			Tok:    p.tok,
			Value:  p.lit,
		}
		p.next()
	} else {
		value = p.parseIdent()
	}
	p.expect(lexer.COMMA)
	spec := &ast.Field{
		Doc:     doc,
		Names:   []*ast.Ident{name},
		Default: value,
		Comment: p.lineComment,
	}
	return spec
}

func (p *parser) parseSignature(scope *ast.Scope) (*ast.FieldList, ast.Type) {
	pos := p.pos
	params := p.parseParameters(scope)
	curPos := p.pos
	if p.file.Position(pos).Line != p.file.Position(curPos).Line {
		return params, nil
	}
	result := p.parseTypeName()
	return params, result
}

func (p *parser) parseParameters(scope *ast.Scope) *ast.FieldList {
	var params []*ast.Field
	lparen := p.expect(lexer.LPAREN)
	if p.tok != lexer.RPAREN {
		params = p.parseParameterList(scope)
	}
	rparen := p.expect(lexer.RPAREN)
	return &ast.FieldList{
		Opening: lparen,
		List:    params,
		Closing: rparen,
	}
}

func (p *parser) parseParameterList(scope *ast.Scope) []*ast.Field {
	var list []*ast.Field
	for p.tok != lexer.RPAREN {
		typ := p.parseTypeName()
		var idents []*ast.Ident
		tok := p.tok
		switch tok {
		case lexer.IDENT:
			idents = []*ast.Ident{p.parseIdent()}
			if p.tok == lexer.COMMA {
				p.next()
			}
		case lexer.COMMA:
			p.next()
		case lexer.RPAREN:
		default:
			p.error(p.pos, "unexpected token "+p.tok.String()+" after type")
			return list
		}
		list = append(list, &ast.Field{
			Type:  typ,
			Names: idents,
		})
	}
	return list
}

func (p *parser) parseFieldOptions() []*ast.Ident {
	// TODO
	return nil
}

func (p *parser) parseTypeName() ast.Type {
	pos := p.pos
	ident := p.parseIdent()
	if p.tok == lexer.PERIOD {
		p.next()
		p.resolve(ident)
		name := p.parseIdent()
		return &ast.StructType{
			Package: ident,
			Name:    name,
		}
	}
	bt, ok := lexer.LookupType(ident.Name)
	if ok {
		switch bt {
		case lexer.Map:
			lessPos := p.pos
			p.expect(lexer.LESS)
			k := p.parseTypeName()
			p.expect(lexer.COMMA)
			v := p.parseTypeName()
			greaterPos := p.pos
			p.expect(lexer.GREATER)
			return &ast.MapType{
				Pos:     pos,
				Less:    lessPos,
				K:       k,
				V:       v,
				Greater: greaterPos,
			}
		case lexer.Array:
			lessPos := p.pos
			p.expect(lexer.LESS)
			t := p.parseTypeName()
			p.expect(lexer.COMMA)
			sizePos := p.pos
			var size ast.Expr
			if p.tok == lexer.INT {
				size = &ast.BasicLit{
					TokPos: sizePos,
					Tok:    lexer.INT,
					Value:  p.lit,
				}
				p.next()
			} else {
				size = p.parseIdent()
			}
			greaterPos := p.pos
			p.expect(lexer.GREATER)
			return &ast.ArrayType{
				Pos:     pos,
				Less:    lessPos,
				T:       t,
				Size:    size,
				Greater: greaterPos,
			}
		case lexer.Vector:
			lessPos := p.pos
			p.expect(lexer.LESS)
			t := p.parseTypeName()
			greaterPos := p.pos
			p.expect(lexer.GREATER)
			return &ast.VectorType{
				Pos:     pos,
				Less:    lessPos,
				T:       t,
				Greater: greaterPos,
			}
		default:
			return &ast.BasicType{Name: ident}
		}
	} else {
		return &ast.StructType{
			Name: ident,
		}
	}
}

type parseSpecFunction func(doc *ast.CommentGroup, keyword lexer.Token, iota int) ast.Spec

func isValidImport(lit string) bool {
	const illegalChars = `!"#$%&'()*,:;<=>?[\]^{|}` + "`\uFFFD"
	s, _ := strconv.Unquote(lit) // go/scanner returns a legal string literal
	for _, r := range s {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
			return false
		}
	}
	return s != ""
}

func (p *parser) parseValueSpec(doc *ast.CommentGroup, keyword lexer.Token, iota int) ast.Spec {
	pos := p.pos
	ident := p.parseIdent()
	var value ast.Expr
	if p.tok == lexer.ASSIGN {
		p.next()
		switch p.tok {
		case lexer.INT, lexer.FLOAT, lexer.STRING:
			value = &ast.BasicLit{TokPos: p.pos, Tok: p.tok, Value: p.lit}
			p.next()
		default:
			value = p.parseIdent()
		}
	}
	p.expectSemi()
	switch keyword {
	case lexer.CONST:
		if value == nil && iota == 0 {
			p.error(pos, "missing constant value")
		}
	}
	spec := &ast.ConstSpec{
		Doc:     doc,
		Name:    ident,
		Value:   value,
		Comment: p.lineComment,
	}
	kind := ast.Const
	p.declare(spec, iota, p.topScope, kind, ident)
	return spec
}

func (p *parser) parseImportSpec(doc *ast.CommentGroup, _ lexer.Token, _ int) ast.Spec {
	var ident *ast.Ident
	switch p.tok {
	case lexer.PERIOD:
		ident = &ast.Ident{Pos: p.pos, Name: "."}
		p.next()
	case lexer.IDENT:
		ident = p.parseIdent()
	}
	pos := p.pos
	var path string
	if p.tok == lexer.STRING {
		path = p.lit
		if !isValidImport(path) {
			p.error(pos, "invalid import path: "+path)
		}
		p.next()
	} else {
		p.expect(lexer.STRING)
	}
	p.expectSemi()

	spec := &ast.ImportSpec{
		Doc:     doc,
		Name:    ident,
		Package: &ast.BasicLit{TokPos: pos, Tok: lexer.STRING, Value: path},
	}

	p.imports = append(p.imports, spec)
	return spec
}

func ParseFile(fset *lexer.FileSet, filename string, src []byte) (f *ast.File, err error) {
	if len(src) == 0 {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
	}
	p := new(parser)
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}
		if f == nil {
			f = &ast.File{
				Name:  new(ast.Ident),
				Scope: ast.NewScope(nil),
			}
		}
		err = p.errors.Err()
	}()
	p.init(fset, filename, src)
	f = p.parseFile()
	err = p.errors.Err()
	return
}

func ParseFiles(fset *lexer.FileSet, files []string) (map[string]*ast.Package, error) {
	pkgs := make(map[string]*ast.Package)
	var firstErr error
	for _, filename := range files {
		if f, err := ParseFile(fset, filename, nil); err == nil {
			name := f.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:    name,
					Scope:   ast.NewScope(nil),
					Imports: make(map[string]*ast.Object),
					Files:   make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[name] = f
			if f.Scope != nil && f.Scope.Objects != nil {
				errors := &errors.ErrorList{}
				for _, obj := range f.Scope.Objects {
					if alt := pkg.Scope.Insert(obj); alt != nil {
						prevDecl := ""
						if pos := alt.Begin(); pos.IsValid() {
							prevDecl = fmt.Sprintf("\n\tprevious declaration at %v", fset.Position(pos))
						}
						pos := fset.Position(obj.Begin())
						msg := fmt.Sprintf("%s redeclared in this block%s", obj.Name, prevDecl)
						errors.Add(&Error{pos, msg})
					}
				}
				if firstErr == nil && errors.Len() > 0 {
					errors.Sort(compareError)
					firstErr = errors.Err()
				}
			}
		} else if firstErr == nil {
			firstErr = err
		}
	}

	return pkgs, firstErr
}

func ParseDir(fset *lexer.FileSet, dir string, suffix string, filter func(os.FileInfo) bool) (map[string]*ast.Package, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	list, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, d := range list {
		if !strings.HasSuffix(d.Name(), suffix) || (filter != nil && !filter(d)) {
			continue
		}
		filename := filepath.Join(dir, d.Name())
		files = append(files, filename)
	}

	return ParseFiles(fset, files)
}
