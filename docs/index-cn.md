---
layout: default
date: 2018-08-19 14:30:11 +0800
title: 初识 midlang
permalink: /cn
---

<a href="/" class="ui labeled icon mini button"><i class="hand point right icon"></i>English</a>

_已经看过本篇文章，需要了解模板 API_ ? 请前往 [API文档](/cn/api) 。

## midlang 是什么？

[midlang][mid-github] 是一个通用的领域特定语言（domain-specific language、[DSL][dsl]）。`midlang` 完成语法解析工作得到语法树，然后交给语言生成插件去生成代码，而生成代码的过程高度使用模板，从而使 `midlang` 的使用者可以方便地定制代码的生成。

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

下载安装包，然后解压安装

<div class="post-content-align">
{% include download.html %}
</div>

```sh
tar zcf <name>.tar.gz
cd <name>

# 安装到 /usr/local/
sudo ./install.sh

# 或安装到别的目录，如用户目录
# PREFIX=$HOME ./install.sh
```

## 开始使用

这里使用一个简单的 `demo` 来展示 [midlang][mid-github] 的最基本用法

###	定义 mid 文件

首先创建一个名为 `demo.mid` 的文本文件，文件内容如下

```c
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

### 使用模板生成代码

定义模板文件 `package.go.temp` 如下

{% raw %}
```
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
```
{% endraw %}

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

### 基本组成元素

#### 字面值

* 标志符: 如 `main`, `x`, `i`
* 整数: 如 `12345`
* 浮点数: 如 `123.45`
* 字符: 如 `'a'`
* 字符串: 如 `"abc"`

**注意**: [标识符][identifier]只能由下划线(`-`)，小写英文字母(`abcdefghijklmnopqrstuvwxyz`)，大写英文字母(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`)和数字(`0123456789`)组成，且第一个字符不能是数字。比如 `main`，`total2`，`hello_world`，`_hello`，`HELLO`，`helloWorld`，`_1` 等都是合法的标识符，而 `$1`，`#var`，`哈罗`，`3d` 等都是不合法的标识符。

#### 运算符

```c
( // 左圆括号
) // 有圆括号
[ // 左方括号
] // 右方括号
{ // 左花括号
} // 右花括号
< // 左尖括号(小于号)
> // 右尖括号(大于号)
, // 逗号
. // 点
; // 分号
: // 冒号
= // 等号
# // 井号
```

#### 内置数据类型

内置的数据类型包括

* 基础数据类型: `any`
，`byte`，`bytes`，`bool`，`string`，`float32`，`float64`，`int`，`int8`，`int16`，`int32`，`int64`，`uint`，`uint8`，`uint16`，`uint32`，`uint64`
* 容器数据类型: `vector`，`array`，`map`

#### 关键字

```c
const    // 常量
enum     // 枚举
extends  // 继承
group    // 分组
import   // 引入包
optional // 可选字段
package  // 定义包
protocol // 结构对象定义
required // 必填字段
service  // 接口定义
struct   // 结构对象定义
```

##### `package`: 定义包

包名需要声明在文件的顶部（即除了注释之外，包名声明必须是第一个语法节点）。包名声明必须以分号 `;` 结尾，且包名必须是一个有效的[标识符][identifier]。如

```
package main;
```

##### `import`: 引入包

##### `const`: 常量

**注意**: 常量定义目前仅支持整数类型。

定义常量是可以单行定义，如下

```c
const C = 1;
```

也可以分组定义，如下

```c
const (
    A = 1;
    B = 2;
    C = 3;
)
```

不管是单行还是分组的方式定义，每个常量末尾都需要使用分号 `;` 结束。

##### `enum`: 枚举

**注意**: 枚举定义目前仅支持整数类型。

枚举定义方式如下

```c
enum Color {
    None = 0,
    Red = 1,
    Blue = 2,
    Green = 3,
}
```

枚举类型需要定义一个名字，如上例中的 `Color`，每个枚举值结尾需要一个逗号 `,`。

##### `struct`: 结构体定义

`struct` 是 `mid` 中由使用者自定义的复杂数据类型，使用时很像 `c` 语言的定义方式。如下例

```c
struct User {
    int64 id;
    string name;
    string email;
    bool is_student;
    vector<string> addresses;
}
```

与 `c` 不同的是右花括号 `}` 后面不需要分号（但是每个字段定义后面需要分号）。

##### `protocol`: 结构体定义

实际上 `protocol` 和 `struct` 除了名称之外，完全一模一样，所以定义方式参照 `struct` 的定义说明即可。既然和 `struct` 完全一样，那为什么需要多出一个 `protocol` 关键字呢？这是由于开发中经常遇到结构体的 2 种类别。第一种仅仅就是定义一个结构，说明包含的数据有什么，它没有更高的含义，另一种则常常具有一种明显的业务含义，比如数据中的一张表的定义，在定义接口时，接口参数的定义。也就是说，当使用者在需要区别对待结构体的意义时，就可以给结构分别冠以 `struct` 和 `protocol` 来区分，而如果使用者的业务不需要区分，那么始终使用 `struct` 或 `protocol` 即可。

##### `extends`: 继承

继承语法适用于 `struct` 和 `protocol`，可以单继承也可以多继承，如

```c
struct User {
    int64 id;
    string name;
    string email;
    bool is_student;
    vector<string> addresses;
}

// 单继承
struct GameUser extends User {
    int zone;
    string nickname;
    int level;
}

struct Action {
    int action_id;
    string ip;
    int64 timestamp;
}

// 多继承
protocol Login extends User,Action {
    string ip;
    int port;
}
```

##### `optional`: 可选字段

`struct` 和 `protocol` 的字段修饰词，使用时放在字段之前，如

```proto
struct User {
	int id;
	optional string nickname;
}

protocol Login {
	string ip;
	optional int port;
}
```

**注意**: 由于 `midlang` 本身只是用于生成代码，所以 `optional` 这样的修饰词对 `mid` 而言不具有任何实际的语法作用，这个修饰词最终的作用仍然由模板决定，如果模板中根本不管这个修饰词，那么这个词就彻底的没用。这个修饰词存在的价值在于给予了结构体的字段一个属性，以供模板使用，一个使用案例就是用于生成 protobuf 代码。

##### `required`: 必填字段

`struct` 和 `protocol` 的字段修饰词，使用时放在字段之前，如

```proto
struct User {
	int id;
	required string nickname;
}

protocol Login {
	string ip;
	required int port;
}
```

##### `service`: 接口定义

接口定义用于声明一组方法，比如

```c
struct User {
    int64 id;
    string name;
    string email;
    bool is_student;
    vector<string> addresses;
}

service UserService{
    getUsers() vector<User>;
    addUser(User user) bool;
    delUser(int64 id) bool;
    findUser(int64 id) User;
}
```

##### `group`: 分组

分组本身并不是一个实体，仅用于对结构体，接口等进行更好的组织。很多时候都不需要使用 `group`，但有时候可能认为将关联性很强的结构体分组定义是很好的组织方式。

```go
group (
    struct LoginRequest {
        int id;
    }

    struct LoginResponse {
        int result;
    }
)

group (
    struct LogoutRequest {
        int id;
    }

    struct LogoutResponse {
        int result;
    }
)
```

## midc 命令行工具的使用

执行 `midc -h` 查看帮助，如下

<pre><code>midlang compiler - compile source files and generate other languages code or documents

Options:

  -h, --help                   display help information
      --suffix=SUFFIX[=.mid]   source file suffix
      --midroot[=$MIDROOT]     mid root directory
  -v, --version                display version information
  -c, --config                 config filename
      --log[=warn]             log level for debugging: trace/debug/info/warn/error/fatal
  -O, --outdir                 output directories for each language, e.g. -Ogo=dir1 -Ocpp=dir2
  -X, --extension              extensions, e.g. -Xmeta -Xcodec
  -E, --env                    custom defined environment variables
  -I, --importpath             import paths for lookuping imports
  -K, --tempkind[=default]     template kind, a directory name
  -T, --template               templates directories for each language, e.g. -Tgo=dir1 -Tjava=dir2
      --id-allocator           id allocator name and options,supported allocators: file
      --id-for                 specific bean kinds which should be allocated a id
</code></pre>

### 最常用的参数

* `-O` 指定输出目录: 使用 `-Ogo=dir1 -Ojava=dir2` 这样的格式对需要生成的语言指定输出目录。
* `-T` 指定模板目录: 使用 `-Tgo=dir1 -Tjava=dir2` 这样的格式对需要生成的语言指定模板目录。
* `-E` 自定义环境变量

### 其他相对较少使用的参数

* `-K` 指定使用内置模板
* `-I` 指定包引入的查找目录
* `-X` 使用内置扩展
* `-c` 指定配置文件
* `--log` 指定日志级别，支持 `trace/debug/info/warn/error/fatal`
* `--suffix` 指定源文件后缀名
* `--midroot` 指定 `mid` 安装根目录

## mid 模板的使用

[mid][mid-github] 使用模板来定制代码的生成，所以掌握模板的书写至关重要。目前 `mid` 使用 [go][go] 语言的[模板][go-template]语法。
`mid` 的内置代码生成插件在使用模板时有一些约定规则。

### 模板文件名规范

之前的示例中使用的模板文件名为 `package.go.temp`。在这个文件名中 `temp` 是后缀，表明这是一个模板文件，`package` 代表这个模板将应用于语法树中的所有 `package` 类型的节点，`go` 则表示输出文件的后缀。

一般地，模板文件名格式为 `<ast_node_type>.<output_file_suffix>[.extra_info].temp`，其中

* `<ast_node_type>` 为节点类型，可取的值为: `package`，`file`，`const`，`enum`，`group`，`struct`，`protocol`，`service`
* `<output_file_suffix>` 为生成文件后缀，如 `go`，`c`，`java`，`js`，`MD`，`txt` 等等
* `[.extra_info]` 为可选的额外信息，可用于助记或区分。如 `package.go.orm.temp`，`package.go.def.temp`

按照上面的说明，模板文件的命名将直接控制文件的生成。内置生成插件使用模板生成文件的过程为

1. 首先根据指定的模板目录取出该目录下所有后缀为 `.temp` 的模板文件
2. 对于每个模板文件，取出文件名指定的节点类型 `ast_node_type` 和生成文件后缀 `output_file_suffix`
3. 遍历语法树中的 `ast_node_type` 类型的所有节点，对于每个节点，传入节点到模板中执行，生成文件输出

由于同一种节点可以有多个模板文件，故而同一个节点可能输出到多个文件中。我们举个简单的例子来说明。

源文件 `demo.mid` 内容如下

```c
package model;

struct User {
    int id;
}

struct Profile {
	string email;
}
```

模板文件 `temp/package.txt.temp` 如下

{% raw %}
```
hello package {{.Name}}
```
{% endraw %}

模板文件 `temp/package.md.temp` 如下

{% raw %}
```
## {{.Name}}
```
{% endraw %}

模板文件 `temp/struct.txt.temp` 如下

{% raw %}
```
hello struct {{.Name}}
```
{% endraw %}

模板文件 `temp/struct.go.temp` 如下

{% raw %}
```
package {{context.Pkg.Name}}

type {{.Name}} struct {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}
	{{end}}
}
```
{% endraw %}

目录结构如下

```
.
├── demo.mid
└── temp
    ├── package.md.temp
    ├── package.txt.temp
    ├── struct.go.temp
    └── struct.txt.temp
```

执行命令 `midc -Ogo=out -Tgo=temp` 后的目录结构如下

```
.
├── demo.mid
├── out
│   └── model
│       ├── Profile.go
│       ├── Profile.txt
│       ├── User.go
│       ├── User.txt
│       ├── model.md
│       └── model.txt
└── temp
    ├── package.md.temp
    ├── package.txt.temp
    ├── struct.go.temp
    └── struct.txt.temp
```

可以看到新增了 6 个文件，这 6 个文件就是输出的文件。这 6 个文件的内容分别为

* model.md

```
## model
```

* model.txt

```
hello package model
```

* User.txt

```
hello struct User
```

* User.go

```go
package model

type User struct {
	Id int
}
```

* Profile.txt

```
hello struct Profile
```

* Profile.go

```go
package model

type Profile struct {
	email string
}
```

### 模板文件头元数据

在模板文件最前面使用

```
---
key1: value1
key2: value2
...
keyN: valueN
---
```

的方式可以定义模板文件的元数据，value 中可以使用模板。比如如下所示的模板文件 `struct.go.temp`

{% raw %}
```
---
file: {{.Name | toLower}}.go
author: me
date: 2008/08/08
---
package {{context.Pkg.Name}}

type {{.Name}} struct {
	{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}
	{{end}}
}
```
{% endraw %}

在这个例子中 `file` 字段是具有特殊含义的，该字段指定了文件输出的名字，这将覆盖掉默认的文件名（语法节点名+后缀），而 `author`，`data` 这两个字段不起任何作用，仅用于模板书写和观看者标注。目前有特殊用途的字段有以下几个

* `file`: 指定输出文件的名字
* `cond`: 文件输出的条件，这个值不等于 `false` 是文件会输出的必要条件之一
* `notexist`: 该字段为 `true` 时要求输出的文件当前并不存在，这可以用于控制当要生成的文件已经存在时就不再输出
* `append`: 该字段为 `true` 时文件追加输出（即如果原文件存在，则在文件后面追加输出内容而不是覆盖）

### 模板内容的书写

上面提到模板生成代码的规则是将语法树节点传入到了模板中供模板使用，那么传入的这个节点数据，模板怎么使用的呢？使用 `.` 即可，如 {% raw %}`{{.Name}}`{% endraw %} 就是获取节点的名称，`Name` 是节点的一个字段。除了 `Name` 之外，节点通常都有一些别的字段和成员方法，不同的节点类型器字段和成员方法不一样，详细的可参见 [API 文档](/cn/api)。在这里挑选几个常用的讲讲。

#### `GenerateDeclsBySubTemplates` 成员方法

这个成员方法只有 `package` 节点和 `file` 节点拥有。`package` 节点为包节点，代表一个包，`file` 节点为文件节点，代表一个 `mid` 源文件。这两种节点下面都可能有许多 `const`，`enum`，`struct`，`protocol` 之类的节点，分别遍历每种类型的节点是一个常见的需求，故而有了 `GenerateDeclsBySubTemplates` 方法。这个方法的作用就是遍历其中的各个 `const`，`enum`，`struct`，`protocol` ，`service` 节点，然后利用每种节点的子模板（如最前前面 `demo` 中的 `T_const`，`T_enum`，...等等）生成输出该节点的文本。

我们回头再看一下前面使用的一个 `demo` 模板 `package.go.temp`

{% raw %}
```
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
```
{% endraw %}

这个模板文件第一行 {% raw %}`package {{.Name}}`{% endraw %} 输出包定义，最后一行 {% raw %}`{{.GenerateDeclsBySubTemplates}}`{% endraw %} 调用 `GenerateDeclsBySubTemplates` 成员方法。而中间那些 `define` 块定义了针对各种节点的子模板，用来控制各种节点的代码输出格式，这些子模板就在 `GenerateDeclsBySubTemplates` 成员方法中遍历这些类型的子节点时被使用。

#### context 全局函数

在上面的例子中我们看到其中有个 `context.BuildType`，这个 `context` 是一个全局函数，用于获取一个 `Context` 对象，而 `BuildType` 是 `Context` 对象的一个成员方法，该成员方法用于将 `mid` 中的数据类型转换成指定语言的相应数据类型。还有很多 `context` 这样的全局函数，而 `Context` 对象也还有其他一些成员变量和成员函数，详情参见 [API 文档](/cn/api)。

### 模板语法基础

目前的模板采用 [go][go] 的 [模板语法][go-template]，对于已经熟悉使用的人来说，可以忽略这一节。这里也只是简单介绍一下，更详细的内容请参考 [go][go] 官方的[模板使用文档][go-template]

#### 注释

使用 {% raw %}`{{/*在这里写注释*/}}`{% endraw %} 的格式书写注释

#### 分支语句

{% raw %}
```
{{if condition}} do something {{end}}
{{if condition}} do something {{else}} do something2 {{end}}
{{if condition}} do something {{else if condition2}} do something2 {{end}}
{{if condition}} do something {{else if condition2}} do something2 {{else}} do something3 {{end}}
```
{% endraw %}

#### 循环语句

{% raw %}
```
{{range $v := $array}} ... {{end}}
```
{% endraw %}

例如前面的 `demo` 中所展示的

{% raw %}
```
{{range $field := .Fields}}{{$field.Name | title}} {{context.BuildType $field.Type}}{{$field.Comment}}
{{end}}
```
{% endraw %}

#### 定义子模板

{% raw %}
```
{{define "name"}} ... {{end}}
```
{% endraw %}

例如前面的 `demo` 中所展示的

{% raw %}
```
{{define "T_const"}}
{{range $decl := .}}
{{$decl.Doc}}const (
	{{range $field := $decl.Consts}}{{$field.Name}} = {{$field.ValueString}}{{$field.Comment}}
	{{end}}
)
{{end}}
{{end}}
```
{% endraw %}

**提示**: 子模板可以类似看作很多语言中函数的概念


[go]: https://golang.org/ "Go"
[go-template]: https://golang.org/pkg/text/template/ "Go template"
[dsl]: https://en.wikipedia.org/wiki/Domain-specific_language "DSL"
[identifier]: https://zh.wikipedia.org/wiki/%E6%A8%99%E8%AD%98%E7%AC%A6 "Identifier"
[mid-github]: https://github.com/midlang/mid "midlang"
