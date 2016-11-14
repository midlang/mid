package generator

import (
	"github.com/midlang/mid/src/mid/build"
)

type Context struct {
	Pkg       *build.Package
	Root      *Template
	Envvars   map[string]string
	buildType BuildTypeFunc
}

func NewContext(pkg *build.Package, rootTemp *Template, envvars map[string]string, buildType BuildTypeFunc) *Context {
	return &Context{
		Pkg:       pkg,
		Root:      rootTemp,
		Envvars:   envvars,
		buildType: buildType,
	}
}

func (ctx *Context) BuildType(typ build.Type) string {
	return ctx.buildType(typ)
}

func (ctx *Context) Env(key string) string {
	return ctx.Envvars[key]
}
