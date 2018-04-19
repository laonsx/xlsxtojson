package export

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const jsonTemplate = `[{{range $item :=.}}
	{　{{range $kv :=$item.value}}{{if $kv.v}}
		{{$kv.key}} : {{$kv.value}}{{if not $kv.end}},{{end}}{{else if $kv.o}}
		{{$kv.key}} : {　{{range $kvo :=$kv.value}}
			{{$kvo.key}} : {{$kvo.value}}{{if not $kvo.end}},{{end}}{{end}}
		}{{if not $kv.end}},{{end}}{{else if $kv.a}}
		{{$kv.key}} : [{{range $kva :=$kv.value}}
			{{$kva.value}}{{if not $kva.end}},{{end}}{{end}}
		]{{if not $kv.end}},{{end}}{{else if $kv.ao}}
		{{$kv.key}} : [　{{range $end,$kvao :=$kv.value}}
			{　{{range $ao :=$kvao}}
				{{$ao.key}} : {{$ao.value}}{{if not $ao.end}},{{end}}{{end}}
			}{{if lt $end $kv.count}},{{end}}{{end}}
		]{{if not $kv.end}},{{end}}{{end}}{{end}}
	}{{if $item.end}},{{end}}{{end}}
]`

func doExportFile(filename, dir string) {

	basename := filepath.Base(filename)
	jsonFileName := strings.TrimSuffix(basename, filepath.Ext(basename))

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {

		fmt.Println(filename, "=> Open excel failed.")

		return
	}

	//var data [][]map[string]interface{}
	var datas []map[string]interface{}
	var specsData []map[string]interface{}

	sheet := xlFile.Sheets[0]

	if len(sheet.Rows) <= 2 {

		fmt.Println(filename, "=> empty")

		return
	}

	for idxRow, row := range sheet.Rows {

		if idxRow == 0 {

			specsData = initSpecsRow(row)
			continue
		}

		if idxRow == 1 {

			addSpecsDesc(row, specsData)
			continue
		}

		if strings.Contains(row.Cells[0].String(), "#") {

			continue
		}

		rowData := initRowData(row, specsData)

		data := map[string]interface{}{
			"value": rowData,
			"end":   (idxRow + 1) != len(sheet.Rows),
		}

		datas = append(datas, data)
	}

	os.MkdirAll(path.Join(dir, "json"), os.ModePerm)

	outputFileName := path.Join(dir, "json", jsonFileName+".json")

	f, err := os.Create(outputFileName)
	if err != nil {

		fmt.Println("create file: ", err)

		return
	}
	defer f.Close()

	t := template.Must(template.New("php").Parse(jsonTemplate))
	t.Execute(f, datas)
}

func isDescRow(cellString string) bool {
	return strings.Contains(cellString, "＃") || strings.Contains(cellString, "#")
}

func initSpecsRow(row *xlsx.Row) []map[string]interface{} {

	specsData := make([]map[string]interface{}, 1)

	for idxCell, cell := range row.Cells {

		if idxCell == 0 {

			continue
		}

		cellInfo := strings.Split(cell.String(), "#")
		cellData := make(map[string]interface{})

		if len(cellInfo) != 2 {

			panic("")
		}

		cellData["key"] = fmt.Sprintf(`"%v"`, cellInfo[0])
		cellData["type"] = cellInfo[1]

		specsData = append(specsData, cellData)
	}

	return specsData
}

func addSpecsDesc(row *xlsx.Row, specsData []map[string]interface{}) {

	for idxCell, cell := range row.Cells {

		if idxCell == 0 || len(cell.String()) == 0 || idxCell > len(specsData) {

			continue
		}

		var descDataInfo []map[string]interface{}

		switch specsData[idxCell]["type"] {

		case "{}", "[{}]":

			cellInfo := strings.Split(cell.String(), ";")

			for _, descInfo := range cellInfo {

				desc := strings.Split(descInfo, "#")

				descData := make(map[string]interface{})

				descData["key"] = fmt.Sprintf(`"%v"`, desc[0])
				descData["type"] = desc[1]

				descDataInfo = append(descDataInfo, descData)
			}

		case "[]":

			desc := strings.Split(cell.String(), "#")

			descData := make(map[string]interface{})

			descData["type"] = desc[1]
			descDataInfo = append(descDataInfo, descData)

		default:

		}

		specsData[idxCell]["desc"] = descDataInfo
	}

}

func initRowData(row *xlsx.Row, specsData []map[string]interface{}) (data []map[string]interface{}) {

	for idxCell, cell := range row.Cells {

		cellData := make(map[string]interface{})

		if idxCell == 0 {

			continue
		}

		cellSpecs := specsData[idxCell]
		switch cellSpecs["type"] {

		case "string":

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = fmt.Sprintf(`"%v"`, cell.String())
			cellData["v"] = true

		case "int":

			cellData["key"] = cellSpecs["key"]
			cellInt, _ := cell.Int64()
			cellData["value"] = cellInt
			cellData["v"] = true

		case "bool":

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = cell.String()
			cellData["v"] = true

		case "date":

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = fmt.Sprintf(`"%v"`, cell.String())
			cellData["v"] = true

		case "{}":

			var itemDatas []map[string]interface{}

			cellInfo := strings.Split(cell.String(), ";")
			cellItemSpecs := cellSpecs["desc"].([]map[string]interface{})

			for index, value := range cellInfo {

				item := make(map[string]interface{})
				item["key"] = cellItemSpecs[index]["key"]
				item["value"] = getItemValue(cellItemSpecs[index]["type"].(string), value)

				if index+1 == len(cellInfo) {

					item["end"] = true
				}

				itemDatas = append(itemDatas, item)
			}

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = itemDatas
			cellData["o"] = true

		case "[]":

			var itemDatas []map[string]interface{}

			cellInfo := strings.Split(cell.String(), "|")
			cellItemSpecs := cellSpecs["desc"].([]map[string]interface{})

			for index, value := range cellInfo {

				if len(value) == 0 {

					continue
				}
				item := make(map[string]interface{})
				//item["key"] = cellItemSpecs[0]["key"]
				item["value"] = getItemValue(cellItemSpecs[0]["type"].(string), value)

				if index+1 == len(cellInfo) {

					item["end"] = true
				}

				itemDatas = append(itemDatas, item)
			}

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = itemDatas
			cellData["a"] = true

		case "[{}]":

			var itemDatas [][]map[string]interface{}

			cellInfo := strings.Split(cell.String(), "|")
			cellItemSpecs := cellSpecs["desc"].([]map[string]interface{})

			for _, value := range cellInfo {

				var items []map[string]interface{}

				itemInfo := strings.Split(value, ";")

				for itemIndex, itemValue := range itemInfo {

					item := make(map[string]interface{})

					item["key"] = cellItemSpecs[itemIndex]["key"]
					item["value"] = getItemValue(cellItemSpecs[itemIndex]["type"].(string), itemValue)

					if itemIndex+1 == len(itemInfo) {

						item["end"] = true
					}

					items = append(items, item)
				}

				itemDatas = append(itemDatas, items)
			}

			cellData["key"] = cellSpecs["key"]
			cellData["value"] = itemDatas
			cellData["ao"] = true
			cellData["count"] = len(cellInfo) - 1

		default:

		}

		if idxCell+1 == len(row.Cells) {

			cellData["end"] = true
		}

		data = append(data, cellData)
	}

	return
}

func getItemValue(itype, cellString string) interface{} {

	switch itype {

	case "string", "date":

		return fmt.Sprintf(`"%s"`, cellString)

	case "int", "bool":

		return cellString

	default:

	}

	return nil
}
