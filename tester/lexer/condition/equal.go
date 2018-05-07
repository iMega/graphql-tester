package condition

import (
	"fmt"
	"regexp"

	"github.com/imega/graphql-tester/tester"
	"github.com/imega/graphql-tester/tester/lexer"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type conditionEqual struct {
	token int
}

func init() {
	lexer.AddUnit(&conditionEqual{})
}

func (u *conditionEqual) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`\s*?equal[^|]+[|\n]?`), lexmachine.Action(u.action)
}

func (u *conditionEqual) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`equal\s+(.*?)\s*?\|?$`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[1], m), nil
}

func (conditionEqual) Scan(token *lexmachine.Token, s *tester.Suite) error {
	return nil
}

func (*conditionEqual) Cmd() tester.CmdFunc {
	return func(in interface{}, val string) (interface{}, error) {
		s, ok := in.(string)
		if !ok {
			return nil, fmt.Errorf("assert equal, is not string, %s", s)
		}
		if s != val {
			return nil, fmt.Errorf("assert equal, not equal %s != %v", s, val)
		}
		return "true", nil
	}
}

func (conditionEqual) SetLexer(l *lexmachine.Lexer) {}
