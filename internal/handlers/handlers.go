package handlers

import (
	"encoding/json"
	"fmt"
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
	"time"
)

const (
	apiDateLayout = "01/02/2006"
)

var Handler *RouteHandler

type RouteHandler struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type jsonResponse struct {
	OK      bool   `json:"ok,omitempty"`
	Message string `json:"message,omitempty"`
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
	start := r.Form.Get("startDate")
	end := r.Form.Get("endDate")
	w.Write([]byte(fmt.Sprintf("Start Date is %s and End Date is %s\n", start, end)))
}

// AvailabilityJSON receives form data from the request and send JSON response
func (rh *RouteHandler) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (rh *RouteHandler) PostReservation(w http.ResponseWriter, r *http.Request) {
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
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	cid := r.Form.Get("check_in_date")
	date, err := time.Parse(apiDateLayout, cid)
	if err != nil {
		rh.App.ErrorLog.Println("failed to convert checkin date", err)
		helpers.ServerError(w, err)
		return
	}
	reservation.CheckInDate = date
	cod := r.Form.Get("check_out_date")
	date, err = time.Parse(apiDateLayout, cod)
	if err != nil {
		rh.App.ErrorLog.Println("failed to convert checkin date", err)
		helpers.ServerError(w, err)
		return
	}
	reservation.CheckOutDate = date
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		rh.App.ErrorLog.Println("failed to convert room id", err)
		helpers.ServerError(w, err)
		return
	}
	reservation.RoomID = roomID
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

	rh.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact renders the contact page
func (rh *RouteHandler) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (rh *RouteHandler) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rh.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rh.App.ErrorLog.Println("Cannot get item from session")
		rh.App.Session.Put(r.Context(), "error", "There are no reservations made at this point")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rh.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
