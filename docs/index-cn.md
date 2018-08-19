---
layout: default
date: 2018-08-19 14:30:11 +0800
title: 文档
permalink: /cn
---

<a href="/" class="ui labeled icon button"><i class="hand point right icon"></i>English</a>

## midlang 是什么？

[midlang][mid-github]是一个通用的领域特定语言（domain-specific language、[DSL][dsl]）。`midlang` 完成语法解析工作得到语法树，然后交给语言生成插件去生成代码，而生成代码的过程高度使用 [go][go] 语言的模板工具，从而使 `midlang` 的使用者可以方便地定制代码的生成。

`midlang` 至少可以用于以下目的:

1. 定义 `API`，生成所需语言的代码，也能用于生成 `API` 的文档。
2. 定义模型(`Model`)，生成所需语言的对象，以及相应的模型操作的成员方法等。
3. 定义模板，根据一组变量生成固定模式的代码。

## 安装 midlang

### 使用源代码进行安装

这种安装方式需要本机有 [go][go] 语言环境，使用以下命令完成安装

```sh
go get -u -v github.com/midlang/mid
cd /path/to/mid
./install.sh
```

### 使用预编译包安装

先从 [https://midlang.org/dl](https://midlang.org/dl) 下载安装包，然后解压安装包拷贝文件

```sh
tar zcf <name>.tar.gz
cd <name>
sudo cp bin/* /usr/local/bin/
sudo cp midconfig /usr/local/etc/
mkdir -p $HOME/.mid
cp -r templates $HOME/.mid/
cp -r extentions $HOME/.mid/
```

## 开始使用

这里使用一个简单的 `demo` 来展示 [midlang][mid-github] 的最基本用法

###	定义 mid 文件

首先创建一个名为 `demo.mid` 的文本文件，文件内容如下

```
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
```

### 使用代码生成模板生成代码

定义模板文件 `package.go.temp` 如下

```
{% raw %}
package {{.Name}}

{{/* 定义常量的生成模板 */}}
{{define "T_const"}}
{{range $decl := .}}
{{$decl.Doc}}const (
	{{range $field := $decl.Consts}}{{$field.Name}} = {{$field.ValueString}}{{$field.Comment}}
	{{end}}
)
{{end}}
{{end}}

{{/* 定义枚举的生成模板 */}}
{{define "T_enum"}}
{{$type := .Name}}
type {{$type}} int
{{.Doc}}const (
	{{range $field := .Fields}}{{$type}}_{{$field.Name}} {{$type}} = {{$field.Value}}{{$field.Comment}}
	{{end}}
)
{{end}}

{{/* 定义 struct 的生成模板 */}}
{{define "T_struct"}}
{{$type := .Name}}
{{.Doc}}type {{$type}} struct {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
	{{end}}
}
{{end}}

{{/* 定义 protocol 的生成模板 */}}
{{define "T_protocol"}}
{{$type := .Name}}
{{.Doc}}type {{$type}} struct {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
	{{end}}
}
{{end}}

{{/* 定义 service 的生成模板 */}}
{{define "T_service"}}
{{$type := .Name}}
{{.Doc}}type {{$type}} interface {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
	{{end}}
}
{{end}}

{{.GenerateDeclsBySubTemplates}}
{% endraw %}
```

将该文件放到 `demo.mid` 所在目录的子目录 `templates` 下，即如下的目录结构:

```
.
├── demo.mid
└── templates
    └── package.go.temp
```

然后执行以下命令生成代码

```sh
midc -I demo.mid -Ogo=generated/go -Tgo=templates
```

生成代码后的目录结构如下:

```
.
├── demo.mid
├── generated
│   └── go
│       └── demo
│           └── demo.go
└── templates
    └── package.go.temp
```

其中 `generated/go/demo/demo.go` 文件的内容为:

```go
package demo

// constants
const (
	A = 1
	B = 2
)

type Status int

// doc: status
const (
	Status_Ok  Status = 0 // ok
	Status_Bad Status = 1 // bad

)

type User struct {
	Id         int64
	Name       string
	OtherNames []string
	Code       [6]byte
}

type UserList struct {
	Users map[int64]User
}

type UserService interface {
	SayHello()
	GetUsers() UserList
	FindUser(uid int64) User
	DelUser(int64) Status
}
```

从上面的过程我们可以看出，使用 `mid` 主要有 3 个步骤:

1. 定义 `*.mid` 源文件
2. 定义生成文件的模板
3. 使用命令行工具 `midc` 完成文件的生成

下面分别围绕这 3 点来详细讲解 `mid` 的使用。

## mid 源文件语法

## mid 模板的使用

## midc 命令行工具的使用

[dsl]: https://en.wikipedia.org/wiki/Domain-specific_language "DSL"
[mid-github]: https://github.com/midlang/mid "midlang"
