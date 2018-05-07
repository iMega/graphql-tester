package lexer

import (
	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type skipUnit struct {
	token int
}

func (u *skipUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`[\n \t]`), lexmachine.Action(u.action)
}

func (u *skipUnit) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	return nil, nil
}

func (skipUnit) Scan(*lexmachine.Token, *tester.Suite) error {
	return nil
}

func (skipUnit) Cmd() tester.CmdFunc {
	return nil
}

func init() {
	AddUnit(&skipUnit{})
}

func (skipUnit) SetLexer(l *lexmachine.Lexer) {}
