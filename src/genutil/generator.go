package genutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/midlang/mid/src/mid/build"
	"github.com/mkideal/log"
	"github.com/mkideal/pkg/errors"
	"github.com/mkideal/pkg/textutil/namemapper"
)

// BuildTypeFunc is a function type which used to build `build.Type` to a string
type BuildTypeFunc func(build.Type) string

var (
	// funcs holds all shared template functions
	funcs template.FuncMap
	// context holds current context information
	context *Context
)

func firstOf(sep, s string) string {
	return nthOf(sep, s, 0)
}

func lastOf(sep, s string) string {
	return nthOf(sep, s, -1)
}

func nthOf(sep, s string, n int) string {
	values := strings.Split(s, sep)
	l := len(values)
	if l == 0 {
		return ""
	}
	return values[(n%l+l)%l]
}

// Init initializes generator
// NOTE: You MUST initialize generator before using generator
//
// buildType is a function for building build.Type to a string
// plugin is the language plugin
// config is runtime config of the plugin
func Init(
	buildType BuildTypeFunc,
	plugin build.Plugin,
	config build.PluginRuntimeConfig,
) {
	// creates context
	context = NewContext(buildType, plugin, config)

	funcs = template.FuncMap{
		// context returns context
		"context": func() *Context { return context },
		// outdir returns output directory
		"outdir": func() string { return context.Config.Outdir },
		// error print error log and returns an error
		"error": func(format string, args ...interface{}) error {
			err := fmt.Errorf(format, args...)
			log.Error("Error: %v", err)
			return err
		},
		// isInt check whether the type is an integer
		"isInt": func(typ string) bool {
			switch typ {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
				return true
			default:
				return false
			}
		},
		// include includes a file
		"include": func(filename string) (string, error) {
			if !filepath.IsAbs(filename) {
				filename = filepath.Join(context.Plugin.TemplatesDir, IncludesDir, filename)
			}
			content, err := ioutil.ReadFile(filename)
			return string(content), err
		},
		// include_template includes a template file with `data`
		// NOTE: include_template ignores meta header
		"includeTemplate": func(filename string, data interface{}) (string, error) {
			if !filepath.IsAbs(filename) {
				filename = filepath.Join(context.Plugin.TemplatesDir, IncludesDir, filename)
			}
			_, temp, err := ParseTemplateFile(filename)
			if err != nil {
				return "", err
			}
			var buf bytes.Buffer
			pwd := context.Pwd
			context.Pwd, _ = filepath.Split(filename)
			err = temp.Execute(&buf, data)
			context.Pwd = pwd
			return buf.String(), err
		},
		"valueAt": func(values []interface{}, index int) interface{} { return values[index] },
		"slice":   func(values ...interface{}) []interface{} { return values },
		// pwd returns current template file directory
		"pwd": func() string { return context.Pwd },
		// joinPath joins file paths
		"joinPath": func(paths ...string) string { return filepath.Join(paths...) },
		// os
		"osenv": func(key string) string { return os.Getenv(key) },
		// values
		"newBool":   func() *Bool { b := Bool(false); return &b },
		"newString": func() *String { s := String(""); return &s },
		"newInt":    func() *Int { i := Int(0); return &i },

		// string operations
		"string":      func(v interface{}) string { return fmt.Sprintf("%v", v) },
		"title":       func(s string) string { return strings.Title(s) },
		"toLower":     func(s string) string { return strings.ToLower(s) },
		"toUpper":     func(s string) string { return strings.ToUpper(s) },
		"contains":    func(sub, s string) bool { return strings.Contains(s, sub) },
		"containsAny": func(chars, s string) bool { return strings.ContainsAny(s, chars) },
		"count":       func(sep, s string) int { return strings.Count(s, sep) },
		"index":       func(sep, s string) int { return strings.Index(s, sep) },
		"lastIndex":   func(sep, s string) int { return strings.LastIndex(s, sep) },
		"join":        func(sep string, strs ...string) string { return strings.Join(strs, sep) },
		"joinStrings": func(sep string, strs []string) string { return strings.Join(strs, sep) },
		"split":       func(sep, s string) []string { return strings.Split(s, sep) },
		"splitN":      func(sep string, n int, s string) []string { return strings.SplitN(s, sep, n) },
		"stringAt":    func(strs []string, index int) string { return strs[index] },
		"repeat":      func(count int, s string) string { return strings.Repeat(s, count) },
		"replace":     func(old, new string, n int, s string) string { return strings.Replace(s, old, new, n) },
		"hasPrefix":   func(prefix string, s string) bool { return strings.HasPrefix(s, prefix) },
		"hasSuffix":   func(suffix string, s string) bool { return strings.HasSuffix(s, suffix) },
		"trimPrefix":  func(prefix string, s string) string { return strings.TrimPrefix(s, prefix) },
		"trimSuffix":  func(suffix string, s string) string { return strings.TrimSuffix(s, suffix) },
		"firstOf":     func(sep, s string) string { return firstOf(sep, s) },
		"lastOf":      func(sep, s string) string { return lastOf(sep, s) },
		"nthOf":       func(sep, s string, n int) string { return nthOf(sep, s, n) },
		"oneof": func(s string, set ...string) bool {
			for _, s2 := range set {
				if s == s2 {
					return true
				}
			}
			return false
		},
		"trimSpace": func(s string) string { return strings.TrimSpace(s) },
		"append":    func(appended string, s string) string { return s + appended },
		"substr": func(startIndex, endIndex int, s string) string {
			n := len(s)
			if n == 0 {
				return ""
			}
			if startIndex < 0 {
				startIndex = startIndex%n + n
			}
			if endIndex < 0 {
				endIndex = endIndex%n + n
			}
			if endIndex > n {
				endIndex = n
			}
			if startIndex > endIndex {
				return ""
			}
			return s[startIndex:endIndex]
		},
		"underScore": func(s string) string { return namemapper.UnderScore(s) },
		"upper":      func(s string) string { return namemapper.Upper(s) },
		"lower":      func(s string) string { return namemapper.Lower(s) },
		"upperCamel": func(s string) string { return namemapper.UpperCamel(s) },
		"lowerCamel": func(s string) string { return namemapper.LowerCamel(s) },
		// logic
		"OR": func(bools ...bool) bool {
			for _, b := range bools {
				if b {
					return true
				}
			}
			return false
		},
		"AND": func(bools ...bool) bool {
			for _, b := range bools {
				if !b {
					return false
				}
			}
			return true
		},
		"NOT": func(b bool) bool { return !b },
		"XOR": func(b1, b2 bool) bool { return b1 != b2 },
	}
}

// GeneratePackage generates codes for package
func GeneratePackage(pkg *build.Package) (files map[string]bool, err error) {
	if context == nil {
		return nil, errors.Throw("generator not initialized")
	}

	infos, err := OpenTemplatesDir(context.Plugin.Lang, context.Plugin.TemplatesDir)
	if err != nil {
		return nil, errors.Throw(err.Error())
	}
	if len(infos) == 0 {
		log.With(context.Plugin.Lang).Warn("no templates found")
		return nil, nil
	}

	// sets context.Pkg
	context.initWithPkg(pkg)
	context.Pwd = context.Plugin.TemplatesDir

	outdir := context.Config.Outdir
	// NOTE: environment variable nopkgdir
	if !context.Config.BoolEnv("nopkgdir") {
		outdir = filepath.Join(outdir, pkg.Name)
	}
	constDecls := make([]*GenDecl, 0)
	files = make(map[string]bool)
	for _, info := range infos {
		log.With(context.Plugin.Lang).Debug("template file: %s", info.Name())
		filename := filepath.Join(context.Plugin.TemplatesDir, info.Name())
		meta, temp, err := ParseTemplateFile(filename)
		if err != nil {
			return files, err
		}
		oldMetaFile := meta.File

		var file io.WriteCloser
		kind, suffix := ParseTemplateFilename(info.Name())
		log.Debug("kind=%s, suffix=%s", kind, suffix)

		// sets context.Root and context.Kind
		context.Root = temp
		context.Kind = kind
		context.Suffix = suffix

		// apply template to specific kind node
		switch kind {
		case "package":
			dftName := pkg.Name + "." + suffix
			ctxPkg := Package{Package: pkg}
			if file, err = ApplyMeta(outdir, meta, ctxPkg, dftName); err == nil {
				err = temp.Execute(file, ctxPkg)
				files[meta.File] = true
				file.Close()
				if err != nil {
					return files, err
				}
			} else {
				return files, err
			}
		case "file":
			for _, f := range pkg.Files {
				_, filename := filepath.Split(f.Filename)
				filename = firstOf(".", filename)
				dftName := filename + "." + suffix
				meta.File = oldMetaFile
				ctxFile := File{File: f}
				if file, err = ApplyMeta(outdir, meta, ctxFile, dftName); err == nil {
					files[meta.File] = true
					err = temp.Execute(file, ctxFile)
					file.Close()
					if err != nil {
						return files, err
					}
				} else {
					return files, err
				}
			}
		case "const":
			if len(constDecls) == 0 {
				for _, f := range pkg.Files {
					for _, c := range f.Decls {
						if len(c.Consts) > 0 {
							constDecls = append(constDecls, NewGenDecl(f, c))
						}
					}
				}
			}
			if len(constDecls) > 0 {
				meta.File = oldMetaFile
				if file, err = ApplyMeta(outdir, meta, constDecls, "constants."+suffix); err == nil {
					files[meta.File] = true
					err = temp.Execute(file, constDecls)
					file.Close()
					if err != nil {
						return files, err
					}
				} else {
					return files, err
				}
			}
		case "group":
			for _, f := range pkg.Files {
				for _, g := range f.Groups {
					dftName := g.Name + "." + suffix
					meta.File = oldMetaFile
					group := NewGroup(f, g)
					if file, err = ApplyMeta(outdir, meta, group, dftName); err == nil {
						files[meta.File] = true
						err = temp.Execute(file, group)
						file.Close()
						if err != nil {
							return files, err
						}
					} else {
						return files, err
					}
				}
			}
		// beans: enum,struct,protocol,service
		default:
			for _, f := range pkg.Files {
				for _, b := range f.Beans {
					if b.Kind == kind {
						dftName := b.Name + "." + suffix
						meta.File = oldMetaFile
						bean := NewBean(f, b)
						if file, err = ApplyMeta(outdir, meta, bean, dftName); err == nil {
							files[meta.File] = true
							err = temp.Execute(file, bean)
							file.Close()
							if err != nil {
								return files, err
							}
						} else {
							return files, err
						}
					}
				}
			}
		}
	}

	return files, nil
}
