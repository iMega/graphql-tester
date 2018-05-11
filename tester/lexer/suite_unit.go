package lexer

import (
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func init() {
	AddUnit(&suiteUnit{})
}

type suiteUnit struct {
	token int
}

func (u *suiteUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`#[^\n]*\n?`), lexmachine.Action(u.action)
}

func (u *suiteUnit) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	res := strings.Trim(string(m.Bytes), " \n")
	return s.Token(u.token, res, m), nil
}

func (suiteUnit) Scan(token *lexmachine.Token, s *tester.Suite) error {
	s.Title = tester.Element{
		Body:        token.Lexeme,
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
		EndLine:     token.EndLine,
		EndColumn:   token.EndColumn,
	}
	return nil
}

func (suiteUnit) Cmd() tester.CmdFunc {
	return nil
}

func (suiteUnit) SetLexer(l *lexmachine.Lexer) {}
