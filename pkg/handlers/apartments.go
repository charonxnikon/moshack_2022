package handlers

import (
	"encoding/json"
	"html/template"
	"moshack_2022/pkg/apartments"
	"moshack_2022/pkg/apartments/excelParser"
	"moshack_2022/pkg/session"
	"net/http"

	"go.uber.org/zap"
)

type ApartmentHandler struct {
	Tmpl          *template.Template
	ApartmentRepo apartments.ApartmentRepo
	Logger        *zap.SugaredLogger
	Sessions      *session.SessionsManager
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

	h.Logger.Infof("header.Filename %v\n", header.Filename)
	h.Logger.Infof("header.Header %#v\n", header.Header)

	userSession, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	aparts, err := excelParser.ParseXLS(file, userSession.UserID)
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

	//w.Write(apartments.MarshalApartments(aparts))
	data, err := json.Marshal(&aparts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *ApartmentHandler) Table(w http.ResponseWriter, r *http.Request) {
	userSession, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	aparts, err := h.ApartmentRepo.GetAllByUserID(userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(&aparts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write(data)
}

func (h *ApartmentHandler) Estimate(w http.ResponseWriter, r *http.Request) {
	//apartmentID, err := strconv.Atoi(r.FormValue("id"))
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusUnauthorized)
	// 		return
	// 	}

	// здесь пишем JSON-rpc к питухону с запросом на посчитать и подобрать аналоги
	// по сути просто пересылаем id квартиры (квартир?)
}

func (h *ApartmentHandler) Reestimate(w http.ResponseWriter, r *http.Request) {
	// аналогично Estimate, только ещё всякие корректировки надо переслать
}

func (h *ApartmentHandler) EstimateAll(w http.ResponseWriter, r *http.Request) {
	// рассчитываем весь пулл
	// мб сразу формируем ексель и предлагаем скачать бесплатно без смс и регистрации?
}
