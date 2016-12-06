package genutil

type Bool bool

func (b Bool) Get() bool          { return bool(b) }
func (b *Bool) Set(v bool) string { *b = Bool(v); return "" }

type String string

func (s String) Get() string          { return string(s) }
func (s *String) Set(v string) string { *s = String(v); return "" }

type Int int

func (i Int) Get() int          { return int(i) }
func (i *Int) Set(v int) string { *i = Int(v); return "" }
