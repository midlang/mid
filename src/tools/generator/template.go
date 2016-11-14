package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mkideal/log"
)

const (
	// IncludesDir contains template files for `Include` function
	IncludesDir = "includes"

	// TemplateFileSuffix is the template filename suffix could be recognized
	TemplateFileSuffix = ".temp"
)

// Funcs holds all shared template functions
var Funcs = template.FuncMap{
	"include": func(filename string) (string, error) {
		// TODO: implements `include`
		return "", nil
	},
	"title": func(s string) string { return strings.Title(s) },
}

// Template wraps template.Template
type Template struct {
	*template.Template
}

// NewTemplate creates a Template by template.Template
func NewTemplate(temp *template.Template) *Template {
	log.Debug("NewTemplate: %s", temp.Name())
	temp = temp.Funcs(Funcs)
	t := &Template{temp}
	return t
}

// Lookup overrides template.Template.Lookup method
func (temp *Template) Lookup(name string) *Template {
	sub := temp.Template.Lookup(name)
	if sub == nil {
		return nil
	}
	return NewTemplate(sub)
}

// OpenTemplatesDir opens a directory for getting all files and directories
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
		return nil, err
	}
	templates := make([]os.FileInfo, 0, len(infos))
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), TemplateFileSuffix) {
			continue
		}
		templates = append(templates, info)
	}
	if len(templates) == 0 {
		log.With(lang).Warn("no templates found")
		return nil, nil
	}
	return templates, nil
}

// TemplateMeta represents meta information of a template file
type TemplateMeta struct {
	Dir    string
	File   string
	Date   string
	Values map[string]string
}

// ParseTemplateFile parses template file
func ParseTemplateFile(filename string) (*TemplateMeta, *Template, error) {
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
		log.Warn("ParseTemplateFile %s: %v", filename, err)
		return nil, nil, err
	}
	return meta, NewTemplate(temp), err
}

// ApplyMeta creates a target file by the template meta
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

// ParseTemplateFilename parses template filename: <kind>[.suffix][.flags].temp
// `kind` maybe `package`,`file`,`const` and other bean kinds like `struct`,`protocol`,`service` etc.
// Examples:
//	struct.go.temp    -> (struct, go)
//	struct.h.temp     -> (struct, h)
//	struct.cpp.temp   -> (struct, cpp)
//	struct.cpp.1.temp -> (struct, cpp)
func ParseTemplateFilename(filename string) (kind, suffix string) {
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
