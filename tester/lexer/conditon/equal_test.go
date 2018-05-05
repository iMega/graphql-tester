package lexer

import (
	"testing"

	lm "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func Test_Lexer_ConditionEqual_Action(t *testing.T) {
	equalAction(t, []byte("equal data\n"), "data")
	equalAction(t, []byte("equal   data\n"), "data")
	equalAction(t, []byte("equal   data   \n"), "data")

	equalAction(t, []byte("equal data|"), "data")
	equalAction(t, []byte("equal   data  |"), "data")

	equalAction(t, []byte("equal     data  "), "data")
	equalAction(t, []byte("equal     data"), "data")
	equalAction(t, []byte("equal data"), "data")
}

func Test_Lexer_ConditionEqual_ActionFunc(t *testing.T) {
	eq := conditionEqual{}
	lex := lm.NewLexer()
	lex.Add(eq.ActionFunc(0))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	eqScanner(t, lex, []byte("equal data\n"), []string{"data"})
	eqScanner(t, lex, []byte("equal   data  \n"), []string{"data"})
	eqScanner(t, lex, []byte("equal data |equal data2\n"), []string{"data", "data2"})
	eqScanner(t, lex, []byte("equal data | equal data2\n"), []string{"data", "data2"})
	eqScanner(t, lex, []byte("equal   data   |    equal data2\n"), []string{"data", "data2"})
	eqScanner(t, lex, []byte("equal data"), []string{"data"})
}

func Test_Lexer_ConditionEqual_Cmd(t *testing.T) {
	equalCmd(t, "data", "data", nil)
}

func equalCmd(t *testing.T, in, val string, expected error) {
	eq := conditionEqual{token: 0}
	f := eq.Cmd()
	_, err := f(in, val)

	//fmt.Printf("%#v\n", err)
	//
	//switch err.(type) {
	//case error:
	//
	//}
	if err != expected {
		t.Fatalf("failed to execute command, %s", err)
	}
}

func equalAction(t *testing.T, in []byte, expected string) {
	eq := conditionEqual{token: 0}
	scanner := &lm.Scanner{}

	match := &machines.Match{
		Bytes: in,
	}
	i, err := eq.action(scanner, match)
	if err != nil {
		t.Fatalf("returns error, %s", err)
	}

	token, ok := i.(*lm.Token)
	if !ok {
		t.Fatal("returns not token instance")
	}

	if token.Value != expected {
		t.Fatalf("returns not equal %s != %s", token.Value, expected)
	}
}

func Test_equal_Cmd_InvalidSelector_ReturnsError(t *testing.T) {
	eq := conditionEqual{token: 0}

	f := eq.Cmd()
	_, err := f([]byte(`data`), "")
	if err == nil {
		t.Fatalf("must return error")
	}
}

func eqScanner(t *testing.T, lex *lm.Lexer, in []byte, expected []string) {
	var idx int
	scanner, err := lex.Scanner(in)
	if err != nil {
		t.Fatalf("failed to create scanner, %s", err)
	}

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err != nil {
			t.Fatalf("failed to scan, %s", err)
		}
		token := tok.(*lm.Token)
		s, _ := token.Value.(string)

		if s != expected[idx] {
			t.Fatal("scanned value is not equal expected")
		}
		idx++
	}
}
