package condition

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/imega/graphql-tester/tester/lexer"
	"github.com/savaki/jq"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type conditionJq struct {
	token int
}

func init() {
	lexer.AddUnit(&conditionJq{})
}

func (u *conditionJq) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`\s*?jq[^|]+[|\n]?`), lexmachine.Action(u.action)
}

func (u *conditionJq) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`jq\s+(.*?)\s*?\|?$`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regexp, %s", err)
	}

	res := re.FindStringSubmatch(string(m.Bytes))
	if len(res) == 0 || len(res[1]) == 0 {
		return nil, fmt.Errorf("failed to match jq selector")
	}

	return s.Token(u.token, res[1], m), nil
}

func (*conditionJq) Scan(token *lexmachine.Token, s *tester.Suite) error {
	return nil
}

func (*conditionJq) Cmd() tester.CmdFunc {
	return func(in interface{}, val string) (interface{}, error) {
		a, ok := in.([]byte)
		if !ok {
			return nil, fmt.Errorf("not ok")
		}

		b, err := extractValueFromJson(val, a)
		if err != nil {
			return nil, fmt.Errorf("failed to extract selector, %s", err)
		}

		return strings.Trim(string(b), "\""), nil
	}
}

func (conditionJq) SetLexer(l *lexmachine.Lexer) {}

func extractValueFromJson(selector string, data []byte) ([]byte, error) {
	op, err := jq.Parse(selector)
	if err != nil {
		return nil, err
	}
	return op.Apply(data)
}
