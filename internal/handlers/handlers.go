package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/driver"
	"github.com/shwethadia/HotelReservation/internal/forms"
	"github.com/shwethadia/HotelReservation/internal/models"
	"github.com/shwethadia/HotelReservation/internal/render"
	"github.com/shwethadia/HotelReservation/internal/repository"
	"github.com/shwethadia/HotelReservation/internal/repository/dbrepo"
)

//Repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	return &Repository{

		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//Creates a new Test repository
func NewTestRepo(a *config.AppConfig) *Repository {

	return &Repository{

		App: a,
		DB:  dbrepo.NewPostgresTestRepo(a),
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {

	Repo = r

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.htm", &models.TemplateData{})
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "RemoteIp", remoteIP)

	render.Template(w, r, "home.page.htm", &models.TemplateData{})
}

//Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})

	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.htm", &models.TemplateData{

		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//PostReservation Handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	/* 	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	   	if !ok {

	   		helpers.ServerError(w, errors.New("can't get from session"))
	   		return
	   	} */

	err := r.ParseForm()
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse the start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse the end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Inavlid data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	/*
		reservation.FirstName = r.Form.Get("first_name")
		reservation.LastName = r.Form.Get("last_name")
		reservation.Phone = r.Form.Get("phone")
		reservation.Email = r.Form.Get("email") */

	reservation := models.Reservation{

		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "Error Message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.htm", &models.TemplateData{

			Form: form,
			Data: data,
		})

		return

	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{

		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//Generals renders the room page and displays form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.htm", &models.TemplateData{})

}

//Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.htm", &models.TemplateData{})

}

//Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.htm", &models.TemplateData{})
}

//Availability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	//w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse the start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse the end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get availability for roomsss")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {

		//No Availability
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms

	reservation_temp := models.Reservation{

		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation_temp)

	render.Template(w, r, "choose-room.page.htm", &models.TemplateData{

		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.PostForm
	if err != nil {

		resp := jsonResponse{

			OK:      false,
			Message: "Internal Server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	available, errr := m.DB.SearchAvailabilityByRoomID(startDate, endDate, roomID)
	if errr != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "   ")
		w.Header().Set("Content-type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
	//return

}

//Contact renders the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.htm", &models.TemplateData{})
}

//ReservationSummary displays the reservation summary page
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.htm", &models.TemplateData{

		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseRoom displays list of availabile rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	/*
		roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {

			helpers.ServerError(w, err)
			return
		}

		fmt.Println("*******************roomID", roomID)
		res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
		if !ok {
			helpers.ServerError(w, err)
			return
		} */

	//split the URL up by /, and grab the 3rd element
	element := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(element[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//BookNow takes URL parameters, builds a session variable and takes user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	//id ,s , e

	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	res.RoomID = ID

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	res.StartDate = startDate
	res.EndDate = endDate

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't get room from db")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
