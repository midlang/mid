package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/midlang/mid/src/mid/build"
	"github.com/mkideal/log"
)

const (
	IncludesDir        = "includes"
	TemplateFileSuffix = ".temp"
)

type BuildTypeFunc func(build.Type) string

func OpenTemplatesDir(lang, dir string) ([]os.FileInfo, error) {
	// open templates directory
	fs, err := os.Open(dir)
	if err != nil {
		log.With(lang).Error("open templates directory %s error: %v", dir, err)
		return nil, err
	}
	defer fs.Close()
	infos, err := fs.Readdir(-1)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	if len(infos) == 0 {
		log.With(lang).Warn("no templates found")
		return nil, nil
	}
	return infos, nil
}

type TemplateMeta struct {
	Dir    string
	File   string
	Date   string
	Values map[string]string
}

func ParseTemplateFile(filename string) (*TemplateMeta, *template.Template, error) {
	meta := &TemplateMeta{}
	// TODO: parse template file meta info
	// e.g.
	//
	// ---
	// dir: aaa
	// file: {{.Name}}.go
	// ---
	temp, err := template.ParseFiles(filename)
	if err != nil {
		return nil, nil, err
	}
	return meta, temp, err
}

var Funcs = func(buildType BuildTypeFunc) map[string]interface{} {
	return map[string]interface{}{
		"include": func(path string, data interface{}) error {
			// TODO
			return nil
		},
		"buildType": buildType,
	}
}

func ApplyMeta(outdir string, meta *TemplateMeta, data interface{}, dftName string) (*os.File, error) {
	//TODO: execute template for meta

	if meta.File == "" {
		meta.File = dftName
	}

	if !filepath.IsAbs(meta.File) {
		meta.File = filepath.Join(outdir, meta.File)
	}
	dir, _ := filepath.Split(meta.File)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error("mkdir %s error: %v", dir, err)
			return nil, err
		}
	}
	file, err := os.OpenFile(meta.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Error("open file %s error: %v", meta.File, err)
		return nil, err
	}
	return file, err
}

func TemplateFileKind(filename string) (kind, suffix string) {
	filename = strings.TrimSuffix(filename, TemplateFileSuffix)
	strs := strings.Split(filename, ".")
	if len(strs) == 0 {
		return filename, ""
	}
	if len(strs) == 1 {
		return strs[0], ""
	}
	return strs[0], strs[1]
}

func GoFmt(filename string) error {
	cmd := exec.Command("gofmt", "-w", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GeneratePackage(
	outdir string,
	pkg *build.Package,
	plugin build.Plugin,
	config build.PluginRuntimeConfig,
	buildType BuildTypeFunc,
	templates []os.FileInfo,
) (files map[string]bool, err error) {
	files = make(map[string]bool)
	for _, info := range templates {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), TemplateFileSuffix) {
			continue
		}
		log.With(plugin.Lang).Debug("template file: %s", info.Name())
		filename := filepath.Join(plugin.TemplatesDir, info.Name())
		meta, temp, err := ParseTemplateFile(filename)
		if err != nil {
			return files, err
		}

		var file *os.File
		temp = temp.Funcs(Funcs(buildType))
		kind, suffix := TemplateFileKind(info.Name())
		log.Debug("kind=%s, suffix=%s", kind, suffix)

		ctx := &Context{
			Pkg:       pkg,
			Root:      temp,
			buildType: buildType,
		}
		switch kind {
		case "package":
			dftName := pkg.Name + "." + suffix
			if file, err = ApplyMeta(outdir, meta, pkg, dftName); err == nil {
				err = temp.Execute(file, Package{Package: pkg, Context: ctx})
				files[meta.File] = true
				file.Close()
				if err != nil {
					return files, err
				}
			}
		case "file":
			for _, f := range pkg.Files {
				dftName := f.Filename + "." + suffix
				if file, err = ApplyMeta(outdir, meta, f, dftName); err == nil {
					files[meta.File] = true
					err = temp.Execute(file, File{
						File:    f,
						Context: ctx,
					})
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
						if file, err = ApplyMeta(outdir, meta, c, "constants."+suffix); err == nil {
							files[meta.File] = true
							err = temp.Execute(file, NewGenDecl(ctx, f, c))
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
						if file, err = ApplyMeta(outdir, meta, b, dftName); err == nil {
							files[meta.File] = true
							err = temp.Execute(file, NewBean(ctx, f, b))
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
