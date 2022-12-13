package acg_http

import "net/http"

func VerifyRequestType(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func VerifyFormRequiredFields(w http.ResponseWriter, r *http.Request) {
	requiredFields := []string{
		"user_id",
		"response_url",
	}
	for _, v := range requiredFields {
		if r.FormValue(v) == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func VerifyRequestForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	VerifyFormRequiredFields(w, r)
}
