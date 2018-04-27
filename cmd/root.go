package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type OptionsCmd struct {
	URL       string
	Path      []string
	Headers   map[string]string
	Verbosity bool
}

var (
	Options OptionsCmd

	url        string
	headersOpt string
	verbosity  bool

	rootCmd = &cobra.Command{
		Use:     "GraphQL Tester",
		Version: getVersion(),
		Run: func(cmd *cobra.Command, args []string) {
			Options = OptionsCmd{
				URL:       strings.Trim(url, "' \""),
				Path:      args,
				Headers:   parseHeaders(headersOpt),
				Verbosity: verbosity,
			}
		},
	}
)

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "usage")
	rootCmd.PersistentFlags().StringVarP(&headersOpt, "headers", "H", "", "usage")
	rootCmd.PersistentFlags().BoolVarP(&verbosity, "verbosity", "v", false, "usage")
}

type headers map[string]string

func parseHeaders(h string) headers {
	res := headers{}
	if len(h) == 0 {
		return res
	}

	hs := strings.Split(h, ";")
	for _, v := range hs {
		if len(v) == 0 {
			continue
		}
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}
		res[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return res
}
