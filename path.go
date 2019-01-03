package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func listDir(dirPth string, suffix string) (files []string, err error) {

	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {

		return nil, err
	}

	pthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)

	for _, f := range dir {

		if f.IsDir() {

			continue
		}

		if strings.HasPrefix(f.Name(), "~$") {

			continue
		}

		if strings.HasSuffix(strings.ToUpper(f.Name()), suffix) {

			files = append(files, dirPth+pthSep+f.Name())
		}
	}

	return files, nil
}
