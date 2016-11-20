midlang
=======

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

1.	Copy all files in directory `bin` to any directory which contained in env `PATH`
2.	Copy file `midconfig` and directory `mid_templates` to your home directory or one of these: `/etc`,`/usr/local/etc`

### install from source

```sh
$ go get github.com/midlang/mid
$ cd /path/to/mid # replace `/path/to/mid` with your actual directory
$ ./install.sh
```

Now, the compiler `midc` and generators installed to `$GOPATH/bin`, file `midconfig` and directory `mid_templates` copied to `$HOME`

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

struct Info {
	extend User; // extend extends another structure
	string desc;
	map<int64,vector<map<int,array<bool,5>>>> xxx;
}
```

### Compile and generate codes

```sh
$ midc -I demo.pid -Ogo=generated/go
```

Try `midc -h` to get help information of compiler `midc`

Language plugins
----------------

Midlang language plugin used to generate code for the program language

Here are all builtin plugins:

-	gengo - golang plugin

You can write yourself plugin instead of using builtin plugin.

Templates
---------

Templates used to generate codes. You can use option `-T<lang>=<template_dir>` to specify templates for specified language, also you can use option `-K <template_kind>` to specify template kind(default template kind is default).

Builtin template kinds: `default`, `beans`

Plugin for editor
-----------------

-	vim-mid - [https://github.com/midlang/vim-mid](https://github.com/midlang/vim-mid)

TODO
----

-	Other generators: cpp,java,rust,swift,python,javascript,...
-	More builtin templates
-	Support extentions while generating codes
-	Build website and documents
