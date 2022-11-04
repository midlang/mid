package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/mid/lexer"
)

func tsFieldDecl(f *build.Field, emptyIfNoName bool) string {
	if len(f.Names) == 0 {
		return buildType(f.Type)
	}
	return strings.Join(f.Names, ", ")
}

func basicIntExprString(expr build.Expr) (string, bool) {
	switch e := expr.(type) {
	case *build.BasicLit:
		if e.Kind == lexer.INT {
			return e.Value, true
		}
	}
	return "", false
}

func buildType(typ build.Type) string {
	switch t := typ.(type) {
	case *build.BasicType:
		builtinType, ok := lexer.LookupType(t.Name)
		if !ok {
			panic("type `" + t.Name + "` not a builtin type")
		}
		switch builtinType {
		case lexer.Bool:
			return "boolean"
		case lexer.Byte:
			return "number"
		case lexer.Bytes:
			return "Unit8Array"
		case lexer.String:
			return "string"
		case lexer.Int,
			lexer.Int8,
			lexer.Int16,
			lexer.Int32,
			lexer.Int64,
			lexer.Uint,
			lexer.Uint8,
			lexer.Uint16,
			lexer.Uint32,
			lexer.Uint64:
			return "number"
		default:
			panic("unknown builtin type `" + t.Name + "`")
		}
	case *build.ArrayType:
		return fmt.Sprintf("%s[]", buildType(t.T))
	case *build.VectorType:
		return fmt.Sprintf("%s[]", buildType(t.T))
	case *build.MapType:
		return fmt.Sprintf("{[key: %s]: %s}", buildType(t.K), buildType(t.V))
	case *build.StructType:
		return t.Name
	case *build.FuncType:
		var buf bytes.Buffer
		buf.WriteByte('(')
		if len(t.Params) > 0 {
			allNoName := true
			for _, field := range t.Params {
				if len(field.Names) > 0 {
					allNoName = false
					break
				}
			}
			for i, field := range t.Params {
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(tsFieldDecl(field, allNoName))
			}
		}
		buf.WriteByte(')')
		return buf.String()
	default:
		return ""
	}
}
