package lexer

import (
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type testUnit struct {
	token int
}

func (u *testUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`===[^\n]*\n`), lexmachine.Action(u.action)
}

func (u *testUnit) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	res := strings.Trim(string(m.Bytes)[3:], " \n")
	return s.Token(u.token, res, m), nil
}

func (*testUnit) Scan(token *lexmachine.Token, s *tester.Suite) error {
	q, _ := token.Value.(string)
	test := tester.Test{
		Title: tester.Element{
			Body:        []byte(q),
			StartLine:   token.StartLine,
			StartColumn: token.StartColumn,
			EndLine:     token.EndLine,
			EndColumn:   token.EndColumn,
		},
	}
	s.Tests = append(s.Tests, test)
	return nil
}

func (testUnit) Cmd() tester.CmdFunc {
	return nil
}

func init() {
	AddUnit(&testUnit{})
}

func (testUnit) SetLexer(l *lexmachine.Lexer) {}
