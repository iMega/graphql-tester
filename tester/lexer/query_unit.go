package lexer

import (
	"regexp"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type queryUnit struct {
	token int
}

func (u *queryUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---\s*query\s*\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *queryUnit) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[0], m), nil
}

func (*queryUnit) Scan(token *lexmachine.Token, s *tester.Suite) error {
	n := len(s.Tests) - 1
	q, _ := token.Value.(string)
	s.Tests[n].Query = tester.Element{
		Body:        []byte(q),
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
		EndLine:     token.EndLine,
		EndColumn:   token.EndColumn,
	}
	return nil
}

func (queryUnit) Cmd() tester.CmdFunc {
	return nil
}

func init() {
	AddUnit(&queryUnit{})
}

func (queryUnit) SetLexer(l *lexmachine.Lexer) {}
