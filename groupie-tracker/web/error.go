package web

import (
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status ErrStatus) {
	t, err := template.ParseFiles("./template/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
