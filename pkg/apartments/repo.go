package apartments

import (
	"errors"

	"gorm.io/gorm"
)

var (
	errNoApartment = errors.New("no apartment found")
)

type ApartmentDB struct {
	db *gorm.DB
}

func NewApartmentRepo(db *gorm.DB) *ApartmentDB {
	return &ApartmentDB{
		db: db,
	}
}

func (adb *ApartmentDB) GetByID(id uint32) (*Apartment, error) {
	aparts := make([]Apartment, 0)
	db := adb.db.Where("id = ?", id).Find(&aparts)
	if db.Error != nil {
		return nil, db.Error
	}
	if len(aparts) != 1 {
		return nil, errNoApartment
	}

	return &aparts[0], nil
}

func (adb *ApartmentDB) GetAllByUserID(userID uint32) ([]Apartment, error) {
	aparts := make([]Apartment, 0)
	db := adb.db.Where("user_id = ?", userID).Find(&aparts)
	if db.Error != nil {
		return nil, db.Error
	}

	return aparts, nil
}

func (adb *ApartmentDB) Add(apartment *Apartment) (uint32, error) {
	result := adb.db.Create(apartment)
	if result.Error != nil {
		return 0, result.Error
	}

	return apartment.ID, nil
}

func (adb *ApartmentDB) Delete(id uint32) (bool, error) {
	apartment := Apartment{}
	result := adb.db.Where("id = ?", id).Delete(&apartment)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil

}
