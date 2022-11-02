package apartments

import (
	"encoding/json"
)

type ApartmentJSON struct {
	Address         string  `json:"Address"`
	Rooms           int16   `json:"Rooms"`
	BuildingSegment string  `json:"BuildingSegment"`
	BuildingFloors  int16   `json:"BuildingFloors"`
	WallMaterial    string  `json:"WallMaterial"`
	ApartmentFloor  int16   `json:"ApartmentFloor"`
	ApartmentArea   float64 `json:"ApartmentArea"`
	KitchenArea     float64 `json:"KitchenArea"`
	Balcony         string  `json:"Balcony"`
	MetroRemoteness int     `json:"MetroRemoteness"`
	Condition       string  `json:"Condition"`
}

func MarshalApartments(apartments []*Apartment) []byte {
	type respBody struct {
		Apartments []ApartmentJSON `json:"apartments"`
	}
	var jsonApartments []ApartmentJSON
	for _, el := range apartments {
		jsonApartments = append(jsonApartments, ApartmentJSON{
			Address: el.Address,
			Rooms:   el.Rooms,
			//BuildingSegment: apartmentTypes.BuildingSegment.GetJSON(el.BuildingSegment),
			BuildingSegment: el.BuildingSegment,
			BuildingFloors:  el.BuildingFloors,
			//WallMaterial:    apartmentTypes.WallMaterial.GetJSON(el.WallMaterial),
			WallMaterial:   el.WallMaterial,
			ApartmentFloor: el.ApartmentFloor,
			ApartmentArea:  el.ApartmentArea,
			KitchenArea:    el.KitchenArea,
			//Balcony:         apartmentTypes.Balcony.GetJSON(el.Balcony),
			Balcony:         el.Balcony,
			MetroRemoteness: el.MetroRemoteness,
			//Condition:       apartmentTypes.Condition.GetJSON(el.Condition),
			Condition: el.Condition,
		})
	}
	data, _ := json.Marshal(respBody{Apartments: jsonApartments})
	return data
}
