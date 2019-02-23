package auth

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/dghubble/gologin"
	github "github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"golang.org/x/oauth2"
	githuboauth2 "golang.org/x/oauth2/github"
)

const (
	sessionName = "goreadme"
)

var (
	sessionSecret = os.Getenv("SESSION_SECRET")
	githubID      = os.Getenv("GITHUB_ID")
	githubSecret  = os.Getenv("GITHUB_SECRET")

	sessionStore   = sessions.NewCookieStore([]byte(sessionSecret), nil)
	sessionUserKey = "user"
)

func Handlers(domain string) (callback, login http.Handler) {
	oauth2Config := &oauth2.Config{
		ClientID:     githubID,
		ClientSecret: githubSecret,
		RedirectURL:  domain + "/github/callback",
		Endpoint:     githuboauth2.Endpoint,
	}

	config := gologin.DefaultCookieConfig

	callback = github.StateHandler(
		config,
		github.CallbackHandler(oauth2Config, issueSession(), nil))
	login = github.StateHandler(
		config,
		github.LoginHandler(oauth2Config, nil))
	return
}

// issueSession issues a cookie session after successful Github login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			logrus.Errorf("Getting user from context: %s", err)
			http.Error(w, "Failed", http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Save(w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// RequireLogin redirects unauthenticated users to the login route.
func RequireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// IsAuthenticated returns true if the user has a signed session cookie.
func IsAuthenticated(r *http.Request) bool {
	if _, err := sessionStore.Get(r, sessionName); err == nil {
		return true
	}
	return false
}

func ID(r *http.Request) string {
	s, err := sessionStore.Get(r, sessionName)
	if err != nil {
		logrus.Errorf("Failed getting user: %s", err)
		return ""
	}
	id, ok := s.Values[sessionUserKey].(string)
	if !ok {
		logrus.Errorf("Failed converting user key: %s", s.Values[sessionUserKey])
		return ""
	}
	return id
}
