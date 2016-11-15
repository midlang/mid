package generator

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/midlang/mid/src/mid/build"
	"github.com/mkideal/log"
)

// BuildTypeFunc is a function type which used to build `build.Type` to a string
type BuildTypeFunc func(build.Type) string

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
		log.With(plugin.Lang).Debug("template file: %s", info.Name())
		filename := filepath.Join(plugin.TemplatesDir, info.Name())
		meta, temp, err := ParseTemplateFile(filename)
		if err != nil {
			return files, err
		}
		oldMetaFile := meta.File

		var file *os.File
		kind, suffix := ParseTemplateFilename(info.Name())
		log.Debug("kind=%s, suffix=%s", kind, suffix)

		ctx := NewContext(pkg, temp, config.Envvars, buildType)
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
				meta.File = oldMetaFile
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
						meta.File = oldMetaFile
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
						meta.File = oldMetaFile
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
