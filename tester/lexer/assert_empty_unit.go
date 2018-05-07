package lexer

import (
	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type empty struct {
	token int
}

func init() {
	AddUnit(&empty{})
}

func (u *empty) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---assert_empty\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *empty) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	selectors, err := parsingAssertions(string(m.Bytes))
	if err != nil {
		return nil, err
	}
	return s.Token(u.token, selectors, m), nil
}

func (*empty) Scan(token *lexmachine.Token, s *tester.Suite) error {
	n := len(s.Tests) - 1
	q, _ := token.Value.([]string)
	s.Tests[n].Assertion = append(s.Tests[n].Assertion, tester.Assert{nil, q})
	return nil
}

func (*empty) Cmd() tester.CmdFunc {
	return nil
}

func (empty) SetLexer(l *lexmachine.Lexer) {}
