package apartments

type Apartment struct {
	ID              uint32
	Address         string
	Rooms           int
	BuildingSegment int
	BuildingFloors  int
	WallMaterial    int
	ApartmentFloor  int
	ApartmentArea   float64
	KitchenArea     float64
	Balcony         int
	MetroRemoteness int
	Condition       int
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
