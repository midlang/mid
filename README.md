midlang
=======

-	[简体中文文档](./README.cn.md)

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
