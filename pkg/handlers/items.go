package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"moshack_2022/pkg/items"
	"moshack_2022/pkg/session"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo items.ItemsRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*items.Item
	}{
		Items: elems,
	})
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(item, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	sess, _ := session.SessionFromContext(r.Context())
	item.CreatedBy = sess.UserID

	lastID, err := h.ItemsRepo.Add(item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByID(uint32(id))
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", item)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(item, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	item.ID = uint32(id)

	ok, err := h.ItemsRepo.Update(item)
	if err != nil {
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", item, ok)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	ok, err := h.ItemsRepo.Delete(uint32(id))
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]bool{
		"success": ok,
	})
	w.Write(respJSON)
}
