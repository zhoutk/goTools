package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/Luxurioust/excelize"
	"encoding/json"
)

type answer struct {
	right bool
	answer string
}

func main() {
	xlsx, err := excelize.OpenFile("./qs.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get value from cell by given sheet index and axis.
	//cell := xlsx.GetCellValue("Sheet1", "B2")
	//fmt.Println(cell)
	// Get sheet index.
	index := xlsx.GetSheetIndex("Sheet1")
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("Sheet" + strconv.Itoa(index))
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var answers []answer
		answers = append(answers, answer{ true,row[3] })
		answers = append(answers, answer{ false,row[4] })
		answers = append(answers, answer{ false,row[5] })
		answers = append(answers, answer{ false,row[6] })
		data, err := json.Marshal(answers)
		if err == nil {
			fmt.Printf("%s\n", data)
		}
	}
}