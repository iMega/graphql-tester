package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const version string = "0.0.1"

type OptionsCmd struct {
	URL     string
	Path    []string
	Headers map[string]string
}

var (
	Options OptionsCmd

	url     string
	headers string

	rootCmd = &cobra.Command{
		Use:     "qltester",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			Options = OptionsCmd{
				URL:     strings.Trim(url, "' \""),
				Path:    args,
				Headers: map[string]string{},
			}

			hs := strings.Split(headers, ";")
			for _, v := range hs {
				kv := strings.Split(v, ":")
				Options.Headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
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
	rootCmd.PersistentFlags().StringVarP(&headers, "headers", "H", "", "usage")
}
