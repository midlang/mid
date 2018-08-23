---
layout: default
date: 2016-12-04T12:06:24+08:00
categories: api
title: API documentation
permalink: /api
---

<a href="/cn/api" class="ui labeled icon mini button"><i class="hand point right icon"></i>中文</a>

## Template functions

### Common functions

<div class="ui styled accordion" style="width: 100%">

  <!-- context -->
  <div class="title"><h5><code><span class="function-name">context</span>()</code>
	Get context object
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{define "T_struct"}}
{{.Doc}}type {{.Name}} struct {
	{{range $field := .Fields}}{{$field.Name}} {{<span class="function-name">context</span>.BuildType $field.Type}}
	{{end}}
}
{{end}}
{% endraw %}</code></pre></div>

  <!-- error -->
  <div class="title"><h5><code><span class="function-name">error</span>(<span class="field-name">format</span> int, <span class="field-name">args</span> ...any)</code>
	Output error message
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">error</span> "Error: %s" "no such file"}}
{% endraw %}</code></pre></div>

  <!-- includeTemplate -->
  <div class="title"><h5><code><span class="function-name">includeTemplate</span>(<span class="field-name">filename</span> string, <span class="field-name">data</span> any)</code>
	Include template file
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">includeTemplate</span> "file.temp" .}}
{% endraw %}</code></pre></div>

  <!-- include -->
  <div class="title"><h5><code><span class="function-name">include</span>(<span class="field-name">filename</span> string)</code>
	Include file
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">include</span> "file.ext"}}
{% endraw %}</code></pre></div>

  <!-- isInt -->
  <div class="title"><h5><code><span class="function-name">isInt</span>(<span class="field-name">type</span> string)</code>
	Determine if type is an integer
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">isInt</span> "int"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "int8"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "uint"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "uint32"}} {{/*true*/}}
{{<span class="function-name">isInt</span> "string"}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- joinPath -->
  <div class="title"><h5><code><span class="function-name">joinPath</span>(<span class="field-name">paths</span> ...string)</code>
	Join path, like the filepath.Join function of go
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">joinPath</span> "path" "to/"}} {{/*path/to*/}}
{% endraw %}</code></pre></div>

  <!-- osenv -->
  <div class="title"><h5><code><span class="function-name">osenv</span>(<span class="field-name">key</span> string)</code>
	Get system environment variables
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">osenv</span> "HOME"}}
{% endraw %}</code></pre></div>

  <!-- outdir -->
  <div class="title"><h5><code><span class="function-name">outdir</span>()</code>
	Get the root directory of the generated file (that is, the directory specified by the -O parameter)
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{joinPath (<span class="function-name">outdir</span>) "subdir"}}
{% endraw %}</code></pre></div>

  <!-- pwd -->
  <div class="title"><h5><code><span class="function-name">pwd</span>()</code>
	Get the directory where the current template file is located
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{joinPath (<span class="function-name">pwd</span>) "subdir"}}
{% endraw %}</code></pre></div>

  <!-- slice -->
  <div class="title"><h5><code><span class="function-name">slice</span>(<span class="field-name">values</span> ...any)</code>
	Make all the parameters into a slice
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{$s := (<span class="function-name">slice</span> "abc" 123 true)}}
{{valueAt $s 0}}
{{valueAt $s 1}}
{{valueAt $s 2}}
{% endraw %}</code></pre></div>

  <!-- valueAt -->
  <div class="title"><h5><code><span class="function-name">valueAt</span>(<span class="field-name">values</span> []any, <span class="field-name">i</span> int)</code>
	Get the ith element of the array values (i starts at 0)
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{$s := (slice "abc" 123 true)}}
{{<span class="function-name">valueAt</span> $s 0}}
{{<span class="function-name">valueAt</span> $s 1}}
{{<span class="function-name">valueAt</span> $s 2}}
{% endraw %}</code></pre></div>

</div>

### String functions

<div class="ui styled accordion" style="width: 100%">

  <!-- append -->
  <div class="title"><h5><code><span class="function-name">append</span>(<span class="field-name">appended</span> string, <span class="field-name">s</span> string)</code>
	Append string (return s + appended)
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">append</span> "King" "Hello"}} {{/*HelloKing*/}}
{{title "hello" | <span class="function-name">append</span> "King"}} {{/*HelloKing*/}}
{% endraw %}</code></pre></div>

  <!-- containsAny -->
  <div class="title"><h5><code><span class="function-name">containsAny</span>(<span class="field-name">chars</span> string, <span class="field-name">s</span> string)</code>
	Check if the string s contains one of the unicode characters in chars
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">containsAny</span> "abcd" "hello"}} {{/*false*/}}
{{<span class="function-name">containsAny</span> "abcd" "bug"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- contains -->
  <div class="title"><h5><code><span class="function-name">contains</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	Check if the substr substring is included in the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">contains</span> "abcd" "bug"}} {{/*false*/}}
{{<span class="function-name">contains</span> "abcd" "helloabcde"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- count -->
  <div class="title"><h5><code><span class="function-name">count</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	Calculate how many substr substrings are contained in the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">count</span> "abc" "bug"}} {{/*0*/}}
{{<span class="function-name">contains</span> "abc" "helloabc"}} {{/*1*/}}
{{<span class="function-name">contains</span> "abc" "helloabcxxabcd"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- firstOf -->
  <div class="title"><h5><code><span class="function-name">firstOf</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	Split the string s with sep as the separator to get the first string
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">firstOf</span> "," "hello,world"}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- hasPrefix -->
  <div class="title"><h5><code><span class="function-name">hasPrefix</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	Determine if the string s has a prefix prefix
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">hasPrefix</span> "hel" "hello,world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- hasSuffix -->
  <div class="title"><h5><code><span class="function-name">hasSuffix</span>(<span class="field-name">suffix</span> string, <span class="field-name">s</span> string)</code>
  Determine if the string s has a suffix suffix
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">hasSuffix</span> "ld" "hello,world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- index -->
  <div class="title"><h5><code><span class="function-name">index</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	Get the index of the substring substr in the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">index</span> "xx" "hello,tele"}} {{/*-1*/}}
{{<span class="function-name">index</span> "el" "hello,tele"}} {{/*1*/}}
{{<span class="function-name">index</span> "llo" "hello,tele"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- joinStrings -->
  <div class="title"><h5><code><span class="function-name">joinStrings</span>(<span class="field-name">sep</span> string, <span class="field-name">strs</span> []string)</code>
	Join the strs array into a string with sep as a separator
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{$s := (slice "hello" "world")}}
{{<span class="function-name">joinStrings</span>"," $s}} {{/*hello,world*/}}
{% endraw %}</code></pre></div>

  <!-- join -->
  <div class="title"><h5><code><span class="function-name">join</span>(<span class="field-name">sep</span> string, <span class="field-name">strs</span> ...string)</code>
	Join the variable length parameter strs into a string with sep as the separator
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">join</span>"," "hello" "world"}} {{/*hello,world*/}}
{% endraw %}</code></pre></div>

  <!-- lastIndex -->
  <div class="title"><h5><code><span class="function-name">lastIndex</span>(<span class="field-name">substr</span> string, <span class="field-name">s</span> string)</code>
	Get the index of the last matched substring substr in the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">index</span> "xx" "hello,tele"}} {{/*-1*/}}
{{<span class="function-name">lastIndex</span> "el" "hello,tele"}} {{/*7*/}}
{{<span class="function-name">index</span> "llo" "hello,tele"}} {{/*2*/}}
{% endraw %}</code></pre></div>

  <!-- lastOf -->
  <div class="title"><h5><code><span class="function-name">lastOf</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	Split the string s with sep as the separator to get the last string
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">lastOf</span> "," "hello,world"}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- lowerCamel -->
  <div class="title"><h5><code><span class="function-name">lowerCamel</span>(<span class="field-name">s</span> string)</code>
	Convert the string s to lower camel-case
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">lowerCamel</span> "helloWorld"}} {{/*helloWorld*/}}
{{<span class="function-name">lowerCamel</span> "HelloWorld"}} {{/*helloWorld*/}}
{{<span class="function-name">lowerCamel</span> "hello_world"}} {{/*helloWorld*/}}
{% endraw %}</code></pre></div>

  <!-- nthOf -->
  <div class="title"><h5><code><span class="function-name">nthOf</span>(<span class="field-name">sep</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	Split the string s with sep as the separator to get the nth string (n starts at 0)
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">nthOf</span> "," 0 "hello,world"}} {{/*hello*/}}
{{<span class="function-name">nthOf</span> "," 1 "hello,world"}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- oneof -->
  <div class="title"><h5><code><span class="function-name">oneof</span>(<span class="field-name">s</span> string, <span class="field-name">set</span> ...string)</code>
	Determine if the string s is one of the string values set
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">oneof</span> "haha" "hello" "world"}} {{/*false*/}}
{{<span class="function-name">oneof</span> "hello" "hello" "world"}} {{/*true*/}}
{{<span class="function-name">oneof</span> "world" "hello" "world"}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- repeat -->
  <div class="title"><h5><code><span class="function-name">repeat</span>(<span class="field-name">count</span> int, <span class="field-name">s</span> string)</code>
	Repeat the string s count times
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">repeat</span> 2 "hello"}} {{/*hellohello*/}}
{% endraw %}</code></pre></div>

  <!-- replace -->
  <div class="title"><h5><code><span class="function-name">replace</span>(<span class="field-name">old</span> string, <span class="field-name">new</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	Replace substring old in string s with new
	</h5></div><div class="content"><p>最多替换 n 个子串，n 为 -1 时替换所有子串。Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">replace</span> "world" "king" 1 "hello,world,world"}} {{/*hello,king,world*/}}
{{<span class="function-name">replace</span> "world" "king" 2 "hello,world,world"}} {{/*hello,king,king*/}}
{{<span class="function-name">replace</span> "world" "king" -1 "hello,world,world"}} {{/*hello,king,king*/}}
{% endraw %}</code></pre></div>

  <!-- splitN -->
  <div class="title"><h5><code><span class="function-name">splitN</span>(<span class="field-name">sep</span> string, <span class="field-name">n</span> int, <span class="field-name">s</span> string)</code>
	Split the string s into at most n substrings with sep as a separator
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">splitN</span> "," 2 "hello,world,world"}} {{/*["hello","world,world"]*/}}
{{<span class="function-name">splitN</span> "," -1 "hello,world,world"}} {{/*["hello","world,world"]*/}}
{% endraw %}</code></pre></div>

  <!-- split -->
  <div class="title"><h5><code><span class="function-name">split</span>(<span class="field-name">sep</span> string, <span class="field-name">s</span> string)</code>
	Split the string s with sep as a separator
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">split</span> "," "hello,world,world"}} {{/*["hello","world","world"]*/}}
{% endraw %}</code></pre></div>

  <!-- stringAt -->
  <div class="title"><h5><code><span class="function-name">stringAt</span>(<span class="field-name">strs</span> []string, <span class="field-name">index</span> int)</code>
	Get the nth string in the string array strs
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{$s := (slice "hello" "world")}}
{{<span class="function-name">stringAt</span> $s 0}} {{/*hello*/}}
{{<span class="function-name">stringAt</span> $s 1}} {{/*world*/}}
{% endraw %}</code></pre></div>

  <!-- string -->
  <div class="title"><h5><code><span class="function-name">string</span>(<span class="field-name">data</span> any)</code>
	Convert data to a string
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">string</span> 123}} {{/*123*/}}
{% endraw %}</code></pre></div>

  <!-- substr -->
  <div class="title"><h5><code><span class="function-name">substr</span>(<span class="field-name">start</span> int, <span class="field-name">end</span> int, <span class="field-name">s</span> string)</code>
	Get the substring of the range [start,end) from the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">substr</span> 0 1 "abcdef"}} {{/*a*/}}
{{<span class="function-name">substr</span> 0 1 "abcdef"}} {{/*a*/}}
{{<span class="function-name">substr</span> 0 3 "abcdef"}} {{/*abc*/}}
{{<span class="function-name">substr</span> 1 3 "abcdef"}} {{/*bc*/}}
{{<span class="function-name">substr</span> 2 0 "abcdef"}} {{/*cdef*/}}
{{<span class="function-name">substr</span> 2 -1 "abcdef"}} {{/*cde*/}}
{% endraw %}</code></pre></div>

  <!-- title -->
  <div class="title"><h5><code><span class="function-name">title</span>(<span class="field-name">s</span> string)</code>
	Capitalize the first letter of the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">title</span> "hello"}} {{/*Hello*/}}
{% endraw %}</code></pre></div>

  <!-- toLower -->
  <div class="title"><h5><code><span class="function-name">toLower</span>(<span class="field-name">s</span> string)</code>
	Convert the string s to lowercase
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">toLower</span> "HELLO"}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- toUpper -->
  <div class="title"><h5><code><span class="function-name">toUpper</span>(<span class="field-name">s</span> string)</code>
	Convert the string s to uppercase
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">toUpper</span> "hello"}} {{/*HELLO*/}}
{% endraw %}</code></pre></div>

  <!-- trimPrefix -->
  <div class="title"><h5><code><span class="function-name">trimPrefix</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	Remove the prefix prefix of the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">trimPrefix</span> "he" "hello"}} {{/*llo*/}}
{% endraw %}</code></pre></div>

  <!-- trimSpace -->
  <div class="title"><h5><code><span class="function-name">trimSpace</span>(<span class="field-name">prefix</span> string, <span class="field-name">s</span> string)</code>
	Remove blank characters both ends of the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">trimSpace</span> "\t\nhello "}} {{/*hello*/}}
{% endraw %}</code></pre></div>

  <!-- trimSuffix -->
  <div class="title"><h5><code><span class="function-name">trimSuffix</span>(<span class="field-name">suffix</span> string, <span class="field-name">s</span> string)</code>
	Remove the suffix suffix of the string s
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">trimSuffix</span> "lo" "hello"}} {{/*hel*/}}
{% endraw %}</code></pre></div>

  <!-- underScore -->
  <div class="title"><h5><code><span class="function-name">underScore</span>(<span class="field-name">s</span> string)</code>
	Convert the string s to an underscore snake
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">underScore</span> "helloWorld"}} {{/*hello_world*/}}
{{<span class="function-name">underScore</span> "HelloWorld"}} {{/*hello_world*/}}
{% endraw %}</code></pre></div>

  <!-- upperCamel -->
  <div class="title"><h5><code><span class="function-name">upperCamel</span>(<span class="field-name">s</span> string)</code>
	Convert the string s to upper camel-case
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">upperCamel</span> "helloWorld"}} {{/*HelloWorld*/}}
{{<span class="function-name">upperCamel</span> "HelloWorld"}} {{/*HelloWorld*/}}
{{<span class="function-name">upperCamel</span> "hello_world"}} {{/*HelloWorld*/}}
{% endraw %}</code></pre></div>

</div>

### Logical functions

<div class="ui styled accordion" style="width: 100%">

  <!-- AND -->
  <div class="title"><h5><code><span class="function-name">AND</span>(<span class="field-name">bools</span> ...bool)</code>
	And operation
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">AND</span> true true}} {{/*true*/}}
{{<span class="function-name">AND</span> true false true}} {{/*false*/}}
{{<span class="function-name">AND</span> true}} {{/*true*/}}
{{<span class="function-name">AND</span> false}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- NOT -->
  <div class="title"><h5><code><span class="function-name">NOT</span>(<span class="field-name">b</span> bool)</code>
	Inverse operation
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">NOT</span> true}} {{/*false*/}}
{{<span class="function-name">NOT</span> false}} {{/*true*/}}
{% endraw %}</code></pre></div>

  <!-- OR -->
  <div class="title"><h5><code><span class="function-name">OR</span>(<span class="field-name">bools</span> ...bool)</code>
	Or operation
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">OR</span> true true}} {{/*true*/}}
{{<span class="function-name">OR</span> true false true}} {{/*true*/}}
{{<span class="function-name">OR</span> true}} {{/*true*/}}
{{<span class="function-name">OR</span> false}} {{/*false*/}}
{% endraw %}</code></pre></div>

  <!-- XOR -->
  <div class="title"><h5><code><span class="function-name">XOR</span>(<span class="field-name">bools</span> ...bool)</code>
	XOR operation
	</h5></div><div class="content"><p>Exmaples</p><pre>
<code>{% raw %}{{<span class="function-name">XOR</span> true true}} {{/*false*/}}
{{<span class="function-name">XOR</span> false false}} {{/*false*/}}
{{<span class="function-name">XOR</span> true false}} {{/*true*/}}
{{<span class="function-name">XOR</span> false true}} {{/*true*/}}
{% endraw %}</code></pre></div>

</div>
