package lexer

import (
	"regexp"

	"fmt"

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
	return []byte(`---\s*?condition\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
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

	//read lines

	err := c.scanConds([]byte(val), &s.Tests[n])
	if err != nil {
		return fmt.Errorf("failed to scan conditions, %s", err)
	}
	return nil
}

func (*condition) Cmd() tester.CmdFunc {
	return nil
}

func (c *condition) scanConds(in []byte, t *tester.Test) error {
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
		return err
	}

	t.Conditions = append(t.Conditions, cond)
	return nil
}
