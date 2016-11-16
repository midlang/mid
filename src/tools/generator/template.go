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
	"github.com/mkideal/pkg/textutil/namemapper"
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
	File   string
	Values map[string]string
}

// ParseTemplateFile parses template file
func ParseTemplateFile(filename string) (*TemplateMeta, *Template, error) {
	meta := &TemplateMeta{
		Values: make(map[string]string),
	}
	// parse template file meta info
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
	temp = temp.Funcs(Funcs)
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
		temp = temp.Funcs(Funcs)
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
