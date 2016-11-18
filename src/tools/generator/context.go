package generator

import (
	"github.com/midlang/mid/src/mid/build"
)

// Context holds current context for generate codes of specified package
type Context struct {
	// Pkg represents current package
	Pkg *build.Package
	// Root represents root template
	Root *Template
	// Plugin represents current plugin
	Plugin build.Plugin
	// Config represents current plugin runtime config
	Config build.PluginRuntimeConfig

	buildType BuildTypeFunc

	// beans holds all beans in current package
	beans map[string]*build.Bean
}

// NewContext creates a context by buildType,plugin,plugin_config
func NewContext(
	buildType BuildTypeFunc,
	plugin build.Plugin,
	config build.PluginRuntimeConfig,
) *Context {
	ctx := &Context{
		buildType: buildType,
		Plugin:    plugin,
		Config:    config,
	}
	return ctx
}

func (ctx *Context) initWithPkg(pkg *build.Package) {
	ctx.Pkg = pkg
	ctx.beans = make(map[string]*build.Bean)
	for _, file := range ctx.Pkg.Files {
		for _, bean := range file.Beans {
			ctx.beans[bean.Name] = bean
		}
	}
}

// BuildType executes buildType function
func (ctx *Context) BuildType(typ build.Type) string {
	return ctx.buildType(typ)
}

// Env gets custom envvar
func (ctx *Context) Env(key string) string {
	if ctx.Config.Envvars == nil {
		return ""
	}
	return ctx.Config.Envvars[key]
}

// FindBean finds bean by name in current package
func (ctx *Context) FindBean(name string) *build.Bean {
	return ctx.beans[name]
}
