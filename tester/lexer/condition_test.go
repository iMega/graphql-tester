package lexer

import (
	"testing"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func Test_Lexer_Condition_ActionFunc(t *testing.T) {
	unit := &condition{}
	lex := lexmachine.NewLexer()
	lex.Add(unit.ActionFunc(0))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	type fixture struct {
		In       []byte
		Expected string
	}

	var fixtures = []fixture{
		{
			In:       []byte("---condition\njq data.value | equal imega\n"),
			Expected: "\njq data.value | equal imega\n",
		},
		{
			In:       []byte("--- condition\njq data.value | equal imega\n"),
			Expected: "\njq data.value | equal imega\n",
		},
		{
			In:       []byte("---condition\njq data.value | equal imega\njq data.value2 | equal value2\n"),
			Expected: "\njq data.value | equal imega\njq data.value2 | equal value2\n",
		},
	}

	for _, f := range fixtures {
		conditionScanner(t, lex, f.In, f.Expected)
	}
}

type testConditionUnit struct {
	token int
}

func (testConditionUnit) SetLexer(l *lexmachine.Lexer) {}

func (testConditionUnit) ActionFunc(token int) ([]byte, lexmachine.Action) {
	return []byte(`jq.*?[\||\n]?`), func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(token, "test value", m), nil
	}
}

func (testConditionUnit) Scan(token *lexmachine.Token, s *tester.Suite) error {
	return nil
}

func (testConditionUnit) Cmd() tester.CmdFunc {
	return func(in interface{}, val string) (interface{}, error) {
		return "value", nil
	}
}

func Test_Lexer_Condition_Scan(t *testing.T) {
	testUnit := &testConditionUnit{}
	lex := lexmachine.NewLexer()
	unit := &condition{
		lexer: lex,
	}
	lex.Add(unit.ActionFunc(0))
	lex.Add(testUnit.ActionFunc(1))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	token := &lexmachine.Token{
		Value: "jq data.value\njq data.value",
	}

	suite := &tester.Suite{
		Tests: []tester.Test{
			{
				Title: tester.Element{
					Body: []byte("test"),
				},
			},
		},
	}
	err := unit.Scan(token, suite)
	if err != nil {
		t.Fatalf("failed to scan, %s", err)
	}

	if len(suite.Tests[0].Conditions) != 2 {
		t.Fatal("conditions are empty")
	}
}

func conditionScanner(t *testing.T, lex *lexmachine.Lexer, in []byte, expected string) {
	scanner, err := lex.Scanner(in)
	if err != nil {
		t.Fatalf("failed to create scanner, %s", err)
	}

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err != nil {
			t.Fatalf("failed to scan, %s", err)
		}
		token := tok.(*lexmachine.Token)
		s, _ := token.Value.(string)

		if s != expected {
			t.Fatal("scanned value is not equal expected")
		}
	}
}
