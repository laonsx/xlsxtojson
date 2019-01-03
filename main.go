package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	var fileDir, outDir string

	flag.StringVar(&fileDir, "f", "./", "File need to export.")
	flag.StringVar(&outDir, "o", "./json", "Output destination.")

	flag.Parse()

	fmt.Println("欢迎使用xlsxtojson工具\n")

	files, _ := listDir(fileDir, "xlsx")
	for _, f := range files {

		doExportFile(f, outDir)
	}

	fmt.Println("\n转换完成，5秒后自动关闭。")
	time.Sleep(5 * time.Second)
}
