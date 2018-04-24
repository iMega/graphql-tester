package lexer

import (
	"regexp"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type expectedResponse struct {
	token int
}

func init() {
	AddUnit(&expectedResponse{})
}

func (u *expectedResponse) ActionFunc(token int) ([]byte, lexmachine.Action) {
	u.token = token
	return []byte(`---expected_response\n(([^\n]+)\n)+`), lexmachine.Action(u.action)
}

func (u *expectedResponse) action(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
	re, err := regexp.Compile(`\n(.*\s)+`)
	if err != nil {
		return nil, err
	}
	res := re.FindStringSubmatch(string(m.Bytes))

	return s.Token(u.token, res[0], m), nil
}

func (u *expectedResponse) Scan(token *lexmachine.Token, s *tester.Suite) {
	n := len(s.Tests) - 1
	q, _ := token.Value.(string)
	s.Tests[n].Response = tester.Element{
		Body:        []byte(q),
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
		EndLine:     token.EndLine,
		EndColumn:   token.EndColumn,
	}
}
