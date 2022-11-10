package apartments

type UserApartment struct {
	ID         uint32 `gorm:"column:id; primaryKey"`
	UserID     uint32 `json:"-"`
	Address    string
	Rooms      string
	Type       string
	Height     int16
	Material   string
	Floor      int16
	Area       float64
	Kitchen    float64
	Balcony    string
	Metro      int
	Condition  string
	Latitude   float64
	Longitude  float64
	TotalPrice float64
	PriceM2    float64
}

type DBApartment struct {
	ID         uint32 `gorm:"column:id; primaryKey"`
	UserID     uint32 `json:"-"`
	Address    string
	Rooms      string
	Type       string
	Height     int16
	Material   string
	Floor      int16
	Area       float64
	Kitchen    float64
	Balcony    string
	Metro      int
	Condition  string
	Latitude   float64
	Longitude  float64
	TotalPrice float64
	PriceM2    float64
}

type ApartmentRepo interface {
	GetUserApartmentByID(id uint32) (*UserApartment, error)
	GetDBApartmentByID(id uint32) (*DBApartment, error)
	GetAllUserApartmentsByUserID(userID uint32) ([]UserApartment, error)
	AddUserApartment(apartment *UserApartment) (uint32, error)
	DeleteUserApartment(id uint32) (bool, error)
}

const (
	userApartmentsTable = "user_apartments"
	dbApartmentsTable   = "db_apartments"
)
