package helpers

import (
	"net/http"
	"runtime/debug"

	"github.com/bhehar/bed-and-breakfast/internal/config"
)

var app *config.AppConfig

// NewHelpers initializes app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError logs client errors using app.InfoLog
// & sends error response to client
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with status", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError logs server errors using app.ErrLog
// & sends error response to client
func ServerError(w http.ResponseWriter, err error) {
	app.ErrLog.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
