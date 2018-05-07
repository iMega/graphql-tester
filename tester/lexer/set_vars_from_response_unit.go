package lexer

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type setVars struct {
	token int
}

func init() {
	AddUnit(&setVars{})
}

func (u *setVars) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---set_vars_from_response\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *setVars) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[0], m), nil
}

func (*setVars) Scan(token *lexmachine.Token, s *tester.Suite) error {
	n := len(s.Tests) - 1
	q, _ := token.Value.(string)
	s.Tests[n].ResponseVars = responseVars([]byte(q))
	return nil
}

func (setVars) Cmd() tester.CmdFunc {
	return nil
}

func responseVars(val []byte) tester.Vars {
	var res = make(tester.Vars, 0)
	r := bufio.NewReader(bytes.NewReader(val))
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}

		split := strings.Split(string(l), "=")
		if len(split) == 2 {
			res[strings.TrimSpace(split[0])] = "." + strings.TrimSpace(split[1])
		}
	}

	return res
}

func (setVars) SetLexer(l *lexmachine.Lexer) {}
