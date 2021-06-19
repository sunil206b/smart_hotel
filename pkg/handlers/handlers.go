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

func (rh *RouteHandler) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rh.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (rh *RouteHandler) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := rh.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
