package chat

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"net/http"
	"strings"

	"github.com/stretchr/objx"
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

func callbackActionHandler(providerStr string, queryMap map[string]interface{}) (string, error) {
	provider, err := gomniauth.Provider(providerStr)
	if err != nil {
		return "", fmt.Errorf("Error when trying to get provider %s: %s", provider, err)
	}

	creds, err := provider.CompleteAuth(queryMap)
	if err != nil {
		return "", fmt.Errorf("Error when trying to complete auth for%s: %s", provider, err)
	}

	user, err := provider.GetUser(creds)
	if err != nil {
		return "", fmt.Errorf("Error when trying to get user from %s: %s", provider, err)
	}

	fmt.Println(user.AvatarURL())
	authCookieValue := objx.New(map[string]interface{}{
		"name":   user.Name(),
		"avatar": user.AvatarURL(),
	}).MustBase64()

	return authCookieValue, nil
}

func loginActionHandler(providerStr string) (string, error) {
	provider, err := gomniauth.Provider(providerStr)
	if err != nil {
		return "", fmt.Errorf("Error when trying to get provider %s: %s", provider, err)
	}

	loginUrl, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		return "", fmt.Errorf("Error when trying to GetBeginAuthURL for %s: %s", provider, err)
	}

	return loginUrl, nil
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
		loginUrl, err := loginActionHandler(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		queryMap := objx.MustFromURLQuery(r.URL.RawQuery)

		authCookieValue, err := callbackActionHandler(provider, queryMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginHandler(w, r)
}
