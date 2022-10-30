package middleware

import (
	"net/http"

	"moshack_2022/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login":                     struct{}{},
		"/registration":              struct{}{},
		"/static/base/css/style.css": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/": struct{}{},
	}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
