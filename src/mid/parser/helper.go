package parser

import (
	"bytes"
	"fmt"

	"github.com/midlang/mid/src/mid/ast"
	xscanner "github.com/midlang/mid/src/mid/external/go/scanner"
	"github.com/midlang/mid/src/mid/lexer"
	"github.com/midlang/mid/src/mid/scanner"
	"github.com/mkideal/pkg/errors"
)

type Error struct {
	Pos lexer.Position
	Msg string
}

func (err Error) Error() string { return fmt.Sprintf("%v: %s", err.Pos, err.Msg) }

func compareError(erri, errj error) bool {
	ei := erri.(*Error)
	ej := errj.(*Error)
	pi := ei.Pos
	pj := ej.Pos
	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	} else if pi.Line != pj.Line {
		return pi.Line < pj.Line
	} else {
		return pi.Column < pj.Column
	}
}

func (p *parser) init(fset *lexer.FileSet, filename string, src []byte) {
	p.file = fset.AddFile(filename, -1, len(src))
	eh := func(s *xscanner.Scanner, msg string) { p.errors.Add(&Error{s.Pos(), msg}) }
	p.scanner = scanner.NewScanner(p.file, bytes.NewReader(src))
	p.scanner.Error = eh
	p.comments = make([]*ast.CommentGroup, 0)
	p.errors = &errors.ErrorList{}

	p.next()
}

func (p *parser) openScope() {
	p.topScope = ast.NewScope(p.topScope)
}

func (p *parser) closeScope() {
	p.topScope = p.topScope.Outer
}

func (p *parser) declare(decl, data interface{}, scope *ast.Scope, kind ast.ObjKind, idents ...*ast.Ident) {
	for _, ident := range idents {
		assert(ident.Obj == nil, "identifier already declared or resolved")
		obj := ast.NewObj(kind, ident.Name)
		obj.Decl = decl
		obj.Data = data
		ident.Obj = obj
		if ident.Name != "_" {
			if alt := scope.Insert(obj); alt != nil {
				prevDecl := ""
				if pos := alt.Begin(); pos.IsValid() {
					prevDecl = fmt.Sprintf("\n\tprevious declaration at %v", p.file.Position(pos))
				}
				p.error(ident.Begin(), fmt.Sprintf("%s redeclared in this block%s", ident.Name, prevDecl))
			}
		}
	}
}

var unresolved = new(ast.Object)

func (p *parser) tryResolve(x ast.Expr, collectUnresolved bool) {
	if x == nil {
		return
	}
	switch n := x.(type) {
	case *ast.Ident:
		p.tryResolveIdent(n, collectUnresolved)
	case *ast.StructType:
		if n != nil {
			p.tryResolveIdent(n.Name, collectUnresolved)
		}
	case *ast.MapType:
		if n.K != nil {
			p.tryResolve(n.K, collectUnresolved)
		}
		if n.V != nil {
			p.tryResolve(n.V, collectUnresolved)
		}
	case *ast.ArrayType:
		if n.T != nil {
			p.tryResolve(n.T, collectUnresolved)
		}
	case *ast.VectorType:
		if n.T != nil {
			p.tryResolve(n.T, collectUnresolved)
		}
	}
}

func (p *parser) tryResolveIdent(ident *ast.Ident, collectUnresolved bool) {
	if ident == nil {
		return
	}
	assert(ident.Obj == nil, "identifier already declared or resolved")
	if ident.Name == "_" {
		return
	}
	// try to resolve the identifier
	for s := p.topScope; s != nil; s = s.Outer {
		if obj := s.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return
		}
	}
	if collectUnresolved {
		ident.Obj = unresolved
		p.unresolved = append(p.unresolved, ident)
	}
}

func (p *parser) resolve(x ast.Expr) {
	p.tryResolve(x, true)
}

func (p *parser) next0() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
	p.pos += lexer.Pos(p.file.Base()) - 1
}

func (p *parser) consumeComment() (comment *ast.Comment, endline int) {
	endline = p.file.Line(p.pos)
	if p.lit[1] == '*' {
		for i := 0; i < len(p.lit); i++ {
			if p.lit[i] == '\n' {
				endline++
			}
		}
	}

	comment = &ast.Comment{Slash: p.pos, Text: p.lit}
	p.next0()

	return
}

func (p *parser) consumeCommentGroup(n int) (comments *ast.CommentGroup, endline int) {
	var list []*ast.Comment
	endline = p.file.Line(p.pos)
	for p.tok == lexer.COMMENT && p.file.Line(p.pos) <= endline+n {
		var comment *ast.Comment
		comment, endline = p.consumeComment()
		list = append(list, comment)
	}

	// add comment group to the comments list
	comments = &ast.CommentGroup{List: list}
	p.comments = append(p.comments, comments)

	return
}

func (p *parser) next() {
	p.leadComment = nil
	p.lineComment = nil
	prev := p.pos
	p.next0()

	if p.tok == lexer.COMMENT {
		var comment *ast.CommentGroup
		var endline int

		if p.file.Line(p.pos) == p.file.Line(prev) {
			// The comment is on same line as the previous token; it
			// cannot be a lead comment but may be a line comment.
			comment, endline = p.consumeCommentGroup(0)
			if p.file.Line(p.pos) != endline {
				// The next token is on a different line, thus
				// the last comment group is a line comment.
				p.lineComment = comment
			}
		}

		// consume successor comments, if any
		endline = -1
		for p.tok == lexer.COMMENT {
			comment, endline = p.consumeCommentGroup(1)
		}

		if endline+1 == p.file.Line(p.pos) {
			// The next token is following on the line immediately after the
			// comment group, thus the last comment group is a lead comment.
			p.leadComment = comment
		}
	}
}

// A bailout panic is raised to indicate early termination.
type bailout struct{}

func (p *parser) error(pos lexer.Pos, msg string) {
	p.errors.Add(&Error{p.file.Position(pos), msg})
}

func (p *parser) errorExpected(pos lexer.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		if p.tok == lexer.SEMICOLON && p.lit == "\n" {
			msg += ", found newline"
		} else {
			msg += ", found '" + p.tok.String() + "'"
			if p.tok.IsLiteral() {
				msg += " " + p.lit
			}
		}
	}
	p.error(pos, msg)
}

func (p *parser) expect(tok lexer.Token) lexer.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *parser) expectOneOf(toks ...lexer.Token) lexer.Pos {
	pos := p.pos
	tok := lexer.ILLEGAL
	for _, want := range toks {
		if want == p.tok {
			tok = want
			break
		}
	}
	if tok == lexer.ILLEGAL {
		var buf bytes.Buffer
		for _, want := range toks {
			if buf.Len() > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(want.String())
		}
		p.errorExpected(pos, "'"+buf.String()+"'")
	}
	p.next() // make progress
	return pos
}

// expectClosing is like expect but provides a better error message
// for the common case of a missing comma before a newline.
//
func (p *parser) expectClosing(tok lexer.Token, context string) lexer.Pos {
	if p.tok != tok && p.tok == lexer.SEMICOLON && p.lit == "\n" {
		p.error(p.pos, "missing ',' before newline in "+context)
		p.next()
	}
	return p.expect(tok)
}

func (p *parser) expectSemi() {
	// semicolon is optional before a closing ')' or '}'
	if p.tok != lexer.RPAREN && p.tok != lexer.RBRACE {
		switch p.tok {
		case lexer.COMMA:
			// permit a ',' instead of a ';' but complain
			p.errorExpected(p.pos, "';'")
			fallthrough
		case lexer.SEMICOLON:
			p.next()
		}
	}
}

func (p *parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == lexer.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(lexer.IDENT) // use expect() error handling
	}
	return &ast.Ident{Pos: pos, Name: name}
}

func (p *parser) parseIdentList() (list []*ast.Ident) {
	list = append(list, p.parseIdent())
	for p.tok == lexer.COMMA {
		p.next()
		list = append(list, p.parseIdent())
	}
	return
}

func (p *parser) atComma(context string, follow lexer.Token) bool {
	if p.tok == lexer.COMMA {
		return true
	}
	if p.tok != follow {
		msg := "missing ','"
		if p.tok == lexer.SEMICOLON && p.lit == "\n" {
			msg += " before newline"
		}
		p.error(p.pos, msg+" in "+context)
		return true // "insert" comma and continue
	}
	return false
}

func assert(cond bool, msg string) {
	if !cond {
		panic("go/parser internal error: " + msg)
	}
}
