package quick_action

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-multierror"
)

// HttpErrorCallback handles errors from Quick Actions.
func HttpErrorCallback(w http.ResponseWriter, r *http.Request, err error) {
	errs, valid := err.(*multierror.Error)
	if !valid {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var errors []string
	for _, err := range errs.WrappedErrors() {
		errors = append(errors, err.Error())
	}

	json, _ := json.Marshal(map[string][]string{"errors": errors})

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(json)
}
