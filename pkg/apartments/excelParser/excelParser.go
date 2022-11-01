package excelParser

import (
	"encoding/json"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"gorm.io/gorm"
	"log"
	"moshack_2022/pkg/apartments"
	"strconv"
)

// TODO: по-хорошему, надо быпереписать так, чтобы можно было легко расширять набор свойств из подпакетов
//Ожидаемый вид таблицы: ...| (номер столбца) Информация |
//(0) Адрес |(1) Кол-во комнат |(2) Тип здания |(3) Кол-во этажей |(4) Материал стен |(5) Этаж квартиры |(6) Площадь квартиры |
//|(7) Площадь кухни |(8) Тип балкона |(9) Удаленность от метро | (10) Состояние |  |  |  |  |

type ApartmentJSON struct {
	Address         string  `json:"Address"`
	Rooms           int     `json:"Rooms"`
	BuildingSegment string  `json:"BuildingSegment"`
	BuildingFloors  int     `json:"BuildingFloors"`
	WallMaterial    string  `json:"WallMaterial"`
	ApartmentFloor  int     `json:"ApartmentFloor"`
	ApartmentArea   float64 `json:"ApartmentArea"`
	KitchenArea     float64 `json:"KitchenArea"`
	Balcony         string  `json:"Balcony"`
	MetroRemoteness int     `json:"MetroRemoteness"`
	Condition       string  `json:"Condition"`
}

type ExcelParser struct {
	FileName   string
	excelSheet *xls.Sheet
	Apartments []apartments.Apartment
}

var excelColumnNames = []string{"Адрес", "Кол-во комнат", "Тип здания", "Кол-во этажей", "Материал стен",
	"Этаж квартиры", "Площадь квартиры", "Площадь кухни", "Тип балкона", "Удаленность от метро", "Состояние"}

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

func insertApartmentToPSQL(apartment *apartments.Apartment, db *gorm.DB) {
	db.Table("apartments").Create(*apartment)
}

func (excel *ExcelParser) Parse(db *gorm.DB) *ExcelParser {

	excel.excelSheet = OpenExcel(excel.FileName)

	//TODO:добавить проверку соответсвия первой строки excelColumnNames
	//TODO:Panic() в работе сервера так себе, пожалуй

	for i := 1; i < excel.excelSheet.GetNumberRows(); i++ {
		row, err := excel.excelSheet.GetRow(i)
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
			BuildingSegment: typeBuildingSegment.GetState(cells[2].GetString()),
			BuildingFloors:  floors,
			WallMaterial:    typeWallMaterial.GetState(cells[4].GetString()),
			ApartmentFloor:  floor,
			ApartmentArea:   aSquare,
			KitchenArea:     kSquare,
			Balcony:         typeBalcony.GetState(cells[8].GetString()),
			MetroRemoteness: metroRemotneness,
			Condition:       typeCondition.GetState(cells[10].GetString()),
		}
		insertApartmentToPSQL(&newApartment, db)

		(*excel).Apartments = append((*excel).Apartments, newApartment)
	}

	return excel
}

func MarshalExcel(apartments []apartments.Apartment) []byte {
	type respBody struct {
		Apartments []ApartmentJSON `json:"apartments"`
	}
	var jsonApartments []ApartmentJSON
	for _, el := range apartments {
		jsonApartments = append(jsonApartments, ApartmentJSON{
			Address:         el.Address,
			Rooms:           el.Rooms,
			BuildingSegment: typeBuildingSegment.GetJSON(el.BuildingSegment),
			BuildingFloors:  el.BuildingFloors,
			WallMaterial:    typeWallMaterial.GetJSON(el.WallMaterial),
			ApartmentFloor:  el.ApartmentFloor,
			ApartmentArea:   el.ApartmentArea,
			KitchenArea:     el.KitchenArea,
			Balcony:         typeBalcony.GetJSON(el.Balcony),
			MetroRemoteness: el.MetroRemoteness,
			Condition:       typeCondition.GetJSON(el.Condition),
		})
	}
	data, _ := json.Marshal(respBody{Apartments: jsonApartments})
	return data
}
