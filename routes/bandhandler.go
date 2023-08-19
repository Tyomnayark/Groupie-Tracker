package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
)

type Data struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func BandHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest)
		return
	}

	var band Artist
	isFound := false
	for _, artist := range Artists {
		if artist.ID == id {
			band = artist
			isFound = true
			break
		}
	}

	if !isFound {
		Error(w, http.StatusNotFound)
		return
	}

	response, err := http.Get(band.Relations)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		Error(w, http.StatusInternalServerError)
		return
	}

	var data Data
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	band.DatesLocations = data.DatesLocations

	tmpl, err := template.ParseFiles("./assets/bandPage.html")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, band)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
}
