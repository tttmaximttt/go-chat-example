package chat

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")

	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	auth.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	fmt.Println(segs)
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusBadRequest,
			)
			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err),
				http.StatusBadRequest,
			)
			return
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginHandler(w, r)
}
