package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Artist struct {
	ID             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Relations      string   `json:"relations"`
	DatesLocations map[string][]string
}

var Artists []Artist

func MainHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		Error(w, http.StatusMethodNotAllowed)
		return
	}
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"

	response, err := http.Get(apiURL)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&Artists)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return

	}
	tmpl, err := template.ParseFiles("./assets/index.html")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Artists)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
}
