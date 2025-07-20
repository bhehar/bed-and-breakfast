package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/bhehar/bed-and-breakfast/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in prod
	testApp.InProd = false
	// logger setup
	testApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.ErrLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	return http.Header{}
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
