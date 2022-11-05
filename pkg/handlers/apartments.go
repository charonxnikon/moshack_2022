package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"moshack_2022/pkg/apartments"
	"moshack_2022/pkg/apartments/excelParser"
	"moshack_2022/pkg/session"
	"net/http"

	"github.com/ybbus/jsonrpc"
	"go.uber.org/zap"
)

type ApartmentHandler struct {
	Tmpl          *template.Template
	ApartmentRepo apartments.ApartmentRepo
	Logger        *zap.SugaredLogger
	Sessions      *session.SessionsManager
	JSONrpcClient jsonrpc.RPCClient
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
	type ApartmentID struct {
		Id uint32
	}
	var apartmentID ApartmentID
	rData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(rData, &apartmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Analogs struct {
		Analogs    []uint32 `json:"Analogs"`
		PriceM2    float64  `json:"PriceM2"`
		TotalPrice float64  `json:"TotalPrice"`
	}
	var analogs Analogs
	analogs.Analogs = make([]uint32, 0)
	response, err := h.JSONrpcClient.Call("get_analogs", &apartmentID.Id)
	if err != nil {
		fmt.Println(err) //
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, ok := response.Result.(string)
	if !ok {
		fmt.Println("not string - ", response.Result) //
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal([]byte(data), &analogs)
	if err != nil {
		fmt.Println(err)             //
		fmt.Println(response.Result) //
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Result from get_analogs:\n", analogs) //

	aparts := make([]*apartments.DBApartment, 0)
	for _, id := range analogs.Analogs {
		apart, err := h.ApartmentRepo.GetDBApartmentByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		PriceM2:    analogs.PriceM2,
		TotalPrice: analogs.TotalPrice,
	}

	wData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(wData)
}

func (h *ApartmentHandler) Reestimate(w http.ResponseWriter, r *http.Request) {
	// аналогично Estimate, только ещё всякие корректировки надо переслать
}

func (h *ApartmentHandler) EstimateAll(w http.ResponseWriter, r *http.Request) {
	// рассчитываем весь пулл
	// мб сразу формируем ексель и предлагаем скачать бесплатно без смс и регистрации?
}
