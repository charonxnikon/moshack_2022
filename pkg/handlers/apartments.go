package handlers

import (
	"encoding/json"
	"html/template"
	"moshack_2022/pkg/apartments"
	"moshack_2022/pkg/apartments/excelParser"
	"moshack_2022/pkg/session"
	"net/http"
	"net/rpc"
	"strconv"

	"go.uber.org/zap"
)

type ApartmentHandler struct {
	Tmpl          *template.Template
	ApartmentRepo apartments.ApartmentRepo
	Logger        *zap.SugaredLogger
	Sessions      *session.SessionsManager
	JSONrpcClient *rpc.Client
}

func (h *ApartmentHandler) Load(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "loadxls.html", nil)
	if err != nil {
		http.Error(w, "Template errror", http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	userSession, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	aparts, err := h.ApartmentRepo.GetAllUserApartmentsByUserID(userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	file, err := excelParser.UnparseXLSX(aparts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = file.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024 * 1024) // 128 MBytes
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
	aparts, err := excelParser.ParseXLSX(file, userSession.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, apart := range aparts {
		_, err := h.ApartmentRepo.AddUserApartment(apart)
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

	aparts, err := h.ApartmentRepo.GetAllUserApartmentsByUserID(userSession.UserID)
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
	apartmentID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	type Analogs struct {
		Analogs    []uint32
		Price      float64
		TotalPrice float64
	}
	var analogs Analogs
	h.JSONrpcClient.Call("get_analogs", apartmentID, &analogs)

	aparts := make([]*apartments.DBApartment, 0)
	for _, id := range analogs.Analogs {
		apart, err := h.ApartmentRepo.GetDBApartmentByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		aparts = append(aparts, apart)
	}

	type Result struct {
		Analogs    []*apartments.DBApartment
		PriceM2    float64
		TotalPrice float64
	}
	res := Result{
		Analogs:    aparts,
		PriceM2:    analogs.Price,
		TotalPrice: analogs.TotalPrice,
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write(data)
}

func (h *ApartmentHandler) Reestimate(w http.ResponseWriter, r *http.Request) {
	// аналогично Estimate, только ещё всякие корректировки надо переслать
}

func (h *ApartmentHandler) EstimateAll(w http.ResponseWriter, r *http.Request) {
	// рассчитываем весь пулл
	// мб сразу формируем ексель и предлагаем скачать бесплатно без смс и регистрации?
}
