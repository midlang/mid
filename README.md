midlang
=======

What's midlang?
---------------

The aim of the first stage is to implement a `Data Define Language` like `protobuf`, but have some differences.

1.	`midlang` generated code is highly customizable. Using go template to generate your codes, even documents.
2.	`midlang` was committed to eradicating boring and tedious code which can be generated, not just as a data interchange format.

The compiler `midc` compile midlang source code to an AST, and then you can visit the AST in template file.

Here is an example template file (`package.go.temp`\):

```markdown
package {{.Name}}

{{define "T_const"}}
{{.Doc}}const (
	{{range $field := .Consts}}{{$field.Name}} = {{$field.ValueString}}
	{{end}}
)
{{end}}

{{define "T_enum"}}
{{$type := .Name}}
type {{$type}} int

{{.Doc}}const (
	{{range $field := .Fields}}{{$type}}_{{$field.Name}} = {{$field.Value}}{{$field.Comment}}
	{{end}}
)
{{end}}

{{define "T_struct"}}
{{$type := .Name}}
{{.Doc}}type {{$type}} struct {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
	{{end}}
}
{{end}}

{{define "T_protocol"}}
{{template "T_struct" .}}
{{end}}

{{define "T_service"}}
{{$type := .Name}}
{{.Doc}}type {{$type}} interface {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
	{{end}}
}
{{end}}

{{.GenerateDeclsBySubTemplates}}
```

The following will be devoted to writing a document on how to use template in midlang.

Mainly includes the following points:

-	template filename form: \<kind\>\[.suffix\][.flags].temp

	-	kind maybe `package`,`file`,`const` and other bean kinds like `struct`,`protocol`,`service` etc.
	-	struct.go.temp -> (struct, go)
	-	struct.h.temp -> (struct, h)
	-	struct.cpp.temp -> (struct, cpp)
	-	struct.cpp.1.temp -> (struct, cpp)

-	template kind `package`: each package as a data be applied to the template file

-	template kind `file`: each source file as a data be applied to the template file

-	template kind `const`: each constant group as a data be applied to the template file

-	template kind `enum`: each enum type as a data be applied to the template file

-	template kind `struct`: each struct as a data be applied to the template file

-	template kind `protocol`: each protocol as a data be applied to the template file

-	template kind `service`: each service as a data be applied to the template file

-	template can use some builtin functions, e.g. `context`,`include`,`include_template`,`osenv` and many utility string functions.(all these functions are defined in `mid/src/tools/generator/generator.go`\)

-	The `context` function returns a object `Context` which has fields `Pkg`, `Plugin`, `Config` and methods `BuildType`, `Env`, `FindBean` etc.

See [builtin templates](https://github.com/midlang/mid/tree/master/templates)

Install
-------

### install from binary package

-	Download binary package from here:

[http://midlang.org/dl](http://midlang.org/dl)

-	Unpack \*.tar.gz, and then `cd` target directory

```sh
$ tar zcf <name>.tar.gz
$ cd <name>
```

-	Install to your local host

1. Copy all files in directory `bin` to any directory which contained in env `PATH`
2. Copy file `midconfig` your home directory or one of these: `/etc`,`/usr/local/etc`
3. Copy directories `templates` and `extensions` to $HOME/.mid or $MIDROOT if exist

### install from source

```sh
$ go get github.com/midlang/mid
$ cd /path/to/mid # replace `/path/to/mid` with your actual directory
$ ./install.sh
```

Getting started
---------------

### Write a source file `demo.mid`

```mid
/**
 * package declaration must has a `;`
 */
package demo;

// constants
const (
	A = 1;
	B = 2;
)

// doc: status
// Status
enum Status {
	Ok = 0, // ok
	Bad = 1, // bad
}

struct User {
	int64 id;
	string name;
	vector<string> otherNames;
	array<byte,6> code;
}

protocol UserList {
	map<int64,User> users;
}

service UserService {
	sayHello()
	getUsers() UserList
	findUser(int64 uid) User
	delUser(int64) Status
}

// extends another structure
struct Info extends User {
	string desc;
	map<int64,vector<map<int,array<bool,5>>>> xxx;
}
```

### Compile and generate codes

```sh
$ midc -I demo.mid -Ogo=generated/go
# Or
# midc -I demo.mid -Ogo=generated/go -K default
```

Try `midc -h` to get help information of compiler `midc`

Try this:

```sh
$ midc -I demo.mid -Ogo=generated/go -K beans
```

Language plugins
----------------

Midlang language plugin used to generate code for the program language

Here are all builtin plugins:

-	gengo - golang plugin
-	gencpp - c++ plugin

You can write yourself plugin instead of using builtin plugin.

Templates
---------

Templates used to generate codes. You can use option `-T<lang>=<template_dir>` to specify templates for specified language, also you can use option `-K <template_kind>` to specify template kind(default template kind is default).

Builtin template kinds: `default`, `beans`

Extensions
----------

Extension is a plugin for builtin template. You can insert code in some specified positions like:

* file_head
* before_import
* in_import
* after_import
* before_const
* after_const
* const_front
* const_back
* before_enum
* enum_front
* enum_back
* after_enum
* before_struct
* struct_front
* struct_back
* after_struct
* before_protocol
* protocol_front
* protocol_back
* after_protocol
* before_service
* service_front
* service_back
* after_service

See [builtin extensions](https://github.com/midlang/mid/tree/master/extensions)

Plugin for editor
-----------------

-	vim-mid - [https://github.com/midlang/vim-mid](https://github.com/midlang/vim-mid)

TODO
----

-	Other generators: java,rust,swift,python,javascript,...
-	More builtin extensions
-	Build documents
