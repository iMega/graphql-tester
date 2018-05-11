package lexer

import (
	"testing"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func Test_Lexer_SuiteUnit_Action(t *testing.T) {
	suiteUnitAction(t, []byte("#SuiteName\n"), "#SuiteName")
	suiteUnitAction(t, []byte("#    SuiteName\n"), "#    SuiteName")
	suiteUnitAction(t, []byte("#    SuiteName    \n"), "#    SuiteName")
}

func Test_Lexer_SuiteUnit_ActionFunc(t *testing.T) {
	unit := suiteUnit{}
	lex := lexmachine.NewLexer()
	lex.Add(unit.ActionFunc(0))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	suiteUnitScanner(t, lex, []byte("#SuiteName\n"), "#SuiteName")
	suiteUnitScanner(t, lex, []byte("#   SuiteName\n"), "#   SuiteName")
	suiteUnitScanner(t, lex, []byte("#   Suite Name    \n"), "#   Suite Name")
}

func Test_Lexer_SuiteUnit_Scan(t *testing.T) {
	unit := suiteUnit{}
	token := &lexmachine.Token{
		Lexeme: []byte("val"),
	}
	suite := &tester.Suite{}
	err := unit.Scan(token, suite)
	if err != nil {
		t.Fatalf("failed to scan test, %s", err)
	}

	if string(suite.Title.Body) != "val" {
		t.Fatalf("failed value in token, %s", err)
	}
}

func suiteUnitAction(t *testing.T, in []byte, expected string) {
	u := suiteUnit{token: 0}
	scanner := &lexmachine.Scanner{}

	match := &machines.Match{
		Bytes: in,
	}
	i, err := u.action(scanner, match)
	if err != nil {
		t.Fatalf("returns error, %s", err)
	}

	token, ok := i.(*lexmachine.Token)
	if !ok {
		t.Fatal("returns not token instance")
	}

	if token.Value != expected {
		t.Fatalf("returns not equal '%s' expected '%s'", token.Value, expected)
	}
}

func suiteUnitScanner(t *testing.T, lex *lexmachine.Lexer, in []byte, expected string) {
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
			t.Fatalf("scanned value is not equal '%s' expected '%s'", s, expected)
		}
	}
}
