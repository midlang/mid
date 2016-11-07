package main

const T_file = `{{.Header}}

package {{.Package}}

import (
	{{range $path := .Imports}}
	"{{$path}}"
	{{end}}
)

{{range $decl := .Decls}}
{{gen_decl $decl}}
{{end}}
`

const T_struct_def = `type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Doc}}{{title $field.Name}} {{$field.Type}} {{$field.Tag}} {{$field.Comment}}
	{{end}}
}`

const T_enum_def = `type {{.Name}} int

const (
	{{range $field := .Fields}}
	{{$field.Doc}}{{$field.Name}} = {{$field.Default}} {{$field.Comment}}
	{{end}}
)
`

const T_protocol_funcs = `
func (p {{.Name}}) ProtoName() string { return {{.Name}} }

{{gen_protocol_marshal_func .}}

{{gen_protocol_unmarshal_func .}}
`
