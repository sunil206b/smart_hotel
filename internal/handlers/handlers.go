package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/driver"
	"github.com/sunil206b/smart_booking/internal/forms"
	"github.com/sunil206b/smart_booking/internal/helpers"
	"github.com/sunil206b/smart_booking/internal/models"
	"github.com/sunil206b/smart_booking/internal/render"
	"github.com/sunil206b/smart_booking/internal/repository"
	"github.com/sunil206b/smart_booking/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiDateLayout         = "01/02/2006"
	restrictionDateLayout = "01/2/2006"
)

var Handler *RouteHandler

type RouteHandler struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type jsonResponse struct {
	OK        bool   `json:"ok,omitempty"`
	Message   string `json:"message,omitempty"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func NewRouteHandler(a *config.AppConfig, db *driver.DB) *RouteHandler {
	return &RouteHandler{
		App: a,
		DB:  dbrepo.NewPostgreRepo(db.SQL, a),
	}
}

func NewHandler(r *RouteHandler) {
	Handler = r
}

// Home renders the home page
func (rh *RouteHandler) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (rh *RouteHandler) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals renders the general quarters page
func (rh *RouteHandler) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the majors suit page
func (rh *RouteHandler) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the rooms available page
func (rh *RouteHandler) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability receives form data from the request and send available rooms data
func (rh *RouteHandler) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	startDate, err := time.Parse(apiDateLayout, start)
	if err != nil {
		rh.App.ErrorLog.Println("failed to convert start date", err)
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(apiDateLayout, end)
	if err != nil {
		rh.App.ErrorLog.Println("failed to convert end date", err)
		helpers.ServerError(w, err)
		return
	}
	rooms, err := rh.DB.SearchAllAvailableRooms(startDate, endDate)
	if err != nil {
		rh.App.ErrorLog.Println("failed to get rooms", err)
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		rh.App.Session.Put(r.Context(), "error", fmt.Sprintf("No Rooms Available within the range of %s-%s", start, end))
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		CheckInDate:  startDate,
		CheckOutDate: endDate,
	}

	rh.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AvailabilityJSON receives form data from the request and send JSON response
func (rh *RouteHandler) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("check_in_date")
	end := r.Form.Get("check_out_date")

	startDate, err := time.Parse(apiDateLayout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(apiDateLayout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := rh.DB.SearchAvailabilityByDatesByRoom(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: start,
		EndDate:   end,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation renders the contact page
func (rh *RouteHandler) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := rh.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from the session"))
		return
	}

	room, err := rh.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	rh.App.Session.Put(r.Context(), "reservation", res)

	checkinDate := res.CheckInDate.Format(apiDateLayout)
	checkoutDate := res.CheckOutDate.Format(apiDateLayout)

	stringMap := make(map[string]string)
	stringMap["check_in_date"] = checkinDate
	stringMap["check_out_date"] = checkoutDate
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (rh *RouteHandler) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rh.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get reservation from session"))
	}
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = rh.DB.CreateReservation(&reservation)
	if err != nil {
		rh.App.ErrorLog.Println("failed to create reservation ", err)
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.CheckInDate,
		EndDate:       reservation.CheckOutDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservation.ID,
		RestrictionID: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err = rh.DB.CreateRoomRestriction(&restriction)
	if err != nil {
		rh.App.ErrorLog.Println("failed to create room restriction ", err)
		helpers.ServerError(w, err)
		return
	}

	htmlMsg := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s, %s: <br>
		This is confirm your reservation from %s to %s.
`, reservation.FirstName, reservation.LastName, reservation.CheckInDate.Format(apiDateLayout), reservation.CheckInDate.Format(apiDateLayout))
	msg := &models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	rh.App.MailChan <- msg

	htmlMsg = fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		A reservation has been made for %s from %s to %s.
`, reservation.Room.RoomName, reservation.CheckInDate.Format(apiDateLayout), reservation.CheckInDate.Format(apiDateLayout))
	msg = &models.MailData{
		To:       "me@here.com",
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	rh.App.MailChan <- msg

	rh.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact renders the contact page
func (rh *RouteHandler) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

//ReservationSummary displays the reservation summary page
func (rh *RouteHandler) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rh.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rh.App.ErrorLog.Println("Cannot get item from session")
		rh.App.Session.Put(r.Context(), "error", "There are no reservations made at this point")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	checkinDate := reservation.CheckInDate.Format(apiDateLayout)
	checkoutDate := reservation.CheckOutDate.Format(apiDateLayout)

	stringMap := make(map[string]string)
	stringMap["check_in_date"] = checkinDate
	stringMap["check_out_date"] = checkoutDate

	rh.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseRoom will allow the user to select the room and send to make reservation page
func (rh *RouteHandler) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, ok := rh.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from the session"))
		return
	}
	res.RoomID = roomID
	rh.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

//BookRoom takes URL parameters and takes user to make reservation page
func (rh *RouteHandler) BookRoom(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	start := r.URL.Query().Get("start")
	startDate, err := time.Parse(apiDateLayout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	end := r.URL.Query().Get("end")
	endDate, err := time.Parse(apiDateLayout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := rh.DB.GetRoomByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res := new(models.Reservation)
	res.RoomID = ID
	res.CheckInDate = startDate
	res.CheckOutDate = endDate
	res.Room.RoomName = room.RoomName

	rh.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

// Login displays the login page
func (rh *RouteHandler) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin allows the user to login
func (rh *RouteHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = rh.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	id, _, err := rh.DB.Authenticate(email, password)
	if err != nil {
		rh.App.Session.Put(r.Context(), "error", "Invalid email or password")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	rh.App.Session.Put(r.Context(), "user_id", id)
	rh.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout allows the user to logout
func (rh *RouteHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = rh.App.Session.Destroy(r.Context())
	_ = rh.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashBoard displays reservation details in the admin tool
func (rh *RouteHandler) AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows all new reservations in the admin tool
func (rh *RouteHandler) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	rs, err := rh.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = rs
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations shows all reservations in the admin tool
func (rh *RouteHandler) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	rs, err := rh.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = rs
	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows single reservation in the admin page
func (rh *RouteHandler) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := rh.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (rh *RouteHandler) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]

	res, err := rh.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = rh.DB.UpdateReservation(&res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rh.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminReservationsCalender shows open calender rooms in the admin tool
func (rh *RouteHandler) AdminReservationsCalender(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	prev := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	prevMonth := prev.Format("01")
	prevMonthYear := prev.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["prev_month"] = prevMonth
	stringMap["prev_month_year"] = prevMonthYear
	stringMap["current_month"] = now.Format("01")
	stringMap["current_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := rh.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("01/2/2006")] = 0
			blockMap[d.Format("01/2/2006")] = 0
		}

		restrictions, err := rh.DB.GetRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, res := range restrictions {
			if res.ReservationID > 0 {
				for d := res.StartDate; d.After(res.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("01/2/2006")] = res.ReservationID
				}
			} else {
				blockMap[res.StartDate.Format("01/2/2006")] = res.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap

		rh.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calender.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminProcessReservation marks a reservation as processed
func (rh *RouteHandler) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := rh.DB.UpdateProcessedReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rh.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminDeleteReservation deletes a reservation
func (rh *RouteHandler) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := rh.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rh.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
