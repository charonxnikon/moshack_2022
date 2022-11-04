package apartments

// ID              uint32  `json:"id"        gorm:"column:id; primaryKey"`
// UserID          uint32  `json:"-"         gorm:"column:user_id"`
// Address         string  `json:"address"   gorm:"column:address"`
// Rooms           int16   `json:"rooms"     gorm:"column:rooms"`
// BuildingSegment string  `json:"type"      gorm:"column:type"`
// BuildingFloors  int16   `json:"height"    gorm:"column:height"`
// WallMaterial    string  `json:"material"  gorm:"column:material"`
// ApartmentFloor  int16   `json:"floor"     gorm:"column:floor"`
// ApartmentArea   float64 `json:"area"      gorm:"column:area"`
// KitchenArea     float64 `json:"kitchen"   gorm:"column:kitchen"`
// Balcony         string  `json:"balcony"   gorm:"column:balcony"`
// MetroRemoteness int     `json:"metro"     gorm:"column:metro"`
// Condition       string  `json:"condition" gorm:"column:condition"`
// Latitude        float64 `json:"latitude"  gorm:"column:latitude"`
// Longitude       float64 `json:"longitude" gorm:"column:longitude"`

type Apartment struct {
	ID        uint32 `gorm:"column:id; primaryKey"`
	UserID    uint32 `json:"-" gorm:"column:user_id"`
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
	GetUserApartmentByID(id uint32) (*Apartment, error)
	GetDBApartmentByID(id uint32) (*Apartment, error)
	GetAllUserApartmentsByUserID(userID uint32) ([]Apartment, error)
	AddUserApartment(apartment *Apartment) (uint32, error)
	DeleteUserApartment(id uint32) (bool, error)
}

const (
	userApartmentsTable = "user_apartments"
	dbApartmentsTable   = "db_apartments"
)
