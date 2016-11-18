package generator

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
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

// Init initialize generator
// NOTE: You MUST initialize generator before using generator
//
// buildType is a function for build build.Type to a string
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
		// include includes a text file
		"include": func(filename string) (string, error) {
			content, err := ioutil.ReadFile(filename)
			return string(content), err
		},
		// include_template includes a template file with `data`
		// NOTE: include_template ignores meta header
		"include_template": func(filename string, data interface{}) (string, error) {
			_, temp, err := ParseTemplateFile(filename)
			if err != nil {
				return "", err
			}
			var buf bytes.Buffer
			err = temp.Execute(&buf, data)
			return buf.String(), err
		},
		// os
		"osenv": func(key string) string { return os.Getenv(key) },
		// string operations
		"title":       func(s string) string { return strings.Title(s) },
		"toLower":     func(s string) string { return strings.ToLower(s) },
		"toUpper":     func(s string) string { return strings.ToUpper(s) },
		"contains":    func(s, sub string) bool { return strings.Contains(s, sub) },
		"containsAny": func(s, chars string) bool { return strings.ContainsAny(s, chars) },
		"count":       func(s, sep string) int { return strings.Count(s, sep) },
		"index":       func(s, sep string) int { return strings.Index(s, sep) },
		"lastIndex":   func(s, sep string) int { return strings.LastIndex(s, sep) },
		"join":        func(strs []string, sep string) string { return strings.Join(strs, sep) },
		"split":       func(s, sep string) []string { return strings.Split(s, sep) },
		"splitN":      func(s, sep string, n int) []string { return strings.SplitN(s, sep, n) },
		"repeat":      func(s string, count int) string { return strings.Repeat(s, count) },
		"replace":     func(s, old, new string, n int) string { return strings.Replace(s, old, new, n) },
		"hasPrefix":   func(s string, prefix string) bool { return strings.HasPrefix(s, prefix) },
		"hasSuffix":   func(s string, suffix string) bool { return strings.HasSuffix(s, suffix) },
		"trimPrefix":  func(s string, prefix string) string { return strings.TrimPrefix(s, prefix) },
		"trimSuffix":  func(s string, suffix string) string { return strings.TrimSuffix(s, suffix) },
		"trimSpace":   func(s string) string { return strings.TrimSpace(s) },
		"append":      func(appended string, origin string) string { return origin + appended },
		"substr": func(s string, startIndex, endIndex int) string {
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
	}

}

// GoFmt formats go code file
func GoFmt(filename string) error {
	const gofmt = "gofmt"
	if _, err := exec.LookPath(gofmt); err != nil {
		// do nothing if failed to lookup `gofmt`
		return nil
	}
	cmd := exec.Command(gofmt, "-w", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GeneratePackage generates code for package
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

	outdir := filepath.Join(context.Config.Outdir, pkg.Name)
	files = make(map[string]bool)
	for _, info := range infos {
		log.With(context.Plugin.Lang).Debug("template file: %s", info.Name())
		filename := filepath.Join(context.Plugin.TemplatesDir, info.Name())
		meta, temp, err := ParseTemplateFile(filename)
		if err != nil {
			return files, err
		}
		oldMetaFile := meta.File

		var file *os.File
		kind, suffix := ParseTemplateFilename(info.Name())
		log.Debug("kind=%s, suffix=%s", kind, suffix)

		// sets context.Root
		context.Root = temp
		switch kind {
		case "package":
			dftName := pkg.Name + "." + suffix
			if file, err = ApplyMeta(outdir, meta, pkg, dftName); err == nil {
				err = temp.Execute(file, Package{Package: pkg})
				files[meta.File] = true
				file.Close()
				if err != nil {
					return files, err
				}
			}
		case "file":
			for _, f := range pkg.Files {
				dftName := f.Filename + "." + suffix
				meta.File = oldMetaFile
				if file, err = ApplyMeta(outdir, meta, f, dftName); err == nil {
					files[meta.File] = true
					err = temp.Execute(file, File{File: f})
					file.Close()
					if err != nil {
						return files, err
					}
				}
			}
		case "const":
			for _, f := range pkg.Files {
				for _, c := range f.Decls {
					if len(c.Consts) > 0 {
						meta.File = oldMetaFile
						if file, err = ApplyMeta(outdir, meta, c, "constants."+suffix); err == nil {
							files[meta.File] = true
							err = temp.Execute(file, NewGenDecl(f, c))
							file.Close()
							if err != nil {
								return files, err
							}
						}
					}
				}
			}
		// beans: enum,struct,...
		default:
			for _, f := range pkg.Files {
				for _, b := range f.Beans {
					if b.Kind == kind {
						dftName := b.Name + "." + suffix
						meta.File = oldMetaFile
						if file, err = ApplyMeta(outdir, meta, b, dftName); err == nil {
							files[meta.File] = true
							err = temp.Execute(file, NewBean(f, b))
							file.Close()
							if err != nil {
								return files, err
							}
						}
					}
				}
			}
		}
	}

	return files, nil
}
