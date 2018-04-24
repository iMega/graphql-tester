package tester

import (
	"os"
	"path/filepath"
)

func FileList(args []string) ([]string, error) {
	var list []string

	for _, arg := range args {
		fi, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					list = append(list, path)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		case mode.IsRegular():
			list = append(list, arg)
		}
	}

	return list, nil
}
