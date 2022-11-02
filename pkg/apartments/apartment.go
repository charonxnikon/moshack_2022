package apartments

type Apartment struct {
	ID              uint32  `json:"id"        gorm:"column:id; primaryKey"`
	UserID          uint32  `json:"-"         gorm:"column:user_id"`
	Address         string  `json:"address"   gorm:"column:address"`
	Rooms           int16   `json:"rooms"     gorm:"column:rooms"`
	BuildingSegment string  `json:"type"      gorm:"column:type"`
	BuildingFloors  int16   `json:"height"    gorm:"column:height"`
	WallMaterial    string  `json:"material"  gorm:"column:material"`
	ApartmentFloor  int16   `json:"floor"     gorm:"column:floor"`
	ApartmentArea   float64 `json:"area"      gorm:"column:area"`
	KitchenArea     float64 `json:"kitchen"   gorm:"column:kitchen"`
	Balcony         string  `json:"balcony"   gorm:"column:balcony"`
	MetroRemoteness int     `json:"metro"     gorm:"column:metro"`
	Condition       string  `json:"condition" gorm:"column:condition"`
	Latitude        float64 `json:"latitude"  gorm:"column:latitude"`
	Longitude       float64 `json:"longitude" gorm:"column:longitude"`
}

type ApartmentRepo interface {
	GetByID(id uint32) (*Apartment, error)
	GetAllByUserID(userID uint32) ([]Apartment, error)
	Add(apartment *Apartment) (uint32, error)
	Delete(id uint32) (bool, error)
}
