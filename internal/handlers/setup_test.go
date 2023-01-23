package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var testSession *scs.SessionManager

func TestMain(m *testing.M) {
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

}
