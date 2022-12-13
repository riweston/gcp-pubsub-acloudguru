package acg_http

import "net/http"

func VerifyRequestType(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func VerifyRequestForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
