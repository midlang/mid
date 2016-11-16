midlang
=======

-	中文文档

midlang 是一种中间语言,主要用于代码生成. midlang 包括编译和生成两个步骤.编译工作由编译器 midc 完成,如果有传递给 midc 语言代码生成插件则会执行相应的代码生成.

生成器做成插件式,目前已实现go代码的生成器,其它语言的生成器正在逐步添加中.生成器使用go模板来生成代码,模版目录中模板文件的组织形式将决定代码的生成的规则,具有超高的定制能力.

安装
----

```sh
> go get github.com/midlang/mid/cmd/midc
> cd github.com/midlang/mid/
> ./install.sh
```

midlang 的基本语法
------------------

一个简单的例子

```mid
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

struct Info {
	extend User; // extend extends another structure
	string desc;
	map<int64,vector<map<int,array<bool,5>>>> xxx;
}
```

-	文件后缀默认为 `.mid`
-	支持单行注释 `//` 和多行注释 `/**/`
-	文件由 `package <PackageName>;` 开头,注意这里行尾分号是必须的.
-	支持 `const` 定义,每个常量行尾必须有分号
-	支持 enum 定义,每个枚举值一行,行尾加逗号
-	支持 `struct`,`protocol` 两种结构体类型,protocol 与 struct 在语法层含义一样,它其实就是struct的别名,这通常只在生成器需要区别对待,而到底如何区别对待则根据完全代码生成的需求而定,编译阶段只不过给了struct和protocol亦不同的Kind属性.struct和protocol的每个字段行尾必须要用分号';'结束.
-	支持接口定义: service
-	支持结构体继承: 使用关键字 `extend`
-	支持的基本数据类型包括:
	-	Void // void
	-	Bool // bool
	-	Byte // byte
	-	Bytes // bytes
	-	String // string
	-	Int // int
	-	Int8 // int8
	-	Int16 // int16
	-	Int32 // int32
	-	Int64 // int64
	-	Uint // uint
	-	Uint8 // uint8
	-	Uint16 // uint16
	-	Uint32 // uint32
	-	Uint64 // uint64
	-	Map // map<K,V>
	-	Vector // vector<T>
	-	Array // array<T,Size>
