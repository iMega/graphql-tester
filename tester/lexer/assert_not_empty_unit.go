package lexer

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type notEmpty struct {
	token int
}

func init() {
	AddUnit(&notEmpty{})
}

func (u *notEmpty) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---assert_not_empty\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *notEmpty) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	selectors, err := parsingAssertions(string(m.Bytes))
	if err != nil {
		return nil, err
	}
	return s.Token(u.token, selectors, m), nil
}

func (*notEmpty) Scan(token *lexmachine.Token, s *tester.Suite) {
	n := len(s.Tests) - 1
	q, _ := token.Value.([]string)
	s.Tests[n].Assertion = append(s.Tests[n].Assertion, tester.Assert{q, nil})
}

func parsingAssertions(body string) ([]string, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return []string{}, err
	}
	res := re.FindStringSubmatch(body)
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(res[0]))
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if len(str) > 0 {
			lines = append(lines, str)
		}
	}
	return lines, scanner.Err()
}
