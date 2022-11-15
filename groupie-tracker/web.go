package main

import (
	"group/models"
	"group/web"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	myLink := "http://localhost:8080"
	log.Println(myLink)
	err := models.GetJson(&web.Links, web.URL)
	if err != nil {
		log.Println(err)
		return
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", web.IndexHandler)
	mux.HandleFunc("/artist/", web.ArtistHandler)

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Println("ListenAndServe: ", err)
		return
	}
}
