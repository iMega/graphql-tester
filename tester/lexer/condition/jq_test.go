package condition

import (
	"testing"

	lm "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

func Test_Lexer_ConditionJq_Action(t *testing.T) {
	jqAction(t, []byte("jq data.value\n"), "data.value")
	jqAction(t, []byte("jq   data.value\n"), "data.value")
	jqAction(t, []byte("jq   data.value   \n"), "data.value")

	jqAction(t, []byte("jq data.value|"), "data.value")
	jqAction(t, []byte("jq   data.value  |"), "data.value")

	jqAction(t, []byte("jq     data.value  "), "data.value")
	jqAction(t, []byte("jq     data.value"), "data.value")
	jqAction(t, []byte("jq data.value"), "data.value")

	jqAction(t, []byte("jq data.value.[0].field\n"), "data.value.[0].field")
	jqAction(t, []byte("jq   data.value.[0].field\n"), "data.value.[0].field")
	jqAction(t, []byte("jq   data.value.[0].field   \n"), "data.value.[0].field")

	jqAction(t, []byte("jq data.value.[0].field|"), "data.value.[0].field")
	jqAction(t, []byte("jq   data.value.[0].field  |"), "data.value.[0].field")

	jqAction(t, []byte("jq     data.value.[0].field  "), "data.value.[0].field")
	jqAction(t, []byte("jq     data.value.[0].field"), "data.value.[0].field")
	jqAction(t, []byte("jq data.value.[0].field"), "data.value.[0].field")

	jqAction(t, []byte("jq data.value.[0:5].field\n"), "data.value.[0:5].field")
	jqAction(t, []byte("jq   data.value.[0:5].field\n"), "data.value.[0:5].field")
	jqAction(t, []byte("jq   data.value.[0:5].field   \n"), "data.value.[0:5].field")

	jqAction(t, []byte("jq data.value.[0:5].field|"), "data.value.[0:5].field")
	jqAction(t, []byte("jq   data.value.[0:5].field  |"), "data.value.[0:5].field")

	jqAction(t, []byte("jq     data.value.[0:5].field  "), "data.value.[0:5].field")
	jqAction(t, []byte("jq     data.value.[0:5].field"), "data.value.[0:5].field")
	jqAction(t, []byte("jq data.value.[0:5].field"), "data.value.[0:5].field")
}

func Test_Lexer_ConditionJq_Action_NoMatch_ReturnsError(t *testing.T) {
	jq := conditionJq{token: 0}
	scanner := &lm.Scanner{}
	match := &machines.Match{
		Bytes: []byte(`data.value`),
	}
	_, err := jq.action(scanner, match)
	if err == nil {
		t.Fatalf("must return error")
	}
}

func Test_Lexer_ConditionJq_ActionFunc(t *testing.T) {
	unit := conditionJq{}
	lex := lm.NewLexer()
	lex.Add(unit.ActionFunc(0))

	if err := lex.CompileDFA(); err != nil {
		t.Fatalf("failed to compile, %s", err)
	}

	jqScanner(t, lex, []byte("jq data.value\n"), []string{"data.value"})
	jqScanner(t, lex, []byte("jq   data.value  \n"), []string{"data.value"})
	jqScanner(t, lex, []byte("jq data.value |jq data\n"), []string{"data.value", "data"})
	jqScanner(t, lex, []byte("jq data.value | jq data\n"), []string{"data.value", "data"})
	jqScanner(t, lex, []byte("jq   data.value   |    jq data\n"), []string{"data.value", "data"})
	jqScanner(t, lex, []byte("jq data.value"), []string{"data.value"})
	jqScanner(t, lex, []byte("jq   data.value.[0:1].test | jq data.value\n"), []string{"data.value.[0:1].test", "data.value"})
}

func Test_Lexer_ConditionJq_Cmd_ReturnsNoError(t *testing.T) {
	jq := conditionJq{token: 0}

	f := jq.Cmd()
	i, err := f([]byte(`{
		"data": {
			"viewer": {
				"login": "iMega"
			}
		}
	}`), "data.viewer.login")
	if err != nil {
		t.Fatalf("failed to execute command, %s", err)
	}

	s, ok := i.(string)
	if !ok {
		t.Fatalf("failed to converting interface")
	}

	if s != "iMega" {
		t.Fatalf("failed to equal")
	}
}

func Test_Lexer_ConditionJq_Cmd_InvalidSelector_ReturnsError(t *testing.T) {
	jq := conditionJq{token: 0}

	f := jq.Cmd()
	_, err := f([]byte(`{
		"data": {
			"viewer": {
				"login": "iMega"
			}
		}
	}`), "data.viewer.invalid_selector")
	if err == nil {
		t.Fatalf("must return error")
	}
}

func jqAction(t *testing.T, in []byte, expected string) {
	jq := conditionJq{token: 0}
	scanner := &lm.Scanner{}

	match := &machines.Match{
		Bytes: in,
	}
	i, err := jq.action(scanner, match)
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

func jqScanner(t *testing.T, lex *lm.Lexer, in []byte, expected []string) {
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
