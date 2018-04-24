package lexer

import (
	"github.com/imega/graphql-tester/tester"
	"github.com/timtadh/lexmachine"
)

type LexUnit interface {
	ActionFunc(token int) ([]byte, lexmachine.Action)
	Scan(token *lexmachine.Token, s *tester.Suite)
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
	return tester.Scan(scan), nil
}

func AddUnit(u LexUnit) {
	lex.Add(u.ActionFunc(idx))
	units[idx] = u
	idx++
}

func scan(in []byte, s *tester.Suite) error {
	scanner, err := lex.Scanner(in)
	if err != nil {
		return err
	}

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err != nil {
			return err
			break
		}
		token := tok.(*lexmachine.Token)
		units[token.Type].Scan(token, s)
	}
	return nil
}
