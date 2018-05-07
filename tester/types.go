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

type CmdFunc func(in interface{}, val string) (interface{}, error)

type ConditionCmd struct {
	Cmd   CmdFunc
	Value string
}

type Condition map[int]*ConditionCmd

type Test struct {
	Title        Element
	Query        Element
	QueryVar     Element
	Response     Element
	ResponseVars Vars
	Assertion    []Assert
	Conditions   []Condition
}

func (t *Test) HasExpectedResponseBody() bool {
	return len(t.Response.Body) > 0
}

type Suite struct {
	Title  Element
	Tests  []Test
	VarSet Vars
}
