package excelParser

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"moshack_2022/pkg/apartments"
	"strconv"

	"github.com/shakinm/xlsReader/xls"
)

// TODO: по-хорошему, надо быпереписать так, чтобы можно было легко расширять набор свойств из подпакетов
//Ожидаемый вид таблицы: ...| (номер столбца) Информация |
//(0) Адрес |(1) Кол-во комнат |(2) Тип здания |(3) Кол-во этажей |(4) Материал стен |(5) Этаж квартиры |(6) Площадь квартиры |
//|(7) Площадь кухни |(8) Тип балкона |(9) Удаленность от метро | (10) Состояние |  |  |  |  |

type ExcelParser struct {
	FileName   string
	excelSheet *xls.Sheet
	Apartments []apartments.Apartment
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
	"Удаленность от метро",
	"Состояние",
}

func invalidValueError(row int, cell int, err error) error {
	return fmt.Errorf("invalid value if row %d, cell %s, error: %s",
		row, excelColumnNames[cell], err)
}

func ParseXLS(file multipart.File, userID uint32) ([]*apartments.Apartment, error) {
	xlsFile, err := xls.OpenReader(file)
	if err != nil {
		return nil, err
	}
	xlsSheet, err := xlsFile.GetSheet(0)
	if err != nil {
		return nil, err
	}

	result := make([]*apartments.Apartment, 0)
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

		rooms, err := strconv.Atoi(cells[1].GetString())
		if err != nil {
			return nil, invalidValueError(rowNum, 1, err)
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

		newApartment := &apartments.Apartment{
			UserID:  userID,
			Address: cells[0].GetString(),
			Rooms:   int16(rooms),
			//Type: apartmentTypes.Type.GetState(cells[2].GetString()),
			Type:   cells[2].GetString(),
			Height: int16(floors),
			//Material:    apartmentTypes.Material.GetState(cells[4].GetString()),
			Material: cells[4].GetString(),
			Floor:    int16(floor),
			Area:     aSquare,
			Kitchen:  kSquare,
			//Balcony:         apartmentTypes.Balcony.GetState(cells[8].GetString()),
			Balcony: cells[8].GetString(),
			Metro:   metroRemotneness,
			//Condition:       apartmentTypes.Condition.GetState(cells[10].GetString()),
			Condition: cells[10].GetString(),
		}

		result = append(result, newApartment)
	}

	return result, nil
}

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

/*

func insertApartmentToPSQL(apartment *apartments.Apartment, db *gorm.DB) {
	db.Table("apartments").Create(*apartment)
}

func (excel ExcelParser) Parse(db *gorm.DB) *ExcelParser {

	excel.excelSheet = OpenExcel(excel.FileName)

	//TODO:добавить проверку соответсвия первой строки excelColumnNames
	//TODO:Panic() в работе сервера так себе, пожалуй

	for rowNum := 1; rowNum < excel.excelSheet.GetNumberRows(); rowNum++ {
		row, err := excel.excelSheet.GetRow(rowNum)
		if err != nil {
			log.Panic(err)
		}
		cells := row.GetCols()
		//TODO: возможно, стоит формировтаь массив всех возникших ошибок и возвращать на фронт, чтобы юзер мог понять что не так
		rooms, err := strconv.Atoi(cells[1].GetString())
		floors, err := strconv.Atoi(cells[3].GetString())
		floor, err := strconv.Atoi(cells[5].GetString())
		aSquare, err := strconv.ParseFloat(cells[6].GetString(), 64)
		kSquare, err := strconv.ParseFloat(cells[7].GetString(), 64)
		metroRemotneness, err := strconv.Atoi(cells[9].GetString())
		if err != nil {
			log.Panic(err)
		}
		if len(cells) != len(excelColumnNames) {
			log.Panic("В екселе больше столбцов чем надо")
		}

		newApartment := apartments.Apartment{
			Address:         cells[0].GetString(),
			Rooms:           rooms,
			Type: apartmentTypes.Type.GetState(cells[2].GetString()),
			Height:  floors,
			Material:    apartmentTypes.Material.GetState(cells[4].GetString()),
			Floor:  floor,
			Area:   aSquare,
			Kitchen:     kSquare,
			Balcony:         apartmentTypes.Balcony.GetState(cells[8].GetString()),
			Metro: metroRemotneness,
			Condition:       apartmentTypes.Condition.GetState(cells[10].GetString()),
		}
		insertApartmentToPSQL(&newApartment, db)

		//TODO:если мы не хотим хранить весь отпаршеный ексель в памяти, то эта страчка не нужна
		(excel).Apartments = append((excel).Apartments, newApartment)
	}

	return &excel
}
*/
