package genutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gopherd/log"
	"github.com/mkideal/pkg/errors"
	"github.com/mkideal/pkg/textutil/namemapper"

	"github.com/midlang/mid/src/mid/build"
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

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func toNumber(x any) float64 {
	switch v := x.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case uintptr:
		return float64(v)
	default:
		return math.NaN()
	}
}

func toInt(x any) int64 {
	switch v := x.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case uintptr:
		return int64(v)
	default:
		return 0
	}
}

func includeTemplate(filename string, data interface{}) (string, error) {
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
		// Common functions

		// context returns context
		"context": func() *Context { return context },
		"debug": func(format string, args ...interface{}) error {
			log.Debug().Printf(format, args...)
			return nil
		},
		"info": func(format string, args ...interface{}) error {
			log.Info().Printf(format, args...)
			return nil
		},
		"warn": func(format string, args ...interface{}) error {
			log.Warn().Printf(format, args...)
			return nil
		},
		// error print error log and returns an error
		"error": func(format string, args ...interface{}) error {
			err := fmt.Errorf(format, args...)
			log.Error().Printf("Error: %v", err)
			return err
		},
		// include_template includes a template file with `data`
		// NOTE: includeTemplate ignores meta header
		"includeTemplate":  includeTemplate,
		"include_template": includeTemplate,
		// include includes a file
		"include": func(filename string) (string, error) {
			if !filepath.IsAbs(filename) {
				filename = filepath.Join(context.Plugin.TemplatesDir, IncludesDir, filename)
			}
			content, err := ioutil.ReadFile(filename)
			return string(content), err
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
		// joinPath joins file paths
		"joinPath": func(paths ...string) string { return filepath.Join(paths...) },
		// osenv gets env
		"osenv": func(key string) string { return os.Getenv(key) },
		// outdir returns output directory
		"outdir": func() string { return context.Config.Outdir },
		// pwd returns current template file directory
		"pwd":     func() string { return context.Pwd },
		"slice":   func(values ...interface{}) []interface{} { return values },
		"valueAt": func(values []interface{}, index int) interface{} { return values[index] },
		"bareFilename": func(filename string) string {
			filename = filepath.Base(filename)
			dotIndex := strings.Index(filename, ".")
			if dotIndex >= 0 {
				return filename[:dotIndex]
			}
			return filename
		},

		// values
		"newBool":   func() *Bool { b := Bool(false); return &b },
		"newString": func() *String { s := String(""); return &s },
		"newInt":    func() *Int { i := Int(0); return &i },

		// String functions

		"append":      func(appended string, s string) string { return s + appended },
		"containsAny": func(chars, s string) bool { return strings.ContainsAny(s, chars) },
		"contains":    func(substr, s string) bool { return strings.Contains(s, substr) },
		"count":       func(substr, s string) int { return strings.Count(s, substr) },
		"firstOf":     func(sep, s string) string { return firstOf(sep, s) },
		"hasPrefix":   func(prefix string, s string) bool { return strings.HasPrefix(s, prefix) },
		"hasSuffix":   func(suffix string, s string) bool { return strings.HasSuffix(s, suffix) },
		"index":       func(substr, s string) int { return strings.Index(s, substr) },
		"joinStrings": func(sep string, strs []string) string { return strings.Join(strs, sep) },
		"join":        func(sep string, strs ...string) string { return strings.Join(strs, sep) },
		"lastIndex":   func(sep, s string) int { return strings.LastIndex(s, sep) },
		"lastOf":      func(sep, s string) string { return lastOf(sep, s) },
		"lowerCamel":  func(s string) string { return namemapper.LowerCamel(s) },
		"nthOf":       func(sep string, n int, s string) string { return nthOf(sep, s, n) },
		"oneof": func(s string, set ...string) bool {
			for _, s2 := range set {
				if s == s2 {
					return true
				}
			}
			return false
		},
		"toNumber":   toNumber,
		"toInt":      toInt,
		"parseInt":   parseInt,
		"parseFloat": parseFloat,
		"parseBool":  parseBool,
		"lt": func(x, y any) bool {
			return toNumber(x) < toNumber(y)
		},
		"le": func(x, y any) bool {
			return toNumber(x) <= toNumber(y)
		},
		"gt": func(x, y any) bool {
			return toNumber(x) > toNumber(y)
		},
		"ge": func(x, y any) bool {
			return toNumber(x) >= toNumber(y)
		},
		"repeat":   func(count int, s string) string { return strings.Repeat(s, count) },
		"replace":  func(old, new string, n int, s string) string { return strings.Replace(s, old, new, n) },
		"splitN":   func(sep string, n int, s string) []string { return strings.SplitN(s, sep, n) },
		"split":    func(sep, s string) []string { return strings.Split(s, sep) },
		"stringAt": func(strs []string, index int) string { return strs[index] },
		"string":   func(data interface{}) string { return fmt.Sprintf("%v", data) },
		"substr": func(startIndex, endIndex int, s string) string {
			n := len(s)
			if n == 0 {
				return ""
			}
			if startIndex < 0 {
				startIndex = startIndex%n + n
			}
			if endIndex <= 0 {
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
		"title":      func(s string) string { return strings.Title(s) },
		"toLower":    func(s string) string { return strings.ToLower(s) },
		"toUpper":    func(s string) string { return strings.ToUpper(s) },
		"trimPrefix": func(prefix string, s string) string { return strings.TrimPrefix(s, prefix) },
		"trimSpace":  func(s string) string { return strings.TrimSpace(s) },
		"trimSuffix": func(suffix string, s string) string { return strings.TrimSuffix(s, suffix) },
		"underScore": func(s string) string { return namemapper.UnderScore(s) },
		"upperCamel": func(s string) string { return namemapper.UpperCamel(s) },
		"enumerate": func(begin, end int) []int {
			var slice = make([]int, end-begin)
			for i := 0; i < len(slice); i++ {
				slice[i] = begin + i
			}
			return slice
		},

		// Logical functions

		"AND": func(bools ...bool) bool {
			for _, b := range bools {
				if !b {
					return false
				}
			}
			return true
		},
		"NOT": func(b bool) bool { return !b },
		"OR": func(bools ...bool) bool {
			for _, b := range bools {
				if b {
					return true
				}
			}
			return false
		},
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
		log.Warn().
			String("plugin", context.Plugin.Lang).
			Print("no templates found")
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
		filename := filepath.Join(context.Plugin.TemplatesDir, info.Name())
		meta, temp, err := ParseTemplateFile(filename)
		if err != nil {
			return files, err
		}
		oldMetaFile := meta.File

		var file io.WriteCloser
		kind, suffix := ParseTemplateFilename(info.Name())

		// sets context.Root and context.Kind
		context.Root = temp
		context.Kind = kind
		context.Suffix = suffix
		context.Filename = ""

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
				context.Filename = f.Filename
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
				context.Filename = ""
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
				context.Filename = f.Filename
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
				context.Filename = ""
			}
		// beans: enum,struct,protocol,service
		default:
			for _, f := range pkg.Files {
				context.Filename = f.Filename
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
				context.Filename = ""
			}
		}
	}

	return files, nil
}
