package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Artist struct {
	ID                   int      `json:"id"`
	Image                string   `json:"image"`
	Name                 string   `json:"name"`
	Members              []string `json:"members"`
	CreationDate         int      `json:"creationDate"`
	FirstAlbum           string   `json:"firstAlbum"`
	Relations            string   `json:"relations"`
	DatesLocations       map[string][]string
	LocationsDatesCoords map[string]*Coordinate
	Geo                  []string
}
type Location struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

var Artists []Artist
var Locations Location

func MainHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		Error(w, http.StatusMethodNotAllowed)
		return
	}
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	apiLoc := "https://groupietrackers.herokuapp.com/api/locations"

	responseLoc, err := http.Get(apiLoc)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	defer responseLoc.Body.Close()

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

	if responseLoc.StatusCode != http.StatusOK {
		Error(w, http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&Artists)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return

	}
	err = json.NewDecoder(responseLoc.Body).Decode(&Locations)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return

	}
	for i := range Artists {
		artist := &Artists[i]
		for _, index := range Locations.Index {
			if index.ID == artist.ID {
				artist.Geo = append(artist.Geo, index.Locations...)
				break
			}
		}
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
