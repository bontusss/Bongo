package session

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
)

type Session struct {
	CookieLifetime string
	CookiePersist  string
	CookieName     string
	CookieDomain   string
	SessionType    string
	CookieSecure   string
}

//todo: Write tests for sessions

func (s *Session) InitSession() *scs.SessionManager {
	var persist, secure bool
	session := scs.New()

	lifetime, err := strconv.Atoi(s.CookieLifetime)
	if err == nil {
		lifetime = 60
	}
	//session lifetime is cal
	session.Lifetime = time.Duration(lifetime) * time.Minute

	if s.CookiePersist == "true" {
		persist = true
	} else {
		persist = false
	}
	session.Cookie.Persist = persist

	if s.CookieSecure == "true" {
		secure = true
	} else {
		secure = false
	}
	session.Cookie.Secure = secure
	session.Cookie.Name = s.CookieName
	session.Cookie.Domain = s.CookieDomain
	session.Cookie.SameSite = http.SameSiteLaxMode

	switch s.SessionType {
	case "redis":
	//todo: implement redis store
	case "mysql", "mariadb":
		//todo: implement mysql store
	case "postgres", "postgresql":
		//todo: implement postgres store
	default:
		//	cookie
	}
	return session
}
