package genutil

import (
	"bytes"

	"github.com/midlang/mid/src/mid/build"
)

// Package wraps build.Package
type Package struct {
	*build.Package
}

// Depcrated API, use `Gen' instead
func (pkg Package) GenerateDeclsBySubTemplates() (string, error) {
	return pkg.Gen()
}

// Gen generates file by predefined sub-template `T_const`, `T_group`, `T_<kind>`
func (pkg Package) Gen() (string, error) {
	buf := new(bytes.Buffer)

	if temp := context.Root.Lookup("T_const"); temp != nil {
		constDecls := make([]*GenDecl, 0)
		for _, f := range pkg.Files {
			for _, c := range f.Decls {
				if len(c.Consts) > 0 {
					constDecls = append(constDecls, NewGenDecl(f, c))
				}
			}
		}
		if len(constDecls) > 0 {
			if err := temp.Execute(buf, constDecls); err != nil {
				return "", err
			}
		}
	}

	if temp := context.Root.Lookup("T_group"); temp != nil {
		for _, f := range pkg.Files {
			context.Filename = f.Filename
			for _, g := range f.Groups {
				if err := temp.Execute(buf, NewGroup(f, g)); err != nil {
					context.Filename = ""
					return "", err
				}
			}
			context.Filename = ""
		}
	}

	for _, f := range pkg.Files {
		context.Filename = f.Filename
		for _, b := range f.Beans {
			if temp := context.Root.Lookup("T_" + b.Kind); temp != nil {
				if err := temp.Execute(buf, NewBean(f, b)); err != nil {
					context.Filename = ""
					return "", err
				}
			}
		}
		context.Filename = ""
	}
	return buf.String(), nil
}

// File wraps build.File
type File struct {
	*build.File
	groups map[string]*build.Group
}

func (f *File) FindGroup(name string) *build.Group {
	if f.groups == nil {
		f.groups = make(map[string]*build.Group)
		pendings := make(map[string]*build.Group)
		for _, g := range f.Groups {
			pendings[g.Name] = g
		}
		for len(pendings) > 0 {
			var first *build.Group
			for _, v := range pendings {
				first = v
				break
			}
			delete(pendings, first.Name)
			f.groups[first.Name] = first
			for _, g := range first.Groups {
				pendings[g.Name] = g
			}
		}
	}
	return f.groups[name]
}

// Depcrated API, use `Gen' instead
func (f File) GenerateDeclsBySubTemplates() (string, error) {
	return f.Gen()
}

// Gen generates file by predefined sub-template `T_const`, `T_<kind>`
func (f File) Gen() (string, error) {
	buf := new(bytes.Buffer)

	if temp := context.Root.Lookup("T_const"); temp != nil {
		for _, c := range f.Decls {
			if len(c.Consts) > 0 {
				if err := temp.Execute(buf, NewGenDecl(f.File, c)); err != nil {
					return "", err
				}
			}
		}
	}

	for _, b := range f.Beans {
		context.Filename = f.Filename
		if temp := context.Root.Lookup("T_" + b.Kind); temp != nil {
			if err := temp.Execute(buf, NewBean(f.File, b)); err != nil {
				context.Filename = ""
				return "", err
			}
		}
		context.Filename = ""
	}
	return buf.String(), nil
}

// GenDecl wraps build.GenDecl
type GenDecl struct {
	*build.GenDecl
	File *File
}

func NewGenDecl(file *build.File, c *build.GenDecl) *GenDecl {
	return &GenDecl{
		GenDecl: c,
		File:    &File{File: file},
	}
}

// Bean wraps build.Bean
type Bean struct {
	*build.Bean
	File *File
}

func NewBean(file *build.File, b *build.Bean) *Bean {
	return &Bean{
		Bean: b,
		File: &File{File: file},
	}
}

// Extends gets extends of bean as a string slice
func (bean *Bean) BuildExtends(ctx *Context) []string {
	extends := bean.Bean.Extends
	if len(extends) == 0 {
		return nil
	}
	var strs []string
	for _, ex := range extends {
		strs = append(strs, ctx.BuildType(ex))
	}
	return strs
}

// AddTag is a chain function for adding tag
func (bean *Bean) AddTag(key, value string, field *build.Field) *build.Field {
	return bean.addTag(key, value, field, true)
}

func (bean *Bean) AddTagNX(key, value string, field *build.Field) *build.Field {
	return bean.addTag(key, value, field, false)
}

func (bean *Bean) addTag(key, value string, field *build.Field, force bool) *build.Field {
	_, found := field.Tag.Lookup(key)
	if !found || force {
		field.Tag.Set(key, value)
	}
	return field
}

// Group wraps build.Group
type Group struct {
	*build.Group
	File *File
}

func NewGroup(file *build.File, g *build.Group) *Group {
	return &Group{
		Group: g,
		File:  &File{File: file},
	}
}
