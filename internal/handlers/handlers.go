package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/forms"
	"github.com/shwethadia/HotelReservation/internal/helpers"
	"github.com/shwethadia/HotelReservation/internal/models"
	"github.com/shwethadia/HotelReservation/internal/render"
)

//Repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//Creates a new repository
func NewRepo(a *config.AppConfig) *Repository {

	return &Repository{

		App: a,
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {

	Repo = r

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about.page.htm", &models.TemplateData{})
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "RemoteIp", remoteIP)

	render.RenderTemplate(w, r, "home.page.htm", &models.TemplateData{})
}

//Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.htm", &models.TemplateData{

		Form: forms.New(nil),
		Data: data,
	})

}

//PostReservation Handles the posting of a reservation form
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

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.htm", &models.TemplateData{

			Form: form,
			Data: data,
		})

		return

	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//Generals renders the room page and displays form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "generals.page.htm", &models.TemplateData{})

}

//Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.htm", &models.TemplateData{})

}

//Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "search-availability.page.htm", &models.TemplateData{})
}

//Availability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      false,
		Message: "Internal server error",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {

		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(out)
	//return

}

//Contact renders the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.htm", &models.TemplateData{})
}

func (m *Repository) ReservatioinSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation summary from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.htm", &models.TemplateData{

		Data: data,
	})
}
