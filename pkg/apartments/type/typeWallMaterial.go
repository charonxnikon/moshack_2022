package excelParser

var typeWallMaterial = enumType{
	numberOfStates: 3,
	possibleInput: [][]string{
		{"Кирпич"},
		{"Панель"},
		{"Монолит"},
	},
	jsonOutput: []string{
		"brickWall",
		"panelWall",
		"monolithWall",
	},
	jsonError: "Error",
}
