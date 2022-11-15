package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Creationdate int      `json:"creationDate"`
	Firstalbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Concertdates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type Display struct {
	Name         string
	Image        string
	Member       []string
	Creationdate int
	Firstalbum   string
	Relate       map[string][]string
}

type Second struct {
	Index []struct {
		Id                int                 `json:"id"`
		Datesandlocations map[string][]string `json:"datesLocations"`
	}
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	tmp, err := template.ParseFiles("template/html.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rspData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/artists")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	jsonerr := json.Unmarshal(rspData, &artist)
	if jsonerr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, artist)
}

func ChangeAPI(s string) ([]byte, int) {
	response, err := http.Get(s)
	if err != nil {
		return nil, 500
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 500
	}
	return responseData, 0
}

func DisplayOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	myid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if myid > 52 || myid == 0 {
		http.NotFound(w, r)
		return
	}
	res, err := template.ParseFiles("template/result.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/artists")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	jsonerr := json.Unmarshal(responseData, &artist)
	if jsonerr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rspData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/relation")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var relate Second
	jsoner := json.Unmarshal(rspData, &relate)

	if jsoner != nil {
		fmt.Println(jsoner.Error())
	}
	myMap := make(map[string][]string)

	var members []string
	image := ""
	name := ""
	creationdate := 0
	firstalbum := ""
	for _, zn := range artist {
		for _, k := range relate.Index {
			if myid == zn.Id && myid == k.Id {
				name = zn.Name
				members = zn.Members
				image = zn.Image
				creationdate = zn.Creationdate
				firstalbum = zn.Firstalbum
				myMap = k.Datesandlocations
			}
		}
	}
	result := Display{
		Name:         name,
		Image:        image,
		Member:       members,
		Creationdate: creationdate,
		Firstalbum:   firstalbum,
		Relate:       myMap,
	}
	res.Execute(w, result)
}
