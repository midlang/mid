package genutil

import (
	"strconv"
)

type Bool bool

func (b Bool) Get() bool          { return bool(b) }
func (b *Bool) Set(v bool) string { *b = Bool(v); return "" }
func (b Bool) String() string     { return strconv.FormatBool(bool(b)) }

type String string

func (s String) Get() string          { return string(s) }
func (s *String) Set(v string) string { *s = String(v); return "" }
func (s String) String() string       { return string(s) }

type Int int64

func (i Int) Get() int64              { return int64(i) }
func (i Int) Int() int                { return int(i) }
func (i *Int) Set(v int64) string     { *i = Int(v); return "" }
func (i *Int) Add(delta int64) int64  { return i.Get() + delta }
func (i *Int) Sub(delta int64) int64  { return i.Get() - delta }
func (i *Int) Incr(delta int64) int64 { i.Set(i.Get() + delta); return i.Get() }
func (i Int) String() string          { return strconv.FormatInt(int64(i), 10) }
