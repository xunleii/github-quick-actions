package quick_action

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
)

// HttpErrorCallback handles errors from Quick Actions.
func HttpErrorCallback(w http.ResponseWriter, r *http.Request, err error) {
	errs, valid := err.(*multierror.Error)
	if !valid {
		// not handled errors
		return
	}

	var errors []string
	for _, err := range errs.WrappedErrors() {
		errors = append(errors, err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	json, err := json.Marshal(map[string][]string{"errors": errors})
	if err != nil {
		zerolog.Ctx(r.Context()).Error().Err(err).Send()
		w.Header().Set("Content-Type", "application/txt")
		json = []byte(errs.Error())
	}

	_, _ = w.Write(json)
	w.WriteHeader(http.StatusInternalServerError)
}
