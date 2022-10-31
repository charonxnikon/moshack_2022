package handlers

import (
	"fmt"
	"html/template"
	"moshack_2022/pkg/apartments"
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
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 * 1024 * 1025)
	file, header, err := r.FormFile("xls_file")
	if err != nil {
		http.Error(w, `File errror`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "header.Filename %v\n", header.Filename)
	fmt.Fprintf(w, "header.Header %#v\n", header.Header)

	// вызов парсера
}

func (h *ApartmentHandler) Table(w http.ResponseWriter, r *http.Request) {
	// здесь в w пишем json'ы
}
