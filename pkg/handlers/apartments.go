package handlers

import (
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
