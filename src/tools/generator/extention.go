package generator

type Pattern string

func (p Pattern) Match(kind string, ctx *Context) bool {
	return false
}

type Rule struct {
	ImportedPkgs   []string
	AddedFields    map[string]string
	AddedFunctions map[string]string
	AddedMethods   map[string]string
}

type Extention struct {
	TemplateDir string
	Name        string
	Rules       map[Pattern]Rule
}
