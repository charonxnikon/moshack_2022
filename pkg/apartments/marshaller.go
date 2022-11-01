package apartments

import (
	"encoding/json"
	apartmentTypes "moshack_2022/pkg/apartments/type"
)

type ApartmentJSON struct {
	Address         string  `json:"Address"`
	Rooms           int     `json:"Rooms"`
	BuildingSegment string  `json:"BuildingSegment"`
	BuildingFloors  int     `json:"BuildingFloors"`
	WallMaterial    string  `json:"WallMaterial"`
	ApartmentFloor  int     `json:"ApartmentFloor"`
	ApartmentArea   float64 `json:"ApartmentArea"`
	KitchenArea     float64 `json:"KitchenArea"`
	Balcony         string  `json:"Balcony"`
	MetroRemoteness int     `json:"MetroRemoteness"`
	Condition       string  `json:"Condition"`
}

func MarshalApartments(apartments []Apartment) []byte {
	type respBody struct {
		Apartments []ApartmentJSON `json:"apartments"`
	}
	var jsonApartments []ApartmentJSON
	for _, el := range apartments {
		jsonApartments = append(jsonApartments, ApartmentJSON{
			Address:         el.Address,
			Rooms:           el.Rooms,
			BuildingSegment: apartmentTypes.BuildingSegment.GetJSON(el.BuildingSegment),
			BuildingFloors:  el.BuildingFloors,
			WallMaterial:    apartmentTypes.WallMaterial.GetJSON(el.WallMaterial),
			ApartmentFloor:  el.ApartmentFloor,
			ApartmentArea:   el.ApartmentArea,
			KitchenArea:     el.KitchenArea,
			Balcony:         apartmentTypes.Balcony.GetJSON(el.Balcony),
			MetroRemoteness: el.MetroRemoteness,
			Condition:       apartmentTypes.Condition.GetJSON(el.Condition),
		})
	}
	data, _ := json.Marshal(respBody{Apartments: jsonApartments})
	return data
}
