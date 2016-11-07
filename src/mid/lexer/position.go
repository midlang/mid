// NOTE: source from go std package "go/token"

package lexer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mkideal/log"
)

type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}

// A source position is represented by a Position value.
// A position is valid if Line > 0.
type Position struct {
	Filename string // filename, if any
	Offset   int    // byte offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count per line)
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos Position) String() string {
	s := pos.Filename
	if s == "" {
		s = "<input>"
	}
	if pos.IsValid() {
		s += fmt.Sprintf(":%d:%d", pos.Line, pos.Column)
	}
	return s
}

type File struct {
	set  *FileSet
	name string
	base int
	size int

	lines []int
	infos []lineInfo
}

func (f *File) Name() string { return f.name }
func (f *File) Base() int    { return f.base }
func (f *File) Size() int    { return f.size }

func (f *File) LineCount() int {
	f.set.mutex.RLock()
	n := len(f.lines)
	f.set.mutex.RUnlock()
	return n
}

func (f *File) AddLine(offset int) {
	f.set.mutex.Lock()
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
	f.set.mutex.Unlock()
}

func (f *File) MergeLine(line int) {
	if line <= 0 {
		panic("illegal line number (line numbering starts at 1)")
	}
	f.set.mutex.Lock()
	defer f.set.mutex.Unlock()
	if line >= len(f.lines) {
		panic("illegal line number")
	}
	copy(f.lines[line:], f.lines[line+1:])
	f.lines = f.lines[:len(f.lines)-1]
}

func (f *File) SetLines(lines []int) bool {
	size := f.size
	for i, offset := range lines {
		if i > 0 && offset <= lines[i-1] || size <= offset {
			return false
		}
	}

	f.set.mutex.Lock()
	f.lines = lines
	f.set.mutex.Unlock()
	return true
}

func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}

	f.set.mutex.Lock()
	f.lines = lines
	f.set.mutex.Unlock()
}

type lineInfo struct {
	Offset   int
	Filename string
	Line     int
}

func (f *File) AddLineInfo(offset int, filename string, line int) {
	f.set.mutex.Lock()
	if i := len(f.infos); i == 0 || f.infos[i-1].Offset < offset && offset < f.size {
		f.infos = append(f.infos, lineInfo{offset, filename, line})
	}
	f.set.mutex.Unlock()
}

func (f *File) Pos(offset int) Pos {
	if offset > f.size {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

func (f *File) Offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal Pos value")
	}
	return int(p) - f.base
}

func (f *File) Line(p Pos) int {
	return f.Position(p).Line
}

func searchLineInfos(a []lineInfo, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].Offset > x }) - 1
}

func (f *File) unpack(offset int, adjusted bool) (filename string, line, column int) {
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}
	if adjusted && len(f.infos) > 0 {
		if i := searchLineInfos(f.infos, offset); i >= 0 {
			alt := &f.infos[i]
			filename = alt.Filename
			if i := searchInts(f.lines, alt.Offset); i >= 0 {
				line += alt.Line - i - 1
			}
		}
	}
	return
}

func (f *File) position(p Pos, adjusted bool) (pos Position) {
	offset := int(p) - f.base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset, adjusted)
	return
}

func (f *File) PositionFor(p Pos, adjusted bool) (pos Position) {
	if p.IsValid() {
		if int(p) < f.base || int(p) > f.base+f.size {
			log.Fatal("illegal Pos value: p=%d, base=%d, size=%d", p, f.base, f.size)
		}
		pos = f.position(p, adjusted)
	}
	return
}

func (f *File) Position(p Pos) (pos Position) {
	return f.PositionFor(p, true)
}

type FileSet struct {
	mutex sync.RWMutex
	base  int
	files []*File
	last  *File
}

func NewFileSet() *FileSet {
	return &FileSet{
		base: 1, // 0 == NoPos
	}
}

func (s *FileSet) Base() int {
	s.mutex.RLock()
	b := s.base
	s.mutex.RUnlock()
	return b

}

func (s *FileSet) AddFile(filename string, base, size int) *File {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if base < 0 {
		base = s.base
	}
	if base < s.base || size < 0 {
		panic("illegal base or size")
	}
	// base >= s.base && size >= 0
	f := &File{s, filename, base, size, []int{0}, nil}
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("token.Pos offset overflow (> 2G of source code in file set)")
	}
	// add the file to the file set
	s.base = base
	s.files = append(s.files, f)
	s.last = f
	return f
}

func (s *FileSet) Iterate(f func(*File) bool) {
	for i := 0; ; i++ {
		var file *File
		s.mutex.RLock()
		if i < len(s.files) {
			file = s.files[i]
		}
		s.mutex.RUnlock()
		if file == nil || !f(file) {
			break
		}
	}
}

func searchFiles(a []*File, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].base > x }) - 1
}

func (s *FileSet) file(p Pos) *File {
	s.mutex.RLock()
	// common case: p is in last file
	if f := s.last; f != nil && f.base <= int(p) && int(p) <= f.base+f.size {
		s.mutex.RUnlock()
		return f
	}
	// p is not in last file - search all files
	if i := searchFiles(s.files, int(p)); i >= 0 {
		f := s.files[i]
		// f.base <= int(p) by definition of searchFiles
		if int(p) <= f.base+f.size {
			s.mutex.RUnlock()
			s.mutex.Lock()
			s.last = f // race is ok - s.last is only a cache
			s.mutex.Unlock()
			return f
		}
	}
	s.mutex.RUnlock()
	return nil
}

func (s *FileSet) File(p Pos) (f *File) {
	if p.IsValid() {
		f = s.file(p)
	}
	return
}

func (s *FileSet) PositionFor(p Pos, adjusted bool) (pos Position) {
	if p.IsValid() {
		if f := s.file(p); f != nil {
			pos = f.position(p, adjusted)
		}
	}
	return
}

func (s *FileSet) Position(p Pos) (pos Position) {
	return s.PositionFor(p, true)
}

func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
