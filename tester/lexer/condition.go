package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type condition struct {
	token int
	lexer *lexmachine.Lexer
}

func init() {
	AddUnit(&condition{})
}

func (u *condition) SetLexer(l *lexmachine.Lexer) {
	u.lexer = l
}

func (u *condition) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---\s*?condition\s*?\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *condition) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[0], m), nil
}

func (c *condition) Scan(token *lexmachine.Token, s *tester.Suite) error {
	n := len(s.Tests) - 1
	val, _ := token.Value.(string)

	lines := strings.Split(val, "\n")
	for _, line := range lines {
		cond, err := c.scanConds([]byte(line))
		if err != nil {
			continue
		}
		s.Tests[n].Conditions = append(s.Tests[n].Conditions, cond)
	}
	return nil
}

func (*condition) Cmd() tester.CmdFunc {
	return nil
}

func (c *condition) scanConds(in []byte) (tester.Condition, error) {
	var (
		cond  = tester.Condition{}
		order int
	)

	err := scan(c.lexer, in, func(t *lexmachine.Token) {
		val, _ := t.Value.(string)
		cond[order] = &tester.ConditionCmd{
			Cmd:   units[t.Type].Cmd(),
			Value: val,
		}
		order++
	})
	if err != nil {
		return cond, err
	}

	if len(cond) == 0 {
		return cond, fmt.Errorf("empty condition")
	}

	return cond, nil
}
