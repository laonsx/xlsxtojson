package export

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

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)

	for _, f := range dir {

		if f.IsDir() {

			continue
		}

		if strings.HasSuffix(strings.ToUpper(f.Name()), suffix) {

			files = append(files, dirPth+PthSep+f.Name())
		}
	}

	return files, nil
}
