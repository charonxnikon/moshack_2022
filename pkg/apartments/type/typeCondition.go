package excelParser

var typeCondition = enumType{
	numberOfStates: 3,
	possibleInput: [][]string{
		{"Без отделки"},
		{"Муниципальный ремонт"},
		{"Современная отделка"},
	},
	jsonOutput: []string{
		"WithoutRenovation",
		"MunicipalRenovation",
		"ModernRenovation",
	},
	jsonError: "Error",
}
