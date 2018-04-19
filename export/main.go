package export

func run() {

	files, _ := listDir("./", "xlsx")

	for _, f := range files {

		doExportFile(f, "./")
	}
}
