package apartments

type Apartment struct {
	ID                 uint32
	Location           string
	Coords             [2]float64
	RoomsNumber        int
	Segment            string
	WallMaterial       string
	Floor              int
	Square             int
	KitchenSquare      int
	HasBalcony         bool
	DistanceFromSubway int
	Condition          string
}

type ApartmentRepo interface {
	//	GetAll() ([]*Apartment, error)
	GetByID(id uint32) (*Apartment, error)
	Add(apartment *Apartment) (uint32, error)
	//	Update(newApartment *Apartment) (bool, error)
	Delete(id uint32) (bool, error)
}
