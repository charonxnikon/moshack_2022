package handlers

import (
	"html/template"
	"net/http"

	"moshack_2022/pkg/session"
	"moshack_2022/pkg/user"

	"go.uber.org/zap"
)

const (
	errorNoUser  = "nouser"
	errorBadPass = "badpass"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
	Sessions *session.SessionsManager
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)

	// _, err := session.SessionFromContext(r.Context())
	// if err == nil {
	// 	http.Redirect(w, r, "/loadxls", http.StatusFound)
	// 	return
	// }

	// // так или как закомментированно ниже?
	// err = h.Tmpl.ExecuteTemplate(w, "login.html", nil)

	// //	err = renderer.Render(h.Tmpl, "mainpage.html", w, nil) // ?
	// if err != nil {
	// 	errString := fmt.Sprintf("Template error: %s", err)
	// 	http.Error(w, errString, http.StatusInternalServerError)
	// 	return
	// }
}

func (h *UserHandler) LoginGET(w http.ResponseWriter, r *http.Request) {
	errorString := r.FormValue("error")
	err := h.Tmpl.ExecuteTemplate(w, "login.html", errorString)
	if err != nil {
		http.Error(w, "Template errror", http.StatusInternalServerError)
	}
}

func (h *UserHandler) LoginPOST(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Authorize(r.FormValue("login"), r.FormValue("password"))
	if err == user.ErrNoUser {
		//http.Error(w, `no user`, http.StatusBadRequest)
		http.Redirect(w, r, "/login?error="+errorNoUser, http.StatusFound)
		return
	}
	if err == user.ErrBadPass {
		//http.Error(w, `bad pass`, http.StatusBadRequest)
		http.Redirect(w, r, "/login?error="+errorBadPass, http.StatusFound)
		return
	}

	sess, _ := h.Sessions.Create(w, u.ID)
	h.Logger.Infof("created session for %v", sess.UserID)
	//http.Redirect(w, r, "/", http.StatusFound)
	http.Redirect(w, r, "/loadxls", http.StatusFound)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "registration.html", nil)
	if err != nil {
		http.Error(w, "Template errror", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := h.UserRepo.AddUser(r.FormValue("login"), r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
