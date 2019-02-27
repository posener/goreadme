package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	gogithub "github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	githuboauth2 "golang.org/x/oauth2/github"
)

const (
	sessionName    = "goreadme"
	sessionUserKey = "user"
)

var (
	cookieConfig = gologin.DefaultCookieConfig
	debug        = os.Getenv("DEBUG") == "1"
)

func init() {
	if debug {
		logrus.Warnf("Auth in debug mode!")
		cookieConfig = gologin.DebugOnlyCookieConfig
	}
}

type Auth struct {
	SessionSecret string
	GithubID      string
	GithubSecret  string
	RedirectUrl   string
	LoginPath     string
	HomePath      string

	sessionStore *sessions.CookieStore
}

func (a *Auth) Init() {
	a.sessionStore = sessions.NewCookieStore([]byte(a.SessionSecret), nil)
}

func (a *Auth) CallbackHandler() http.Handler {
	return github.StateHandler(cookieConfig,
		github.CallbackHandler(a.config(), http.HandlerFunc(a.loginSuccess), http.HandlerFunc(a.loginFailed)))
}

func (a *Auth) LoginHandler() http.Handler {
	return github.StateHandler(cookieConfig, github.LoginHandler(a.config(), nil))
}

// loginSuccess issues a cookie session after successful Github login
func (a *Auth) loginSuccess(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Login succeeded")
	u, err := github.UserFromContext(r.Context())
	if err != nil {
		logrus.Errorf("Getting user from context: %s", err)
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(u)
	if err != nil {
		logrus.Errorf("Marshaling user: %+v: %s", u, err)
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}
	logrus.Infof("UserData: %s", string(b))

	session := a.sessionStore.New(sessionName)
	session.Values[sessionUserKey] = string(b)
	session.Save(w)
	http.Redirect(w, r, a.HomePath, http.StatusFound)
}

func (a *Auth) loginFailed(w http.ResponseWriter, r *http.Request) {
	err := gologin.ErrorFromContext(r.Context())
	logrus.Infof("Login failed: %s", err)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (a *Auth) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.GithubID,
		ClientSecret: a.GithubSecret,
		RedirectURL:  a.RedirectUrl,
		Endpoint:     githuboauth2.Endpoint,
	}
}

// RequireLogin redirects unauthenticated users to the login route.
func (a *Auth) RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !a.IsAuthenticated(r) {
			http.Redirect(w, r, a.LoginPath, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// IsAuthenticated returns true if the user has a signed session cookie.
func (a *Auth) IsAuthenticated(r *http.Request) bool {
	_, err := a.sessionStore.Get(r, sessionName)
	return err == nil
}

func (a *Auth) User(r *http.Request) *gogithub.User {
	s, err := a.sessionStore.Get(r, sessionName)
	if err != nil {
		logrus.Errorf("Failed getting user: %s", err)
		return nil
	}
	jsonData, ok := s.Values[sessionUserKey].(string)
	if !ok {
		logrus.Errorf("Failed converting user key: %s", s.Values[sessionUserKey])
		return nil
	}
	var u gogithub.User
	err = json.Unmarshal([]byte(jsonData), &u)
	if err != nil {
		logrus.Errorf("Failed marhsalling user data %s: %s", jsonData, err)
		return nil
	}

	return &u
}
