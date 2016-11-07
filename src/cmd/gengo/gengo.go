package main

import (
	//"bytes"
	"fmt"
	"os"
	"strings"

	//"github.com/midlang/mid/src/mid/ast"
	"github.com/midlang/mid/src/mid/build"
	"github.com/mkideal/log"
)

func main() {
	defer log.Uninit(log.InitConsole(log.LvWARN))
	log.SetLevel(log.LvTRACE)

	config, builder, err := build.ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	log.WithJSON(config).Trace("config")
	log.WithJSON(builder).Trace("builder")
}

func goFieldDecl(f *build.Field) string {
	if len(f.Names) == 0 {
		return "_ " + buildType(f.Type)
	}
	return strings.Join(f.Names, ", ") + " " + buildType(f.Type)
}

func buildType(typ build.Type) string {
	/*switch t := typ.(type) {
	case *build.BasicType:
		return t.Name
	case *build.ArrayType:
	return fmt.Sprintf("[%s]%s", build.BuildExpr(t.Size), buildType(t.T))
	case *ast.VectorType:
		return fmt.Sprintf("[]%s", buildType(t.T))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", buildType(t.K), buildType(t.V))
	case *ast.StructType:
		if t.Package != nil {
			return t.Package.Name + "." + t.Name.Name
		}
		return t.Name.Name
	case *ast.FuncType:
		var buf bytes.Buffer
		buf.WriteByte('(')
		if t.Params != nil && len(t.Params.List) > 0 {
			for i, field := range t.Params.List {
				if i > 0 {
					buf.WriteByte(',')
				}
				bf := build.BuildField(field)
				buf.WriteString(goFieldDecl(bf))
			}
		}
		buf.WriteByte(')')
		if t.Result != nil {
			buf.WriteString(buildType(t.Result))
		}
		return buf.String()
	}*/
	return ""
}
