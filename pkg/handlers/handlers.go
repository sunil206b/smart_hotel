package handlers

import (
	"github.com/sunil206b/smart_booking/pkg/config"
	"github.com/sunil206b/smart_booking/pkg/models"
	"github.com/sunil206b/smart_booking/pkg/render"
	"net/http"
)

var Handler *RouteHandler

type RouteHandler struct {
	App *config.AppConfig
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
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (rh *RouteHandler) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := rh.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the general quarters page
func (rh *RouteHandler) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the majors suit page
func (rh *RouteHandler) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the rooms available page
func (rh *RouteHandler) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Reservations renders the contact page
func (rh *RouteHandler) Reservations(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (rh *RouteHandler) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}