package chat

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"net/http"
	"strings"

	"crypto/md5"
	"github.com/stretchr/objx"
	"io"
)

type authHandler struct {
	next http.Handler
}

func (auth *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	auth.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func callbackActionHandler(providerStr string, queryMap map[string]interface{}) (string, error) {
	m := md5.New()

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

	io.WriteString(m, strings.ToLower(user.Email()))
	userId := fmt.Sprintf("%x", m.Sum(nil))

	authCookieValue := objx.New(map[string]interface{}{
		"userId": userId,
		"name":   user.Name(),
		"avatar": user.AvatarURL(),
		"email":  user.Email(),
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

// initialAuthHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func initialAuthHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
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

func InitialAuthHandler(w http.ResponseWriter, r *http.Request) {
	initialAuthHandler(w, r)
}
