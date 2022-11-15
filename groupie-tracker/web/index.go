package web

import (
	"fmt"
	"group/models"
	"log"
	"net/http"
	"text/template"
)

var Links models.Link

const URL = "https://groupietrackers.herokuapp.com/api"

type ErrStatus struct {
	StatusCode int
	StatusText string
}

type GroupInfo struct {
	BaseInfo models.Group
	Relation models.Relations
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	fmt.Println("index Handler")
	t, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Println("parse fail")
		ErrorHandler(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}
	fmt.Println("get json")
	Artists := []models.Group{}
	err = models.GetJson(&Artists, Links.Artists)
	if err != nil {
		ErrorHandler(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}
	fmt.Printf("---------------\n%d\n-----------------\n", len(Artists))
	err = t.Execute(w, Artists)
	if err != nil {
		log.Println("index", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Println("finish index\n\n")
}
