package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/bhehar/bed-and-breakfast/internal/driver"
	"github.com/bhehar/bed-and-breakfast/internal/forms"
	"github.com/bhehar/bed-and-breakfast/internal/helpers"
	"github.com/bhehar/bed-and-breakfast/internal/models"
	"github.com/bhehar/bed-and-breakfast/internal/render"
	"github.com/bhehar/bed-and-breakfast/internal/repository"
	"github.com/bhehar/bed-and-breakfast/internal/repository/dbrepo"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestingRepo creates a new Repository
func NewTestingRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers() sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (rp *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(r.Context(), "error", "could not load reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	strMap := make(map[string]string)
	strMap["start_date"] = sd
	strMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strMap,
	})
}

// PostReservation handeles reservation form
func (rp *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(r.Context(), "error", "could not load reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't parse form in request body")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	// validations
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLen("first_name", 3)
	form.ValidateEmail()

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res
		http.Error(w, "the form is not valid", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationId, err := rp.DB.InsertReservation(res)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "could not insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}
	err = rp.DB.InsertRoomRestriction(restriction)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "could not insert room restriction into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rp.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the General's Quarters Room Page
func (rp *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Major's Suite room page
func (rp *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the availability page
func (rp *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles requests about room availability
func (rp *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	startDate, err := time.Parse(time.DateOnly, r.Form.Get("start_date"))
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(time.DateOnly, r.Form.Get("end_date"))
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := rp.DB.SearchAvailabilityByDates(startDate, endDate)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		rp.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	rp.App.Session.Put(r.Context(), "reservation", reservation)
	data := make(map[string]any)
	data["rooms"] = rooms
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    int    `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Availability handles request about room availability and returns JSON
func (rp *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	startDate, _ := time.Parse(time.DateOnly, r.Form.Get("start_date"))
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	endDate, _ := time.Parse(time.DateOnly, r.Form.Get("end_date"))
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	isAvailable, err := rp.DB.SearchAvailabilityByDatesAndRoomId(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "error connecting to database",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        isAvailable,
		Message:   "available!",
		RoomId:    roomID,
		StartDate: startDate.Format(time.DateOnly),
		EndDate:   endDate.Format(time.DateOnly),
	}

	// ignore error since all aspects of json are handled within the function
	out, _ := json.MarshalIndent(resp, "", "	")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (rp *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
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

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	strMap := make(map[string]string)
	strMap["start_date"] = sd
	strMap["end_date"] = ed

	rp.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = res
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: strMap,
	})
}

func (rp *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(exploded[2])
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomName := r.URL.Query().Get("roomName")

	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomId
	res.Room = models.Room{
		RoomName: roomName,
	}
	rp.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (rp *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	var res models.Reservation
	res.RoomID = roomId

	res.StartDate, _ = time.Parse(time.DateOnly, startDate)
	res.EndDate, _ = time.Parse(time.DateOnly, endDate)

	room, err := rp.DB.GetRoomById(roomId)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room = room
	rp.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
