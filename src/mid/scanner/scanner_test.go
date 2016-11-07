package scanner

import (
	"fmt"
	"strings"
	"testing"

	"github.com/midlang/mid/src/mid/external/go/scanner"
	"github.com/midlang/mid/src/mid/lexer"
)

func TestScanner(t *testing.T) {
	text := fmt.Sprintf(`package scanner
	
const x = 1
const y = "2"

const (
	a = 1.1
	b = 'c'
)

const rawstr = %s

// comment

/**
 * comments
 */

protocol Hello {
	string name;
}
`, "`abcdef`")

	src := strings.NewReader(text)

	fmt.Println(text)

	fset := lexer.NewFileSet()
	file := fset.AddFile("demo.mid", -1, len(text))
	s := NewScanner(file, src)
	s.Mode = scanner.GoTokens & (^(scanner.SkipComments))
	s.Filename = file.Name()
	fmt.Println("s.Mode&ScanRawStrings != 0:", s.Mode&scanner.ScanRawStrings != 0)
	for {
		pos, tok, lit := s.Scan()
		if tok == lexer.EOF {
			break
		}
		fmt.Printf("pos,tok,lit=(%v,%v,%v)\n", pos, tok, lit)
	}
}
