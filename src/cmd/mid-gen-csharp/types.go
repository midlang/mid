package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/mid/lexer"
)

func csFieldDecl(f *build.Field, emptyIfNoName bool) string {
	if len(f.Names) == 0 {
		if emptyIfNoName {
			return buildType(f.Type)
		}
		return "_ " + buildType(f.Type)
	}
	return strings.Join(f.Names, ", ") + " " + buildType(f.Type)
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
			return "bool"
		case lexer.Byte:
			return "byte"
		case lexer.Bytes:
			return "byte[]"
		case lexer.String:
			return "string"
		case lexer.Int:
			return "int"
		case lexer.Int8:
			return "int8"
		case lexer.Int16:
			return "int16"
		case lexer.Int32:
			return "int32"
		case lexer.Int64:
			return "int64"
		case lexer.Uint:
			return "uint"
		case lexer.Uint8:
			return "uint8"
		case lexer.Uint16:
			return "uint16"
		case lexer.Uint32:
			return "uint32"
		case lexer.Uint64:
			return "uint64"
		case lexer.Float32:
			return "float"
		case lexer.Float64:
			return "double"
		default:
			panic("unknown builtin type `" + t.Name + "`")
		}
	case *build.ArrayType:
		size, ok := build.IntFromExpr(t.Size)
		if !ok {
			panic("array.Size not an integer")
		}
		_ = size
		return fmt.Sprintf("%s[]", buildType(t.T))
	case *build.VectorType:
		return fmt.Sprintf("%s[]", buildType(t.T))
	case *build.MapType:
		return fmt.Sprintf("Dictionary<%s, %s>", buildType(t.K), buildType(t.V))
	case *build.StructType:
		if t.Package != "" {
			return t.Package + "." + t.Name
		}
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
				buf.WriteString(csFieldDecl(field, allNoName))
			}
		}
		buf.WriteByte(')')
		if t.Result != nil {
			buf.WriteString(buildType(t.Result))
		}
		return buf.String()
	default:
		return ""
	}
}
