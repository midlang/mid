package generator

import (
	"github.com/midlang/mid/src/mid/build"
)

type Context struct {
	Pkg       *build.Package
	Root      *Template
	Envvars   map[string]string
	buildType BuildTypeFunc
	beans     map[string]*build.Bean
}

func NewContext(pkg *build.Package, rootTemp *Template, envvars map[string]string, buildType BuildTypeFunc) *Context {
	ctx := &Context{
		Pkg:       pkg,
		Root:      rootTemp,
		Envvars:   envvars,
		buildType: buildType,
	}
	ctx.beans = make(map[string]*build.Bean)
	for _, file := range ctx.Pkg.Files {
		for _, bean := range file.Beans {
			ctx.beans[bean.Name] = bean
		}
	}
	return ctx
}

func (ctx *Context) BuildType(typ build.Type) string {
	return ctx.buildType(typ)
}

func (ctx *Context) Env(key string) string {
	return ctx.Envvars[key]
}

func (ctx *Context) FindBean(name string) *build.Bean {
	return ctx.beans[name]
}
