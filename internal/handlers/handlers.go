package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/models"
	"github.com/sunil206b/smart_booking/internal/render"
	"log"
	"net/http"
)

var Handler *RouteHandler

type RouteHandler struct {
	App *config.AppConfig
}

type jsonResponse struct {
	OK      bool   `json:"ok,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewRouteHandler(a *config.AppConfig) *RouteHandler {
	return &RouteHandler{
		App: a,
	}
}

func NewHandler(r *RouteHandler) {
	Handler = r
}

// Home renders the home page
func (rh *RouteHandler) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rh.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (rh *RouteHandler) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := rh.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the general quarters page
func (rh *RouteHandler) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the majors suit page
func (rh *RouteHandler) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the rooms available page
func (rh *RouteHandler) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
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
		OK: true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil{
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservations renders the contact page
func (rh *RouteHandler) Reservations(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (rh *RouteHandler) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}