package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/imega/graphql-tester/cmd"
	"github.com/imega/graphql-tester/tester"
	"github.com/imega/graphql-tester/tester/lexer"
	_ "github.com/imega/graphql-tester/tester/lexer/condition"
)

func main() {
	cmd.Execute()

	wg := sync.WaitGroup{}
	buf := make(chan tester.MessageCh)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tester.PrinterWatch(buf)
	}()

	opts := tester.Options{
		URL:     cmd.Options.URL,
		Headers: cmd.Options.Headers,
		Path:    cmd.Options.Path,
		StdOut:  buf,
		Verbose: cmd.Options.Verbosity,
	}

	scan, err := lexer.Compile()
	if err != nil {
		log.Fatalf("err: %s", err)
	}

	if err := tester.RunNew(opts, scan); err != nil {
		fmt.Printf("failed to run tester, %s", err)
	}

	wg.Wait()

	os.Exit(0)
}
