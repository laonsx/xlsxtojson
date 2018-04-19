package export

func Run() {

	files, _ := listDir("./", "xlsx")

	for _, f := range files {

		doExportFile(f, "./")
	}
}
