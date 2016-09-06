package main

import (
	"fmt"
	"log"

	"github.com/sunreaver/gotools/system"
	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := system.CurPath() + system.SystemSep() + "Book1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Println("open err:", err)
		panic(0)
	}
	over := map[string]int{}
	repeat := 0
	for _, sheet := range xlFile.Sheets {
		if sheet.Name != "Sheet1" {
			continue
		}
		for rnum, row := range sheet.Rows {
			for index, cell := range row.Cells {
				if index != 2 || len(cell.Value) == 0 {
					continue
				}
				if _, ok := over[cell.Value]; ok {
					repeat++
					fmt.Println("Repeat:", cell.Value, rnum, "=", over[cell.Value])
				}
				over[cell.Value] = rnum
			}
		}
	}

	fmt.Println("Row:", len(over), "repeat num:", repeat)
}
