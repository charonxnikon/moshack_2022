package handlers

import (
	"fmt"
	"html/template"
	"moshack_2022/pkg/apartments"
	"moshack_2022/pkg/apartments/excelParser"
	"net/http"

	"go.uber.org/zap"
)

type ApartmentHandler struct {
	Tmpl          *template.Template
	ApartmentRepo apartments.ApartmentRepo
	Logger        *zap.SugaredLogger
}

func (h *ApartmentHandler) Load(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "loadxls.html", nil)
	if err != nil {
		http.Error(w, "Template errror", http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024 * 1024)
	if err != nil {
		http.Error(w, "File errror: file is too much", http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("xls_file")
	if err != nil {
		http.Error(w, "File errror", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "header.Filename %v\n", header.Filename) //
	fmt.Fprintf(w, "header.Header %#v\n", header.Header)    //

	fmt.Println("qq")
	aparts, err := excelParser.ParseXLS(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, apart := range aparts {
		_, err := h.ApartmentRepo.Add(apart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// что записать в w в ответ?
}

func (h *ApartmentHandler) Table(w http.ResponseWriter, r *http.Request) {
	//TODO:
	//написать запрос, возвращающтй слайс апартаментов и вызвать к нему
	//apartments.MarshalApartments()
	// здесь в w пишем json'ы
}
