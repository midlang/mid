package lexer

import (
	"strconv"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	LPAREN    // (
	RPAREN    // )
	LBRACK    // [
	RBRACK    // ]
	LBRACE    // {
	RBRACE    // }
	LESS      // <
	GREATER   // >
	COMMA     // ,
	PERIOD    // .
	SEMICOLON // ;
	COLON     // :
	ASSIGN    // =
	AT        // @ @function(args)
	DOLLAR    // $ env variable
	SHARP     // #
	operator_end

	keyword_beg
	PACKAGE  // package
	IMPORT   // import
	ENUM     // enum
	CONST    // const
	STRUCT   // struct
	PROTOCOL // protocol
	SERVICE  // service
	REQUIRED // required
	OPTIONAL // optional
	EXTENDS  // extends
	keyword_end
)

const Group = "group"

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	LPAREN:    "(",
	RPAREN:    ")",
	LBRACK:    "[",
	RBRACK:    "]",
	LBRACE:    "{",
	RBRACE:    "}",
	LESS:      "<",
	GREATER:   ">",
	COMMA:     ",",
	PERIOD:    ".",
	SEMICOLON: ";",
	COLON:     ":",
	ASSIGN:    "=",
	AT:        "@",
	DOLLAR:    "$",
	SHARP:     "#",

	PACKAGE:  "package",
	IMPORT:   "import",
	ENUM:     "enum",
	CONST:    "const",
	STRUCT:   "struct",
	PROTOCOL: "protocol",
	SERVICE:  "service",
	REQUIRED: "required",
	OPTIONAL: "optional",
	EXTENDS:  "extends",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var (
	keywords  map[string]Token
	operators map[string]Token
)

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	operators = make(map[string]Token)
	for i := operator_beg + 1; i < operator_end; i++ {
		operators[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func LookupOperator(lit string) (Token, bool) {
	if op, ok := operators[lit]; ok {
		return op, true
	}
	return ILLEGAL, false
}

func (tok Token) IsLiteral() bool  { return literal_beg < tok && tok < literal_end }
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }
func (tok Token) IsKeyword() bool  { return keyword_beg < tok && tok < keyword_end }

type BuiltinType int

const (
	Any     BuiltinType = iota // any
	Bool                       // bool
	Byte                       // byte
	Bytes                      // bytes
	String                     // string
	Int                        // int
	Int8                       // int8
	Int16                      // int16
	Int32                      // int32
	Int64                      // int64
	Uint                       // uint
	Uint8                      // uint8
	Uint16                     // uint16
	Uint32                     // uint32
	Uint64                     // uint64
	Float32                    // float32
	Float64                    // float64
	Map                        // map<K,V>
	Vector                     // vector<T>
	Array                      // array<T,Size>
)

var builtinTypes = [...]string{
	Any:     "any",
	Bool:    "bool",
	Byte:    "byte",
	Bytes:   "bytes",
	String:  "string",
	Int:     "int",
	Int8:    "int8",
	Int16:   "int16",
	Int32:   "int32",
	Int64:   "int64",
	Uint:    "uint",
	Uint8:   "uint8",
	Uint16:  "uint16",
	Uint32:  "uint32",
	Uint64:  "uint64",
	Float32: "float32",
	Float64: "float64",
	Map:     "map",
	Vector:  "vector",
	Array:   "array",
}

var revBuiltinTypes = make(map[string]BuiltinType)

func init() {
	for t, s := range builtinTypes {
		revBuiltinTypes[s] = BuiltinType(t)
	}
}

func LookupType(ident string) (BuiltinType, bool) {
	t, ok := revBuiltinTypes[ident]
	return t, ok
}

func (bt BuiltinType) String() string {
	if int(bt) >= 0 && int(bt) < len(builtinTypes) {
		return builtinTypes[int(bt)]
	}
	return "<unknown>"
}

func (bt BuiltinType) IsInt() bool {
	switch bt {
	case Byte,
		Int,
		Int8,
		Int16,
		Int32,
		Int64,
		Uint,
		Uint8,
		Uint16,
		Uint32,
		Uint64:
		return true
	}
	return false
}

func (bt BuiltinType) IsFloat() bool     { return bt == Float32 || bt == Float64 }
func (bt BuiltinType) IsNumber() bool    { return bt.IsInt() || bt.IsFloat() }
func (bt BuiltinType) IsContainer() bool { return bt == Array || bt == Map || bt == Vector }
