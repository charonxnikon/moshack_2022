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

// type AnalogsAdjustments struct {
// 	Id        uint32        `json:"id"`
// 	Analogs   []uint32      `json:"analogs"`
// 	Tender    [1]float64    `json:"tender"`
// 	Floor     [3][3]float64 `json:"floor"`
// 	Area      [6][6]float64 `json:"area"`
// 	Kitchen   [3][3]float64 `json:"kitchen"`
// 	Balcony   [2][2]float64 `json:"balcony"`
// 	Metro     [6][6]float64 `json:"metro"`
// 	Condition [3][3]float64 `json:"condition"`
// }

type AnalogIDAdjustments struct {
	AnalogID  uint32  `json:"analog_id"`
	Tender    float64 `json:"tender"`
	Floor     float64 `json:"floor"`
	Area      float64 `json:"area"`
	Kitchen   float64 `json:"kitchen"`
	Balcony   float64 `json:"balcony"`
	Metro     float64 `json:"metro"`
	Condition float64 `json:"condition"`
}

type AnalogApartmentAdjustments struct {
	Apartment   apartments.DBApartment
	Adjustments AnalogIDAdjustments
}

type AnalogsPrices struct {
	Analogs    []AnalogIDAdjustments `json:"analogs_with_adjustments"`
	PriceM2    float64               `json:"price_m2"`
	TotalPrice float64               `json:"total_price"`
}

func (h *ApartmentHandler) Load(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "loadxls.html", nil)
	if err != nil {
		h.Logger.Error(fmt.Errorf("error from Load handler: error from ExecuteTemplate: %s", err))
		http.Error(w, "Template errror", http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from Download handler: " + fmt.Sprintf(errf, args...)
	}

	userSession, err := h.Sessions.Check(r)
	if err != nil {
		h.Logger.Error(errWrapper("error from session.Check: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	aparts, err := h.ApartmentRepo.GetAllUserApartmentsByUserID(userSession.UserID)
	if err != nil {
		h.Logger.Error(errWrapper("error from GetAllUserApartmentsByUserID: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	file, err := excelParser.UnparseXLSX(aparts)
	if err != nil {
		h.Logger.Error(errWrapper("error from excelParser.UnparseXLSX: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = file.Write(w)
	if err != nil {
		h.Logger.Error(errWrapper("error from Write to file: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ApartmentHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from ParseFile handler: " + fmt.Sprintf(errf, args...)
	}

	err := r.ParseMultipartForm(128 * 1024 * 1024) // 128 MBytes
	if err != nil {
		h.Logger.Error(errWrapper("error from ParseMultipartForm: %s", err))
		http.Error(w, "File errror: file is too much", http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("xls_file")
	if err != nil {
		h.Logger.Error(errWrapper("error from FormFile: %s", err))
		http.Error(w, "File errror", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	h.Logger.Infof("header.Filename %v\n", header.Filename)
	h.Logger.Infof("header.Header %#v\n", header.Header)

	userSession, err := h.Sessions.Check(r)
	if err != nil {
		h.Logger.Error(errWrapper("error from sessions.Check: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	aparts, err := excelParser.ParseXLSX(file, userSession.UserID)
	if err != nil {
		h.Logger.Error(errWrapper("error from excelParser.ParseXLSX: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, apart := range aparts {
		_, err := h.ApartmentRepo.AddUserApartment(apart)
		if err != nil {
			h.Logger.Error(errWrapper("error from ApartmentRepo.AddUserApartment: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	go h.JSONrpcClient.Call("update_pull", &userSession.UserID)

	data, err := json.Marshal(&aparts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *ApartmentHandler) Table(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from Table handler: " + fmt.Sprintf(errf, args...)
	}

	userSession, err := h.Sessions.Check(r)
	if err != nil {
		h.Logger.Error(errWrapper("error from session.Check: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	aparts, err := h.ApartmentRepo.GetAllUserApartmentsByUserID(userSession.UserID)
	if err != nil {
		h.Logger.Error(errWrapper("error from GetAllUserApartmentsByUserID: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(&aparts)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Marshal: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write(data)
}

func (h *ApartmentHandler) Estimate(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from Estimate handler: " + fmt.Sprintf(errf, args...)
	}

	rData, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error(errWrapper("error from io.ReadAll: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	type ApartmentID struct {
		Id uint32
	}
	var apartmentID ApartmentID
	err = json.Unmarshal(rData, &apartmentID)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Unmarshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rpcResponse, err := h.JSONrpcClient.Call("get_analogs", &apartmentID.Id)
	if err != nil {
		h.Logger.Error(errWrapper("error from jsonrpc call: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rpcResponse.Result == nil {
		h.Logger.Error(errWrapper("error from jsonrpc call: result is nil"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, ok := rpcResponse.Result.(string)
	if !ok {
		h.Logger.Error(errWrapper("error from jsonrpc call: result is not string"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Estimate handler rpc - req: %v, res: %v", apartmentID, rpcResponse)

	var analogs AnalogsPrices
	analogs.Analogs = make([]AnalogIDAdjustments, 0)
	err = json.Unmarshal([]byte(data), &analogs)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Unmarshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aparts := make([]AnalogApartmentAdjustments, 0)
	for _, analog := range analogs.Analogs {
		apart, err := h.ApartmentRepo.GetDBApartmentByID(analog.AnalogID)
		if err != nil {
			h.Logger.Error(errWrapper("error from GetDBApartmentByID: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aparts = append(aparts, AnalogApartmentAdjustments{*apart, analog})
	}

	type Response struct {
		AnalogsWithAdjustments []AnalogApartmentAdjustments `json:"analogs_with_adjustments"`
		PriceM2                float64                      `json:"price_m2"`
		TotalPrice             float64                      `json:"total_price"`
	}

	res := Response{
		AnalogsWithAdjustments: aparts,
		PriceM2:                analogs.PriceM2,
		TotalPrice:             analogs.TotalPrice,
	}

	wData, err := json.Marshal(res)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Marshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(wData)
}

func (h *ApartmentHandler) Reestimate(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from Reestimate handler: " + fmt.Sprintf(errf, args...)
	}

	type Request struct {
		AnalogsWithAdjustments []AnalogApartmentAdjustments `json:"analogs_with_adjustments"`
	}
	var rParams Request
	rData, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error(errWrapper("error from io.ReadAll: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(rData, &rParams)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Unmarshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rpcResponse, err := h.JSONrpcClient.Call("recalculate_price_expert_flat", &rParams)
	if err != nil {
		h.Logger.Error(errWrapper("error from jsonrpc call: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rpcResponse.Result == nil {
		h.Logger.Error(errWrapper("error from jsonrpc call: result is nil"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, ok := rpcResponse.Result.(string)
	if !ok {
		h.Logger.Error(errWrapper("error from jsonrpc call: result is not string"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Reestimate handler rpc - req: %v, res: %v", rParams, rpcResponse)

	var analogs AnalogsPrices
	analogs.Analogs = make([]AnalogIDAdjustments, 0)
	err = json.Unmarshal([]byte(data), &analogs)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Unmarshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aparts := make([]AnalogApartmentAdjustments, 0)
	for _, analog := range analogs.Analogs {
		apart, err := h.ApartmentRepo.GetDBApartmentByID(analog.AnalogID)
		if err != nil {
			h.Logger.Error(errWrapper("error from GetDBApartmentByID: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aparts = append(aparts, AnalogApartmentAdjustments{*apart, analog})
	}

	type Response struct {
		AnalogsWithAdjustments []AnalogApartmentAdjustments `json:"analogs_with_adjustments"`
		PriceM2                float64                      `json:"price_m2"`
		TotalPrice             float64                      `json:"total_price"`
	}

	res := Response{
		AnalogsWithAdjustments: aparts,
		PriceM2:                analogs.PriceM2,
		TotalPrice:             analogs.TotalPrice,
	}

	wData, err := json.Marshal(res)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Marshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(wData)
}

func (h *ApartmentHandler) EstimateAll(w http.ResponseWriter, r *http.Request) {
	errWrapper := func(errf string, args ...interface{}) string {
		return "error from EstimateAll handler: " + fmt.Sprintf(errf, args...)
	}

	userSession, err := h.Sessions.Check(r)
	if err != nil {
		h.Logger.Error(errWrapper("error from session.Check: %s", err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	type UserIDAdjastments struct {
		UserID    uint32   `json:"user_id"`
		Samples   []uint32 `json:"samples"`
		Tender    float64  `json:"tender"`
		Floor     float64  `json:"floor"`
		Area      float64  `json:"area"`
		Kitchen   float64  `json:"kitchen"`
		Balcony   float64  `json:"balcony"`
		Metro     float64  `json:"metro"`
		Condition float64  `json:"condition"`
	}

	var rParams UserIDAdjastments
	rData, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error(errWrapper("error from io.ReadAll: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(rData, &rParams)
	if err != nil {
		h.Logger.Error(errWrapper("error from json.Unmarshal: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rParams.UserID = userSession.UserID

	rpcResponse, err := h.JSONrpcClient.Call("calculate_pull", &rParams)
	if err != nil {
		h.Logger.Error(errWrapper("error from jsonrpc call: %s", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("EstimateAll handler rpc - req: %v, res: %v", rParams, rpcResponse)
}
