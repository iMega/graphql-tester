package lexer

import (
	"fmt"

	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
)

type LexUnit interface {
	ActionFunc(token int) ([]byte, lexmachine.Action)
	Scan(token *lexmachine.Token, s *tester.Suite) error
	Cmd() tester.CmdFunc
	SetLexer(l *lexmachine.Lexer)
}

var (
	lex   = lexmachine.NewLexer()
	units = make(map[int]LexUnit)
	idx   int
)

func Compile() (tester.Scan, error) {
	if err := lex.CompileDFA(); err != nil {
		return nil, err
	}
	return tester.Scan(scanner), nil
}

func AddUnit(u LexUnit) {
	lex.Add(u.ActionFunc(idx))
	u.SetLexer(lex)
	units[idx] = u
	idx++
}

func scanner(in []byte, s *tester.Suite) error {
	return scan(lex, in, func(t *lexmachine.Token) {
		units[t.Type].Scan(t, s)
	})
}

func scan(lex *lexmachine.Lexer, in []byte, receiver func(t *lexmachine.Token)) error {
	scanner, err := lex.Scanner(in)
	if err != nil {
		return fmt.Errorf("failed to create scanner, %s", err)
	}

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err != nil {
			return err
			break
		}
		token := tok.(*lexmachine.Token)
		receiver(token)
	}
	return nil
}
