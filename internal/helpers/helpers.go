package helpers

import (
	"fmt"
	"github.com/sunil206b/smart_booking/internal/config"
	"net/http"
	"runtime/debug"
)

//postgres://mticitmc:b8D0BnW2KVi8O_32dVj6fdkB_BLTEd0v@batyr.db.elephantsql.com/mticitmc

var appConfig *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	appConfig = a
}

func ClientError(w http.ResponseWriter, status int) {
	appConfig.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
