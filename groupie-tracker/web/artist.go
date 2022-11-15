package web

import (
	"group/models"
	"net/http"
	"regexp"
	"text/template"
)

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/artist.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		ErrorHandler(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}

	path := r.URL.Path
	re := regexp.MustCompile(`artist\/ID=\d+$`)
	match := re.MatchString(path)
	if !match {
		ErrorHandler(w, r, ErrStatus{http.StatusNotFound, http.StatusText(http.StatusNotFound)})
		return
	}

	re = regexp.MustCompile(`\d+$`)
	groupID := re.FindString(path)

	var group models.Group
	groupURL := Links.Artists + "/" + groupID
	err = models.GetJson(&group, groupURL)
	if err != nil {
		ErrorHandler(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError) + err.Error()})
		return
	}

	relation, err := relationParse(groupID)
	data := GroupInfo{group, relation}
	err = t.Execute(w, data)
	if err != nil {
		ErrorHandler(w, r, ErrStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		return
	}
}

func relationParse(id string) (models.Relations, error) {
	var relation models.Relations
	url := Links.Relation + "/" + id
	err := models.GetJson(&relation, url)
	if err != nil {
		return relation, err
	}

	for k, v := range relation.DatesLocations {
		re := regexp.MustCompile(`(-|_)`)
		temp := re.ReplaceAllString(k, " ")
		relation.DatesLocations[temp] = v
		delete(relation.DatesLocations, k)
	}

	return relation, nil
}
