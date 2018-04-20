package main

import "flag"

func main() {

	var fileDir, outDir string

	flag.StringVar(&fileDir, "f", "./", "File need to export.")
	flag.StringVar(&outDir, "o", "./json", "Output destination.")

	flag.Parse()

	files, _ := listDir(fileDir, "xlsx")
	for _, f := range files {

		doExportFile(f, outDir)
	}
}
