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
	EXTEND   // extend
	keyword_end
)

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

	PACKAGE:  "package",
	IMPORT:   "import",
	ENUM:     "enum",
	CONST:    "const",
	STRUCT:   "struct",
	PROTOCOL: "protocol",
	SERVICE:  "service",
	REQUIRED: "required",
	OPTIONAL: "optional",
	EXTEND:   "extend",
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
	Void   BuiltinType = iota // void
	Bool                      // bool
	Byte                      // byte
	Bytes                     // bytes
	String                    // string
	Int                       // int
	Int8                      // int8
	Int16                     // int16
	Int32                     // int32
	Int64                     // int64
	Uint                      // uint
	Uint8                     // uint8
	Uint16                    // uint16
	Uint32                    // uint32
	Uint64                    // uint64
	Map                       // map<K,V>
	Vector                    // vector<T>
	Array                     // array<T,Size>
)

func LookupType(ident string) (BuiltinType, bool) {
	switch ident {
	case "", "void":
		return Void, true
	case "bool":
		return Bool, true
	case "byte":
		return Byte, true
	case "bytes":
		return Bytes, true
	case "string":
		return String, true
	case "int":
		return Int, true
	case "int8":
		return Int8, true
	case "int16":
		return Int16, true
	case "int32":
		return Int32, true
	case "int64":
		return Int64, true
	case "uint":
		return Uint, true
	case "uint8":
		return Uint8, true
	case "uint16":
		return Uint16, true
	case "uint32":
		return Uint32, true
	case "uint64":
		return Uint64, true
	case "map":
		return Map, true
	case "vector":
		return Vector, true
	case "array":
		return Array, true
	}
	return 0, false
}

func (bt BuiltinType) String() string {
	switch bt {
	case Void:
		return "void"
	case Bool:
		return "bool"
	case Byte:
		return "byte"
	case Bytes:
		return "bytes"
	case String:
		return "string"
	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case Uint:
		return "uint"
	case Uint8:
		return "uint8"
	case Uint16:
		return "uint16"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Map:
		return "map"
	case Vector:
		return "vector"
	case Array:
		return "array"
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
