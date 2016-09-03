midlang 设计规范
==============

`midlang` 被设计成算法和数据描述语言，不关心操作系统环境。midlang 不能直接运行，而只能翻译成其他语言，比如 c/c++/go/rust/java/javascript/swift/python/php/ruby...。
`midlang` 主要用于代码的跨语言共享，避免基本算法用各种语言重复写。

* midlang 的特性非常简单而且目前不支持多线程。
* midlang 中不能操作文件，网络，键盘，鼠标之类的设备。
* midlang 不支持范型，除了内置类型 `array`，`vector`，`map` 以外。
* midlang 对 c 的支持有些限制(具体什么限制后面单独说明)。
* midlang 本身没有运行时和垃圾回收的概念，所以 midlang 中离开作用域的变量都无效。
* midlang 暂时不支持类型自动推导。
* midlang 不支持反射。但是可以定义注解，那么在支持反射的语言中就可以使用此注解。
* midlang 不支持无符号整型，因为有些语言不支持，比如java。
* midlang 中可以定义单元测试，目前是被翻译成 go 来执行的测试(前面说过，midlang 必须翻译成别的语言才能执行)。
* midlang 不能调用其他语言的代码，只能生成代码来被其他语言调用
* midlang 提供内置的标准库，标准库都已经编译成了其他语言以供引入(include/import...)

应当牢记：midlang 是一种面向数据和算法的抽象描述语言，她不同于一般的编程语言具有运行的能力，也不像 protobuf 那样的语言只有数据描述能力。

## 关键字
---

	byte
	int
	int16
	int32
	int64
	string
	array
	vector
	map
	struct
	enum
	interface
	func
	void
	public
	private
	ref
	let
	var
	if
	else
	for
	in
	break
	continue
	return
	switch
	case
	package
	import

自定义关键字: 通过 midlang 编译器 `zcc` 的 `alias` 参数或配置可以给 `struct`, `enum`, `interface` 定义别名，这些别名也将成为关键字

## 类型系统
---

1. 基础类型: `byte`,`int`,`int16`,`int32`,`int64`,`string`,`array`,`vector`,`map`

2. 复杂类型: 
	* `struct` - 结构体, 可以拥有成员方法
	* `enum` - 枚举，只能是整数枚举或字符串枚举
	* `interface` - 接口，定义一组方法
	* 以上类型均可以在编译器中进行别名定义，比如 `alias protocol=struct`,`alias service=interface`

3. 引用,常量,变量: `ref`,`let`,`var`

4. 函数/成员方法

## 语法规范
---

### 符号
	=
	>
	<
	!=
	&&
	||
	!
	&
	|
	^
	>>
	<<
	(
	)
	{
	}
	[
	]
	"
	'
	.
	*
	/
	+
	-
	%

### 声明

* 常量: `let <Name> <Type> = <Value>`
* 变量: `var <Name> <Type> [= <Value>]`
* struct: `struct <Name> { [Fields] }`
* enum: `enum <Name> { [KeyValues] }`
* interface: `interface <Name> { <Methods> }`
* func: `[const] [public | private] func <FuncName>([Arguments]) [ReturnType] { [FuncBody] }`

### 表达式

* 括号表达式: `()`,`[]`,`{}`,`<>`
* 一元运算符号: `+`,`-`,`++`,`--`,`*`,`&`
* 二元运算符号: `+`,`-`,`*`,`/`,`%`,`&`,`>`,`>=`,`<`,`<=`,`==`,`!=`
* 赋值: `=`

### 注释和标记

注释支持 `/**/`,`//`,`#`。

* 多行注释: 开始结束对 `/*` `*/`
* 单行注释: 行中遇到 `#` 或 `//` 则表示该行后面为注释

标记使用 `@` 符号，标记的使用格式为 `@<Name>([Arguments])`，比如 `@url(https://github.com/midlang)`
