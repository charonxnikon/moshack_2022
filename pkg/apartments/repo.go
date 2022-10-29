package apartments

import (
	"errors"

	"gorm.io/gorm"
)

var (
	errNoApartment = errors.New("No apartment found")
)

type ApartmentDB struct {
	db *gorm.DB
}

func NewApartmentRepo(db *gorm.DB) *ApartmentDB {
	return &ApartmentDB{
		db: db,
	}
}

// func GetAll() ([]*Apartment, error)

func (adb *ApartmentDB) GetByID(id uint32) (*Apartment, error) {
	apartments := make([]Apartment, 0)
	adb.db.Where("id = ?", id).Find(&apartments)
	if len(apartments) != 1 {
		return nil, errNoApartment
	}

	return &apartments[0], nil
}

func (adb *ApartmentDB) Add(apartment *Apartment) (uint32, error) {
	result := adb.db.Create(apartment)
	if result.Error != nil {
		return 0, result.Error
	}

	return apartment.ID, nil
}

// Update(newApartment *Apartment) (bool, error)

func (adb *ApartmentDB) Delete(id uint32) (bool, error) {
	apartment := Apartment{}
	result := adb.db.Where("id = ?", id).Delete(&apartment)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil

}
