package main

import (
	"encoding/json"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type apartment struct {
	Address string `json:"Address"`
	Rooms   int    `json:"Rooms"`
}

type ExcelParser struct {
	fileName   string
	excelSheet *xls.Sheet
	Apartments []apartment
}

var excelColumnNames = []string{"Адрес", "Комнаты"}

func OpenExcel(name string) *xls.Sheet {
	excelFile, err := xls.OpenFile(name)
	if err != nil {
		fmt.Println("@@excelFileErr@@")
		log.Panic(err.Error())
	}
	excelSheet, err := excelFile.GetSheet(0)
	if err != nil {
		fmt.Println("@@exclerSheetErr@@")
		log.Panic(err.Error())
	}
	return excelSheet
}

func insertApartmentToPSQL(json *apartment, db *gorm.DB) {
	db.Table("apartments").Create(*json)
}

func (excel *ExcelParser) parse(db *gorm.DB) *ExcelParser {

	excel.excelSheet = OpenExcel(excel.fileName)

	//TODO:добавить проверку соответсвия первой строки excelColumnNames
	//TODO:Panic() в работе сервера так себе, пожалуй

	for i := 1; i < excel.excelSheet.GetNumberRows(); i++ {
		row, err := excel.excelSheet.GetRow(i)
		if err != nil {
			log.Panic(err)
		}
		cells := row.GetCols()
		rooms, err := strconv.Atoi(cells[1].GetString())
		if err != nil {
			log.Panic(err)
		}
		if len(cells) != len(excelColumnNames) {
			log.Panic("В екселе больше столбцов чем надо")
		}

		newJSON := apartment{
			Address: cells[0].GetString(),
			Rooms:   rooms,
		}
		insertApartmentToPSQL(&newJSON, db)

		(*excel).Apartments = append((*excel).Apartments, newJSON)
	}

	return excel
}

func (excel *ExcelParser) marshalExcel() []byte {
	type respBody struct {
		Apartments []apartment `json:"apartments"`
	}
	data, _ := json.Marshal(respBody{Apartments: excel.Apartments})
	return data
}
