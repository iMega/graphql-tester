package lexer

import (
	"regexp"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type queryVarUnit struct {
	token int
}

func (u *queryVarUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---query_var\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *queryVarUnit) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[0], m), nil
}

func (*queryVarUnit) Scan(token *lexmachine.Token, s *tester.Suite) error {
	n := len(s.Tests) - 1
	q, _ := token.Value.(string)
	s.Tests[n].QueryVar = tester.Element{
		Body:        []byte(q),
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
		EndLine:     token.EndLine,
		EndColumn:   token.EndColumn,
	}
	return nil
}

func (queryVarUnit) Cmd() tester.CmdFunc {
	return nil
}

func init() {
	AddUnit(&queryVarUnit{})
}

func (queryVarUnit) SetLexer(l *lexmachine.Lexer) {}
