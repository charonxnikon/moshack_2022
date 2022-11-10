package excelParser

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"moshack_2022/pkg/apartments"
	"os"
	"strconv"

	"github.com/shakinm/xlsReader/xls"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExcelParser struct {
	FileName   string
	excelSheet *xls.Sheet
	Apartments []apartments.UserApartment
}

var excelColumnNames = []string{
	"Адрес",
	"Кол-во комнат",
	"Тип здания",
	"Кол-во этажей",
	"Материал стен",
	"Этаж квартиры",
	"Площадь квартиры",
	"Площадь кухни",
	"Тип балкона",
	"Удалённость от метро",
	"Состояние",
	"Стоимость",
	"Стоимость кв.м.",
}

func invalidValueError(row int, cell int, err error) error {
	return fmt.Errorf("invalid value if row %d, cell %s, error: %s",
		row, excelColumnNames[cell], err)
}

func ParseXLSX(file multipart.File, userID uint32) ([]*apartments.UserApartment, error) {
	xlsFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := xlsFile.Close(); err != nil {
			log.Println("Error during closing excel file:" + err.Error())
		}
	}()

	result := make([]*apartments.UserApartment, 0)
	rows, err := xlsFile.GetRows(xlsFile.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	for rowNum, row := range rows {
		if err != nil {
			return nil, err
		}
		if rowNum == 0 {
			continue
		}
		if len(row) != len(excelColumnNames)-2 { // -2 because of price columns
			errStr := fmt.Sprintf("invalid number of conumns in xls file: expected %d, got %d",
				len(excelColumnNames),
				len(row))
			return nil, errors.New(errStr)
		}

		floors, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, invalidValueError(rowNum, 3, err)
		}

		floor, err := strconv.Atoi(row[5])
		if err != nil {
			return nil, invalidValueError(rowNum, 5, err)
		}

		aSquare, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			return nil, invalidValueError(rowNum, 6, err)
		}

		kSquare, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			return nil, invalidValueError(rowNum, 7, err)
		}

		metroRemotneness, err := strconv.Atoi(row[9])
		if err != nil {
			return nil, invalidValueError(rowNum, 9, err)
		}

		newApartment := &apartments.UserApartment{
			UserID:    userID,
			Address:   row[0],
			Rooms:     row[1],
			Type:      row[2],
			Height:    int16(floors),
			Material:  row[4],
			Floor:     int16(floor),
			Area:      aSquare,
			Kitchen:   kSquare,
			Balcony:   row[8],
			Metro:     metroRemotneness,
			Condition: row[10],
		}

		result = append(result, newApartment)
	}

	return result, nil
}

func ParseXLS(file multipart.File, userID uint32) ([]*apartments.UserApartment, error) {
	xlsFile, err := xls.OpenReader(file)
	if err != nil {
		return nil, err
	}
	xlsSheet, err := xlsFile.GetSheet(0)
	if err != nil {
		return nil, err
	}

	result := make([]*apartments.UserApartment, 0)
	for rowNum := 1; rowNum < xlsSheet.GetNumberRows(); rowNum++ {
		row, err := xlsSheet.GetRow(rowNum)
		if err != nil {
			return nil, err
		}

		cells := row.GetCols()
		if len(cells) != len(excelColumnNames) {
			errStr := fmt.Sprintf("invalid number of conumns in xls file: expected %d, got %d",
				len(excelColumnNames),
				len(cells))
			return nil, errors.New(errStr)
		}
		floors, err := strconv.Atoi(cells[3].GetString())
		if err != nil {
			return nil, invalidValueError(rowNum, 3, err)
		}

		floor, err := strconv.Atoi(cells[5].GetString())
		if err != nil {
			return nil, invalidValueError(rowNum, 5, err)
		}

		aSquare, err := strconv.ParseFloat(cells[6].GetString(), 64)
		if err != nil {
			return nil, invalidValueError(rowNum, 6, err)
		}

		kSquare, err := strconv.ParseFloat(cells[7].GetString(), 64)
		if err != nil {
			return nil, invalidValueError(rowNum, 7, err)
		}

		metroRemotneness, err := strconv.Atoi(cells[9].GetString())
		if err != nil {
			return nil, invalidValueError(rowNum, 9, err)
		}

		newApartment := &apartments.UserApartment{
			UserID:    userID,
			Address:   cells[0].GetString(),
			Rooms:     cells[1].GetString(),
			Type:      cells[2].GetString(),
			Height:    int16(floors),
			Material:  cells[4].GetString(),
			Floor:     int16(floor),
			Area:      aSquare,
			Kitchen:   kSquare,
			Balcony:   cells[8].GetString(),
			Metro:     metroRemotneness,
			Condition: cells[10].GetString(),
		}

		result = append(result, newApartment)
	}

	return result, nil
}

func UnparseXLSX(aparts []apartments.UserApartment) (*xlsx.File, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Результат рассчёта")
	if err != nil {
		return nil, err
	}

	newRow := sheet.AddRow()
	for _, cellName := range excelColumnNames {
		newCell := newRow.AddCell()
		newCell.SetString(cellName)
	}

	for _, apartment := range aparts {
		newRow := sheet.AddRow()
		newRow.AddCell().SetString(apartment.Address)
		newRow.AddCell().SetString(apartment.Rooms)
		newRow.AddCell().SetString(apartment.Type)
		newRow.AddCell().SetInt(int(apartment.Height))
		newRow.AddCell().SetString(apartment.Material)
		newRow.AddCell().SetInt(int(apartment.Floor))
		newRow.AddCell().SetFloat(apartment.Area)
		newRow.AddCell().SetFloat(apartment.Kitchen)
		newRow.AddCell().SetString(apartment.Balcony)
		newRow.AddCell().SetInt(apartment.Metro)
		newRow.AddCell().SetString(apartment.Condition)
		newRow.AddCell().SetFloat(apartment.TotalPrice)
		newRow.AddCell().SetFloat(apartment.PriceM2)
	}

	return file, nil
}

func MakeXLSX_OLD(db *gorm.DB, localSave bool, filePath string) (interface{}, error) {
	const sheetName = "Квартиры"
	if localSave {
		if _, err := os.Stat(filePath); err == nil {
			//указанный файл уже существует
			err := os.Rename(filePath, filePath+"_older")
			if err != nil {
				return nil, err
			}
		}
	}
	var aparts []apartments.UserApartment
	db.Find(&aparts)

	f := xlsx.NewFile()
	sheet, err := f.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}
	newRow := sheet.AddRow()
	for _, cellName := range excelColumnNames {
		newCell := newRow.AddCell()
		newCell.SetString(cellName)
	}

	for _, apartment := range aparts {
		newRow := sheet.AddRow()
		newRow.AddCell().SetString(apartment.Address)
		newRow.AddCell().SetString(apartment.Rooms)
		newRow.AddCell().SetString(apartment.Type)
		newRow.AddCell().SetInt(int(apartment.Height))
		newRow.AddCell().SetString(apartment.Material)
		newRow.AddCell().SetInt(int(apartment.Floor))
		newRow.AddCell().SetFloat(apartment.Area)
		newRow.AddCell().SetFloat(apartment.Kitchen)
		newRow.AddCell().SetString(apartment.Balcony)
		newRow.AddCell().SetInt(apartment.Metro)
		newRow.AddCell().SetString(apartment.Condition)
	}

	if localSave {
		err = f.Save(filePath)
	}

	if false {
		err := os.Remove(filePath)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}
	return nil, nil
}
