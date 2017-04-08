package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/midlang/mid/src/mid/build"
	"github.com/midlang/mid/src/mid/lexer"
)

const (
	Env_unordered_map = "cpp:unordered_map"
)

var (
	config build.PluginRuntimeConfig
)

func cppFieldDecl(f *build.Field) string {
	if len(f.Names) == 0 {
		return ""
	}
	return buildType(f.Type) + " " + strings.Join(f.Names, ", ")
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
			return "bool"
		case lexer.Byte:
			return "unsigned char"
		case lexer.Bytes:
			return "unsigned char*"
		case lexer.String:
			return "std::string"
		case lexer.Int:
			return "int"
		case lexer.Int8:
			return "int8_t"
		case lexer.Int16:
			return "int16_t"
		case lexer.Int32:
			return "int32_t"
		case lexer.Int64:
			return "int64_t"
		case lexer.Uint:
			return "uint_t"
		case lexer.Uint8:
			return "uint8_t"
		case lexer.Uint16:
			return "uint16_t"
		case lexer.Uint32:
			return "uint32_t"
		case lexer.Uint64:
			return "uint64_t"
		case lexer.Float32:
			return "float"
		case lexer.Float64:
			return "double"
		default:
			panic("unknown builtin type `" + t.Name + "`")
		}
	case *build.ArrayType:
		size, ok := basicIntExprString(t.Size)
		if !ok {
			panic("array.Size not a integer")
		}
		return fmt.Sprintf("std::array<%s,%s> ", buildType(t.T), size)
	case *build.VectorType:
		return fmt.Sprintf("std::vector<%s> ", buildType(t.T))
	case *build.MapType:
		if config.BoolEnv(Env_unordered_map) {
			return fmt.Sprintf("std::unordered_map<%s,%s> ", buildType(t.K), buildType(t.V))
		} else {
			return fmt.Sprintf("std::map<%s,%s> ", buildType(t.K), buildType(t.V))
		}
	case *build.StructType:
		if t.Package != "" {
			return t.Package + "::" + t.Name
		}
		return t.Name
	case build.TypeBase, *build.TypeBase:
		return "void"
	case *build.FuncType:
		var buf bytes.Buffer
		if t.Result != nil {
			buf.WriteString(buildType(t.Result))
		} else {
			buf.WriteString("void")
		}
		buf.WriteByte('(')
		if len(t.Params) > 0 {
			for i, field := range t.Params {
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(cppFieldDecl(field))
			}
		}
		buf.WriteByte(')')
		return buf.String()
	default:
		return ""
	}
}
