package parser

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/midlang/mid/src/mid/ast"
	"github.com/midlang/mid/src/mid/lexer"
)

type traceVisitor struct {
	w       io.Writer
	prefix  []byte
	indents []byte
}

func (v *traceVisitor) Fprintf(w io.Writer, format string, args ...interface{}) {
	log.Output(2, string(v.indents)+fmt.Sprintf(format, args...))
}

func (v *traceVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Field:
		if len(n.Names) > 0 {
			v.Fprintf(v.w, "%T: %s\n", n, n.Names[0].Name)
		}
	case *ast.FieldList:
		v.Fprintf(v.w, "%T: %d\n", n, len(n.List))
	case *ast.Comment:
	case *ast.CommentGroup:
		if n != nil {
			v.Fprintf(v.w, "CommentGroup: %s\n", n.Text())
		}
	case *ast.File:
	case *ast.Package:
	case *ast.Ident:
		if n != nil {
			v.Fprintf(v.w, "%T: %s\n", n, n.Name)
		} else {
			v.Fprintf(v.w, "%T(nil)\n", n)
		}
	case *ast.BasicLit:
		if n != nil {
			v.Fprintf(v.w, "%T: %s\n", n, n.Value)
		} else {
			v.Fprintf(v.w, "%T(nil)\n", n)
		}
	case *ast.BasicType:
		if n != nil && n.Name != nil {
			v.Fprintf(v.w, "%T: %s\n", n, n.Name.Name)
		} else {
			v.Fprintf(v.w, "%T(nil)\n", n)
		}
	case *ast.ArrayType:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.MapType:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.VectorType:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.StructType:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.GenDecl:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.BeanDecl:
		if n != nil {
			v.Fprintf(v.w, "%T: %s\n", n, n.Kind)
		} else {
			v.Fprintf(v.w, "%T(nil)\n", n)
		}
	case *ast.ImportSpec:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.ConstSpec:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.BadDecl:
		v.Fprintf(v.w, "%T\n", n)
	case *ast.BadExpr:
		v.Fprintf(v.w, "%T\n", n)
	}
	//v.Fprintf(v.w, "%T\n", n)
	return v
}

func (v *traceVisitor) In() {
	v.indents = append(v.indents, v.prefix...)
}

func (v *traceVisitor) Out() {
	v.indents = v.indents[:len(v.indents)-len(v.prefix)]
}

func TestParser(t *testing.T) {
	w := os.Stdout
	src := []byte(`package demo;

import "a/b/c";
import x "o/p/q";

struct User {
	int32 id;
}

protocol UserInfo {
	int16 a,b;
	bool x;
	string y;
	int32 z;
	vector<int32> list;
	array<int32,6> arr;
}

service HelloWorld {
	name()
	say(string s)
	abc(int32 a, bool b)
}

// doc
enum Type {
	A = 1,
	B = 2,
}
`)

	fset := lexer.NewFileSet()
	file, err := ParseFile(fset, "demo.mid", src)
	if err != nil {
		t.Errorf("parse error: %v", err)
		return
	}
	ast.Walk(file, &traceVisitor{w: w, prefix: []byte(". . ")})

	for _, unresolved := range file.Unresolved {
		log.Printf("unresolved ident: %s (pos: %v)", unresolved.Name, fset.Position(unresolved.Pos))
	}
}
