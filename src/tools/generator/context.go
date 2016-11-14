package generator

import (
	"strings"

	"github.com/midlang/mid/src/mid/build"
)

type Context struct {
	Pkg       *build.Package
	Root      *Template
	buildType BuildTypeFunc
}

func NewContext(pkg *build.Package, rootTemp *Template, buildType BuildTypeFunc) *Context {
	return &Context{
		Pkg:       pkg,
		Root:      rootTemp,
		buildType: buildType,
	}
}

func (ctx *Context) BuildType(typ build.Type) string {
	return ctx.buildType(typ)
}

func (ctx *Context) Title(s string) string                    { return strings.Title(s) }
func (ctx *Context) ToLower(s string) string                  { return strings.ToLower(s) }
func (ctx *Context) ToUpper(s string) string                  { return strings.ToUpper(s) }
func (ctx *Context) Contains(s, sub string) bool              { return strings.Contains(s, sub) }
func (ctx *Context) ContainsAny(s, chars string) bool         { return strings.ContainsAny(s, chars) }
func (ctx *Context) Count(s, sep string) int                  { return strings.Count(s, sep) }
func (ctx *Context) Index(s, sep string) int                  { return strings.Index(s, sep) }
func (ctx *Context) LastIndex(s, sep string) int              { return strings.LastIndex(s, sep) }
func (ctx *Context) Join(strs []string, sep string) string    { return strings.Join(strs, sep) }
func (ctx *Context) Split(s, sep string) []string             { return strings.Split(s, sep) }
func (ctx *Context) SplitN(s, sep string, n int) []string     { return strings.SplitN(s, sep, n) }
func (ctx *Context) Repeat(s string, count int) string        { return strings.Repeat(s, count) }
func (ctx *Context) Replace(s, old, new string, n int) string { return strings.Replace(s, old, new, n) }
