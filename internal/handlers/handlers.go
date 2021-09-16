package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/driver"
	"github.com/shwethadia/HotelReservation/internal/forms"
	"github.com/shwethadia/HotelReservation/internal/helpers"
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

	/* 	restriction := models.RoomRestriction{

		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	} */

	restriction := models.RoomRestriction{

		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomId,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notifications to user

	htmlMessage := fmt.Sprintf(`
	
			<strong>Reservation Confirmation</strong>
			Dear %s:<br>
			This is confirm your reservation from %s to %s.
			`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{

		To:       reservation.Email,
		From:     "shwetha@gmail.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.htm",
	}

	m.App.MailChan <- msg

	//send notification to property owner

	htmlMessage = fmt.Sprintf(`
			<strong>Reservation Notification</strong>
			A reservation has been made for %s from  %s to %s.
			`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.MailData{

		To:      "shwetha@gmail.com",
		From:    "shwetha@gmail.com",
		Subject: "Reservation Notification",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

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

//ShowLoogin shows the login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.page.htm", &models.TemplateData{

		Form: forms.New(nil),
	})
}

//PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {

		//Take user back to page
		render.Template(w, r, "login.page.htm", &models.TemplateData{

			Form: form,
		})

		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in suucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

//Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-dashboard.page.htm", &models.TemplateData{})
}

//AdminNewReservations shows all new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.htm", &models.TemplateData{

		Data: data,
	})
}

//AdminAllReservations lists all the reservations in admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.htm", &models.TemplateData{

		Data: data,
	})
}

//AdminShowReservation shows the reservation in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	exploaded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploaded[4])
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	log.Println(id)

	src := exploaded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	//Get reservations from database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.htm", &models.TemplateData{

		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	//Assume that there is no month /year specified

	now := time.Now()

	if r.URL.Query().Get("y") != "" {

		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	fmt.Println(rooms)

	for _, x := range rooms {

		//Create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {

			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0

		}

		//get all the resetrictions for the current room

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)

		if err != nil {

			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {

			if y.ReservationID > 0 {

				//Its a reservation

				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {

					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {

				//Its a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)

	}

	render.Template(w, r, "admin-reservations-calendar.page.htm", &models.TemplateData{

		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//AdminPostShowReservation updates the reservation by ID
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploaded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploaded[4])
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	log.Println(id)

	src := exploaded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//AdminProcessReservation Marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	err := m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		log.Println(err)
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

//AdminDeleteReservation Marks a reservation as processed
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = m.DB.DeleteReservation(id)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {

		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminPostReservationsCalendar handles post of reservation calender
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {

		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	//Process Blocks

	rooms, err := m.DB.AllRooms()
	if err != nil {

		helpers.ServerError(w, err)
		return

	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {

		//Get the block maps.Loop through entire map, if we have an entry in the map
		//That does not exist in our posted data, and if the restriction id >0, then it is a block to remove

		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {

			//ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {

				//Only pay attention to values >0 and that are not in the form post
				//The rest are just placeholders for days without blocks

				if val > 0 {

					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {

						//delete the restrictions by id

						err := m.DB.DeleteBlockForRoom(value)
						if err != nil {

							log.Println(err)
						}

					}
				}
			}
		}

	}

	//now handle new blocks
	for name, _ := range r.PostForm {

		if strings.HasPrefix(name, "add_block") {

			exploded := strings.Split(name, "_")

			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])

			//insert a new block

			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {

				log.Println(err)
			}

		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)

}
