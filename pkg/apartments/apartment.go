package apartments

type Apartment struct {
	ID        uint32
	UserID    uint32
	Address   string
	Rooms     int16
	Type      string
	Height    int16
	Material  string
	Floor     int16
	Area      float64
	Kitchen   float64
	Balcony   string
	Metro     int
	Condition string
	Latitude  float64
	Longitude float64
}

type ApartmentRepo interface {
	//	GetAll() ([]*Apartment, error)
	GetByID(id uint32) (*Apartment, error)
	Add(apartment *Apartment) (uint32, error)
	//	Update(newApartment *Apartment) (bool, error)
	Delete(id uint32) (bool, error)
}
