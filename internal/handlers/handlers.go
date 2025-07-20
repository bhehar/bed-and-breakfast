package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/bhehar/bed-and-breakfast/internal/forms"
	"github.com/bhehar/bed-and-breakfast/internal/helpers"
	"github.com/bhehar/bed-and-breakfast/internal/models"
	"github.com/bhehar/bed-and-breakfast/internal/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers() sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handeles reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// validations
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLen("first_name", 3)
	form.ValidateEmail()

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the General's Quarters Room Page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Major's Suite room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles requests about room availability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("startDate")
	end := r.Form.Get("endDate")
	val := fmt.Sprintf("start date: %v -- end date: %v", start, end)
	w.Write([]byte(val))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Availability handles request about room availability and returns JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "available!"}

	out, err := json.MarshalIndent(resp, "", "	")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (rp *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary returns the summary of a reservation
func (rp *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		msg := "can't get reservation from session"
		// rp.App.InfoLog.Println(msg)
		// helpers.ServerError(w, errors.New(msg))
		rp.App.Session.Put(r.Context(), "error", msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	rp.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = res
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
