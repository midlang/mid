---
layout: default
date: 2016-12-04T12:06:24+08:00
categories: api
title: API 文档
permalink: /cn/api
---

<a href="/api" class="ui labeled icon mini button"><i class="hand point right icon"></i>English</a>

## 模板函数

### 通用函数

<div class="ui styled accordion" style="width: 100%">

  <!-- context -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">context</span>()</code>
	获取上下文对象
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{define "T_struct"}}
{{.Doc}}type {{.Name}} struct {
	{{range $field := .Fields}}{{$field.Name}} {{<span class="function-name">context</span>.BuildType $field.Type}}
	{{end}}
}
{{end}}
{% endraw %}</code></pre></div>

  <!-- error -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">error</span>(<span class="field-name">format</span> int, <span class="field-name">args</span> ...any)</code>
	输出错误信息
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">error</span> "Error: %s" "no such file"}}
{% endraw %}</code></pre></div>

  <!-- includeTemplate -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">includeTemplate</span>(<span class="field-name">filename</span> string, <span class="field-name">data</span> any)</code>
	引入模板文件
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">includeTemplate</span> "file.temp" .}}
{% endraw %}</code></pre></div>

  <!-- include -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">include</span>(<span class="field-name">filename</span> string)</code>
	引入文件
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">include</span> "file.ext"}}
{% endraw %}</code></pre></div>

  <!-- isInt -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">isInt</span>(<span class="field-name">type</span> string)</code>
	判断 type 是否为整数
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">isInt</span> "int"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "int8"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "uint"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "uint32"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "string"}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- joinPath -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">joinPath</span>(<span class="field-name">paths</span> ...string)</code>
	拼接路径，如同 go 的 filepath.Join 函数
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">joinPath</span> "path" "to/"}} {{/*path/to*/}}
{% endraw %}</code></pre></div>

  <!-- osenv -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">osenv</span>(<span class="field-name">key</span> string)</code>
	获取系统环境变量
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">osenv</span> "HOME"}}
{% endraw %}</code></pre></div>

  <!-- outdir -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">outdir</span>()</code>
	获取生成文件的根目录(即 -O 参数指定的目录)
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{joinPath (<span class="function-name">outdir</span>) "subdir"}}
{% endraw %}</code></pre></div>

  <!-- pwd -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">pwd</span>()</code>
	获取当前模板文件所在目录
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{joinPath (<span class="function-name">pwd</span>) "subdir"}}
{% endraw %}</code></pre></div>

  <!-- slice -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">slice</span>(<span class="field-name">values</span> ...any)</code>
	将所有参数组成一个 slice（切片：变长数组） 返回
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{$s := (<span class="function-name">slice</span> "abc" 123 true)}}
{{valueAt $s 0}}
{{valueAt $s 1}}
{{valueAt $s 2}}
{% endraw %}</code></pre></div>

  <!-- valueAt -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">valueAt</span>(<span class="field-name">values</span> []any, <span class="field-name">index</span> int)</code>
	获取数组 values 的第 index 个元素（index 从 0 开始）
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{$s := (slice "abc" 123 true)}}
{{<span class="function-name">valueAt</span> $s 0}}
{{<span class="function-name">valueAt</span> $s 1}}
{{<span class="function-name">valueAt</span> $s 2}}
{% endraw %}</code></pre></div>

</div>

### 字符串处理函数

<div class="ui styled accordion" style="width: 100%">

  <!-- append -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">append</span>(<span class="field-name">appended</span> string, <span class="field-name">s</span> string)</code>
	追加字符串（返回 s + appended）
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">append</span> "King" "Hello"}} {{/*HelloKing*/}}
{{title "hello" | <span class="function-name">append</span> "King"}} {{/*HelloKing*/}}
{% endraw %}</code></pre></div>

  <!-- containsAny -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">containsAny</span>(<span class="field-name">chars</span> string, <span class="field-name">s</span> string)</code>
	检查字符串 s 中是否包含 chars 中的某一个 unicode 字符
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">containsAny</span> "abcd" "hello"}} {{/*false*/}}
{{<span class="function-name">containsAny</span> "abcd" "bug"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- contains -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">contains</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	检查字符串 s 中是否包含 substr 子串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">contains</span> "abcd" "bug"}} {{/*false*/}}
{{<span class="function-name">contains</span> "abcd" "helloabcde"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- count -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">count</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	计算字符串 s 中包含多少个 substr 子串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">count</span> "abc" "bug"}} {{/*0*/}}
{{<span class="function-name">contains</span> "abc" "helloabc"}} {{/*1*/}}
{{<span class="function-name">contains</span> "abc" "helloabcxxabcd"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- firstOf -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">firstOf</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	将字符串 s 以 sep 做分隔符分割后取得第一个字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">firstOf</span> "," "hello,world"}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- hasPrefix -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">hasPrefix</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	判断字符串 s 是否有前缀 prefix
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">hasPrefix</span> "hel" "hello,world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- hasSuffix -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">hasSuffix</span>(<span class="field-name">suffix</span> string, <span class="field-name">s</span> string)</code>
	判断字符串 s 是否有后缀 suffix
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">hasSuffix</span> "ld" "hello,world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- index -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">index</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	获得子串 substr 在字符串 s 中的索引
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">index</span> "xx" "hello,tele"}} {{/*-1*/}}
{{<span class="function-name">index</span> "el" "hello,tele"}} {{/*1*/}}
{{<span class="function-name">index</span> "llo" "hello,tele"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- joinStrings -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">joinStrings</span>(<span class="field-name">sep</span> string, <span class="field-name">strs</span> []string)</code>
	以 sep 为间隔将 strs 数组拼接成一个字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{$s := (slice "hello" "world")}}
{{<span class="function-name">joinStrings</span>"," $s}} {{/*hello,world*/}}
{% endraw %}</code></pre></div>

  <!-- join -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">join</span>(<span class="field-name">sep</span> string, <span class="field-name">strs</span> ...string)</code>
	以 sep 为间隔将变长参数 strs 拼接成一个字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">join</span>"," "hello" "world"}} {{/*hello,world*/}}
{% endraw %}</code></pre></div>

  <!-- lastIndex -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">lastIndex</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	获得最后一个匹配到的子串 substr 在字符串 s 中的索引
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">index</span> "xx" "hello,tele"}} {{/*-1*/}}
{{<span class="function-name">lastIndex</span> "el" "hello,tele"}} {{/*7*/}}
{{<span class="function-name">index</span> "llo" "hello,tele"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- lastOf -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">lastOf</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	将字符串 s 以 sep 做分隔符分割后取得最后一个字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">lastOf</span> "," "hello,world"}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- lowerCamel -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">lowerCamel</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 转换成小驼峰命名
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">lowerCamel</span> "helloWorld"}} {{/*helloWorld*/}}
{{<span class="function-name">lowerCamel</span> "HelloWorld"}} {{/*helloWorld*/}}
{{<span class="function-name">lowerCamel</span> "hello_world"}} {{/*helloWorld*/}}
{% endraw %}</code></pre></div>

  <!-- nthOf -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">nthOf</span>(<span class="field-name">sep</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	将字符串 s 以 sep 做分隔符分割后取得第 n 字符串（n 从 0 开始）
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">nthOf</span> "," 0 "hello,world"}} {{/*hello*/}}
{{<span class="function-name">nthOf</span> "," 1 "hello,world"}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- oneof -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">oneof</span>(<span class="field-name">s</span> string, <span class="field-name">set</span> ...string)</code>
	判断字符串 s 是否是字符串集 set 中的一个
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">oneof</span> "haha" "hello" "world"}} {{/*false*/}}
{{<span class="function-name">oneof</span> "hello" "hello" "world"}} {{/*true*/}}
{{<span class="function-name">oneof</span> "world" "hello" "world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- repeat -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">repeat</span>(<span class="field-name">count</span> int, <span class="field-name">s</span> string)</code>
	将字符串 s 重复 count 次
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">repeat</span> 2 "hello"}} {{/*hellohello*/}}
{% endraw %}</code></pre></div>

  <!-- replace -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">replace</span>(<span class="field-name">old</span> string, <span class="field-name">new</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	将字符串 s 中的 old 子串替换成 new
	</h5></div><div class="content"><p>最多替换 n 个子串，n 为 -1 时替换所有子串。使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">replace</span> "world" "king" 1 "hello,world,world"}} {{/*hello,king,world*/}}
{{<span class="function-name">replace</span> "world" "king" 2 "hello,world,world"}} {{/*hello,king,king*/}}
{{<span class="function-name">replace</span> "world" "king" -1 "hello,world,world"}} {{/*hello,king,king*/}}
{% endraw %}</code></pre></div>

  <!-- splitN -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">splitN</span>(<span class="field-name">sep</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	以 sep 做分隔符将字符串 s 分割成至多 n 个子串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">splitN</span> "," 2 "hello,world,world"}} {{/*["hello","world,world"]*/}}
{{<span class="function-name">splitN</span> "," -1 "hello,world,world"}} {{/*["hello","world,world"]*/}}
{% endraw %}</code></pre></div>

  <!-- split -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">split</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	以 sep 做分隔符将字符串 s 分割
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">split</span> "," "hello,world,world"}} {{/*["hello","world","world"]*/}}
{% endraw %}</code></pre></div>

  <!-- stringAt -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">stringAt</span>(<span class="field-name">strs</span> []string, <span class="field-name">index</span> int)</code>
	获取字符串数组 strs 中的第 n 个字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{$s := (slice "hello" "world")}}
{{<span class="function-name">stringAt</span> $s 0}} {{/*hello*/}}
{{<span class="function-name">stringAt</span> $s 1}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- string -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">string</span>(<span class="field-name">data</span> any)</code>
	将 data 转换成字符串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">string</span> 123}} {{/*123*/}}
{% endraw %}</code></pre></div>

  <!-- substr -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">substr</span>(<span class="field-name">start</span> int, <span class="field-name">end</span> int, <span class="field-name">s</span> string)</code>
	从字符串 s 中获取范围 [start,end) 的子串
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">substr</span> 0 1 "abcdef"}} {{/*a*/}}
{{<span class="function-name">substr</span> 0 1 "abcdef"}} {{/*a*/}}
{{<span class="function-name">substr</span> 0 3 "abcdef"}} {{/*abc*/}}
{{<span class="function-name">substr</span> 1 3 "abcdef"}} {{/*bc*/}}
{{<span class="function-name">substr</span> 2 0 "abcdef"}} {{/*cdef*/}}
{{<span class="function-name">substr</span> 2 -1 "abcdef"}} {{/*cde*/}}
{% endraw %}</code></pre></div>

  <!-- title -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">title</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 的首字母大写
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">title</span> "hello"}} {{/*Hello*/}}
{% endraw %}</code></pre></div>

  <!-- toLower -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">toLower</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 的转成小写
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">toLower</span> "HELLO"}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- toUpper -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">toUpper</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 的转成大写
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">toUpper</span> "hello"}} {{/*HELLO*/}}
{% endraw %}</code></pre></div>

  <!-- trimPrefix -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">trimPrefix</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	将字符串 s 的前缀 prefix 去除
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">trimPrefix</span> "he" "hello"}} {{/*llo*/}}
{% endraw %}</code></pre></div>

  <!-- trimSpace -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">trimSpace</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	将字符串 s 前后的空字符去除
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">trimSpace</span> "\t\nhello "}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- trimSuffix -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">trimSuffix</span>(<span class="field-name">suffix</span> string, <span class="field-name">s</span> string)</code>
	将字符串 s 的后缀 suffix 去除
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">trimSuffix</span> "lo" "hello"}} {{/*hel*/}}
{% endraw %}</code></pre></div>

  <!-- underScore -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">underScore</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 转换成下划线蛇形命名
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">underScore</span> "helloWorld"}} {{/*hello_world*/}}
{{<span class="function-name">underScore</span> "HelloWorld"}} {{/*hello_world*/}}
{% endraw %}</code></pre></div>

  <!-- upperCamel -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">upperCamel</span>(<span class="field-name">s</span> string)</code>
	将字符串 s 转换成大驼峰命名
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">upperCamel</span> "helloWorld"}} {{/*HelloWorld*/}}
{{<span class="function-name">upperCamel</span> "HelloWorld"}} {{/*HelloWorld*/}}
{{<span class="function-name">upperCamel</span> "hello_world"}} {{/*HelloWorld*/}}
{% endraw %}</code></pre></div>

</div>

### 逻辑运算函数

<div class="ui styled accordion" style="width: 100%">

  <!-- AND -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">AND</span>(<span class="field-name">bools</span> ...bool)</code>
	与运算
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">AND</span> true true}} {{/*true*/}}
{{<span class="function-name">AND</span> true false true}} {{/*false*/}}
{{<span class="function-name">AND</span> true}} {{/*true*/}}
{{<span class="function-name">AND</span> false}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- NOT -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">NOT</span>(<span class="field-name">b</span> bool)</code>
	取反运算
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">NOT</span> true}} {{/*false*/}}
{{<span class="function-name">NOT</span> false}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- OR -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">OR</span>(<span class="field-name">bools</span> ...bool)</code>
	或运算
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">OR</span> true true}} {{/*true*/}}
{{<span class="function-name">OR</span> true false true}} {{/*true*/}}
{{<span class="function-name">OR</span> true}} {{/*true*/}}
{{<span class="function-name">OR</span> false}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- XOR -->
  <div class="title"><i class="pointing right icon"></i><h5><code><span class="function-name">XOR</span>(<span class="field-name">bools</span> ...bool)</code>
	异或运算
	</h5></div><div class="content"><p>使用示例</p><pre>
<code>{% raw %}{{<span class="function-name">XOR</span> true true}} {{/*false*/}}
{{<span class="function-name">XOR</span> false false}} {{/*false*/}}
{{<span class="function-name">XOR</span> true false}} {{/*true*/}}
{{<span class="function-name">XOR</span> false true}} {{/*true*/}}
{% endraw %}</code></pre></div>

</div>
