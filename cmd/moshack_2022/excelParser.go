package main

import (
	"encoding/json"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type apartmentJSON struct {
	Address string `json:"Address"`
	Rooms   string `json:"Rooms"`
}

type apartmentDB struct {
	Address string
	Rooms   int
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

func insertApartmentToPSQL(json apartmentJSON, db *gorm.DB) {
	rooms, err := strconv.Atoi(json.Rooms)
	if err != nil {
		panic(err)
	}
	newApartment := apartmentDB{
		Address: json.Address,
		Rooms:   rooms,
	}
	db.Table("apartments").Create(newApartment)
}

func ExcelParser(excelSheet *xls.Sheet, db *gorm.DB) []apartmentJSON {
	var resJSON []apartmentJSON

	//TODO:добавить проверку соответсвия первой строки excelColumnNames
	//TODO:Panic() в работе сервера так себе, пожалуй

	for i := 1; i < excelSheet.GetNumberRows(); i++ {
		row, err := excelSheet.GetRow(i)
		if err != nil {
			log.Panic(err)
		}
		cells := row.GetCols()
		if len(cells) != len(excelColumnNames) {
			log.Panic("В екселе больше столбцов чем надо")
		}

		newJSON := apartmentJSON{
			Address: cells[0].GetString(),
			Rooms:   cells[1].GetString(),
		}
		insertApartmentToPSQL(newJSON, db)

		resJSON = append(resJSON, apartmentJSON{
			Address: cells[0].GetString(),
			Rooms:   cells[1].GetString(),
		})
	}

	return resJSON
}

func printJSON(jsonArray []apartmentJSON) {
	for _, el := range jsonArray {
		str, err := json.MarshalIndent(el, "", "  ")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(str))
	}
}
