package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mkideal/log"
	"github.com/mkideal/pkg/errors"
)

const (
	// IncludesDir contains template files for `include` function
	IncludesDir = "includes"

	// TemplateFileSuffix is the template filename suffix could be recognized
	TemplateFileSuffix = ".temp"
)

// Template wraps template.Template
type Template struct {
	*template.Template
}

// NewTemplate creates a Template by template.Template
func NewTemplate(temp *template.Template) *Template {
	log.Debug("NewTemplate: %s", temp.Name())
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
		return nil, errors.Throw(fmt.Sprintf("open templates directory %s error: %v", dir, err))
	}
	defer fs.Close()
	infos, err := fs.Readdir(-1)
	if err != nil {
		return nil, errors.Throw(err.Error())
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
	return templates, nil
}

// TemplateMeta represents meta information of a template file
type TemplateMeta struct {
	File   string
	Values map[string]string
}

// ParseTemplateFile parses template file
func ParseTemplateFile(filename string) (*TemplateMeta, *Template, error) {
	meta := &TemplateMeta{
		Values: make(map[string]string),
	}
	// parse template file meta header
	// e.g.
	//
	// ---
	// dir: aaa
	// file: {{.Name}}.go
	// ---
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("ParseTemplateFile %s: %v", filename, err)
		return nil, nil, err
	}
	advance, token, _ := bufio.ScanLines(data, true)
	const metaHeaderFlag = "---"
	if string(token) == metaHeaderFlag {
		ended := false
		line := 1
		for advance < len(data) {
			tmp, tok, _ := bufio.ScanLines(data[advance:], true)
			advance += tmp
			if tmp == 0 {
				break
			}
			if string(tok) == metaHeaderFlag {
				ended = true
				break
			}
			line++
			kv := strings.SplitN(string(tok), ":", 2)
			if len(kv) != 2 {
				err = fmt.Errorf("%s:%d: not a key value pair split by `:`", filename, line)
				return nil, nil, err
			}
			kv[0] = strings.TrimSpace(kv[0])
			meta.Values[kv[0]] = kv[1]
			log.Debug("%s:%d: key value pair: <%s:%s>", filename, line, kv[0], kv[1])
		}
		if !ended {
			err = fmt.Errorf("%s: unexpected meta header end", filename)
			return nil, nil, err
		}
		data = data[advance:]
	}
	temp := template.New(filename)
	temp = temp.Funcs(funcs)
	temp, err = temp.Parse(string(data))
	if err != nil {
		err = fmt.Errorf("ParseTemplateFile %s: %v", filename, err)
		return nil, nil, err
	}
	return meta, NewTemplate(temp), err
}

// ApplyMeta creates a target file by the template meta
func ApplyMeta(outdir string, meta *TemplateMeta, data interface{}, dftName string) (*os.File, error) {
	// execute template for meta
	values := make(map[string]string)
	for k, v := range meta.Values {
		temp := template.New(k)
		temp = temp.Funcs(funcs)
		temp, err := temp.Parse(v)
		if err != nil {
			log.Error("ApplyMeta: %v", err)
			return nil, err
		}
		var buf bytes.Buffer
		if err = temp.Execute(&buf, data); err != nil {
			log.Error("ApplyMeta: %v", err)
			return nil, err
		}
		v = strings.TrimSpace(buf.String())
		values[k] = v
		log.Debug("meta key value pair: <%s,%s>", k, v)
	}
	meta.Values = values

	// pick `file` value to meta.File
	if value, ok := meta.Values["file"]; ok {
		meta.File = value
		delete(meta.Values, "file")
	}
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
