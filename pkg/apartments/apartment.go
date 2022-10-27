package apartments

type Apartment struct {
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
