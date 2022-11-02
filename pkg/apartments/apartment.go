package apartments

type Apartment struct {
	ID              uint32
	UserID          uint32
	Address         string
	Rooms           int16
	BuildingSegment string
	BuildingFloors  int16
	WallMaterial    string
	ApartmentFloor  int16
	ApartmentArea   float64
	KitchenArea     float64
	Balcony         string
	MetroRemoteness int
	Condition       string
	Latitude        float64
	Longitude       float64
}

type ApartmentRepo interface {
	//	GetAll() ([]*Apartment, error)
	GetByID(id uint32) (*Apartment, error)
	Add(apartment *Apartment) (uint32, error)
	//	Update(newApartment *Apartment) (bool, error)
	Delete(id uint32) (bool, error)
}
