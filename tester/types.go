package tester

type Element struct {
	Body        []byte
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type Vars map[string]string

type Assert struct {
	NotEmpty []string
	Empty    []string
}

type Test struct {
	Title        Element
	Query        Element
	QueryVar     Element
	Response     Element
	ResponseVars Vars
	Assertion    []Assert
}

type Suite struct {
	Title  Element
	Tests  []Test
	VarSet Vars
}
