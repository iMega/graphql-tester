package main

import (
	"fmt"
	"log"

	"github.com/imega/graphql-tester/cmd"
	"github.com/imega/graphql-tester/tester"
	"github.com/imega/graphql-tester/tester/lexer"
)

func main() {
	cmd.Execute()

	buf := make(chan tester.MessageCh)
	go tester.PrinterWatch(buf)

	opts := tester.Options{
		URL:     cmd.Options.URL,
		Headers: cmd.Options.Headers,
		Path:    cmd.Options.Path,
		StdOut:  buf,
		Verbose: false,
	}

	scan, err := lexer.Compile()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	if err := tester.RunNew(opts, scan); err != nil {
		fmt.Printf("failed to run tester, %s", err)
	}
}
