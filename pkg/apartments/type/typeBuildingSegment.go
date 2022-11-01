package excelParser

var typeBuildingSegment = enumType{
	numberOfStates: 3,
	possibleInput: [][]string{
		{"Новое"},
		{"Современное"},
		{"Старое"},
	},
	jsonOutput: []string{
		"newBuilding",
		"modernBuilding",
		"oldBuilding",
	},
	jsonError: "Error",
}
