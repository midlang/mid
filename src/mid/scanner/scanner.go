package scanner

import (
	"io"

	"github.com/midlang/mid/src/mid/external/go/scanner"
	"github.com/midlang/mid/src/mid/lexer"
)

const (
	DefaultScanMode = scanner.GoTokens
)

type Scanner struct {
	scanner.Scanner
}

func NewScanner(file *lexer.File, src io.Reader) *Scanner {
	s := new(Scanner)
	s.Scanner.Init(file, src)
	return s
}

func (s *Scanner) Scan() (pos lexer.Pos, tok lexer.Token, lit string) {
	r := s.Scanner.Scan()
	pos = lexer.Pos(s.Scanner.Pos().Offset)
	tok = lexer.EOF
	if r == scanner.EOF {
		return
	}
	lit = s.Scanner.TokenText()
	switch r {
	case scanner.Ident:
		tok = lexer.Lookup(lit)
	case scanner.Int:
		tok = lexer.INT
	case scanner.Float:
		tok = lexer.FLOAT
	case scanner.Char:
		tok = lexer.CHAR
	case scanner.String:
		tok = lexer.STRING
	case scanner.Comment:
		tok = lexer.COMMENT
	default:
		if op, ok := lexer.LookupOperator(lit); ok {
			tok = op
		} else {
			tok = lexer.ILLEGAL
		}
	}
	return
}
