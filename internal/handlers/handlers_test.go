package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shwethadia/HotelReservation/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{

	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {

			t.Errorf("for %s , expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}

}

func TestRepository_Reservation(t *testing.T) {

	reservation := models.Reservation{

		RoomID: 1,
		Room: models.Room{

			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session(reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCTX(req *http.Request) context.Context {

	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {

		log.Println(err)
	}

	return ctx

}

func TestRepository_PostReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Reservation handler returned wrong response code got for missing  %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid start date
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Post Reservation handler returned wrong response code got for invalid start date:got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid end date
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "invalid")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Post Reservation handler returned wrong response code got for invalid end date:got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid room Id

	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Post Reservation handler returned wrong response code got for invalid room Id:got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid data
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "sh")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {

		t.Errorf("Post Reservation handler returned wrong response code got for invalid data:got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for failure to insert reservation into database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Post Reservation handler failed when trying to fail inserting reservation :got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for failure to insert restriction into database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "shwetha")
	postedData.Add("last_name", "dia")
	postedData.Add("email", "s@gmail.com")
	postedData.Add("phone", "5555-555-555")
	postedData.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Post Reservation handler failed when trying to fail inserting reservation :got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	//first case - rooms are not available

	postData := url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-02")

	//Create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))

	//get context with session
	ctx := getCTX(req)

	req = req.WithContext(ctx)

	//Set the request header
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	//make hander handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	//get response recorder
	rr := httptest.NewRecorder()

	//make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse

	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {

		t.Error(("failed to parse json"))
	}

}
