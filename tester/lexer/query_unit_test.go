package lexer

import (
	"testing"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func Test_Lexer_QueryUnit_Action(t *testing.T) {
	queryUnitAction(t, []byte("---query\nmutation{}\n"), "query\nmutation{}")
	queryUnitAction(t, []byte("---    query\nmutation{}\n"), "query\nmutation{}")
	queryUnitAction(t, []byte("---    query\nmutation\n{\n}\n"), "query\nmutation\n{\n}")
}

func Test_Lexer_QueryUnit_ActionFunc(t *testing.T) {
	unit := queryUnit{}
	lex := lexmachine.NewLexer()
	lex.Add(unit.ActionFunc(0))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	queryUnitScanner(t, lex, []byte("---query\nmutation{}\n"), "\nmutation{}\n")
	queryUnitScanner(t, lex, []byte("--- query \nmutation{}\n"), "\nmutation{}\n")
	queryUnitScanner(t, lex, []byte("--- query\nmutation{}\n"), "\nmutation{}\n")
	queryUnitScanner(t, lex, []byte("--- query\nmutation\n{\n}\n"), "\nmutation\n{\n}\n")
}

func Test_Lexer_QueryUnit_Scan(t *testing.T) {
	unit := queryUnit{}
	token := &lexmachine.Token{
		Value: "val",
	}
	suite := &tester.Suite{
		Tests: []tester.Test{
			{
				Title: tester.Element{},
			},
		},
	}
	err := unit.Scan(token, suite)
	if err != nil {
		t.Fatalf("failed to scan test, %s", err)
	}

	if len(suite.Tests) == 0 {
		t.Fatalf("failed to add test, %s", err)
	}

	if string(suite.Tests[0].Query.Body) != "val" {
		t.Fatalf("failed value in token, %s", err)
	}
}

func queryUnitAction(t *testing.T, in []byte, expected string) {
	u := testUnit{token: 0}
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
		t.Fatalf("returns not equal '%s' != '%s'", token.Value, expected)
	}
}

func queryUnitScanner(t *testing.T, lex *lexmachine.Lexer, in []byte, expected string) {
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
