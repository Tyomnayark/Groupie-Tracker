package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Coordinate struct {
	Latitude   string   `json:"lat"`
	Longitude  string   `json:"lon"`
	Importance float64  `json:"importance"`
	Dates      []string `json:"-"`
}

func (c *Coordinate) GenerateGoogleMapsURLWithMarker() string {
	return fmt.Sprintf("https://www.google.com/maps?q=%s,%s&markers=%s,%s", c.Latitude, c.Longitude, c.Latitude, c.Longitude)
}

var geoCoderUrl = "https://geocode.maps.co/search"

func cleanLocation(location string) string {
	// Replace "_" and "-" with "+"
	location = strings.ReplaceAll(location, "_", "+")
	location = strings.ReplaceAll(location, "-", "+")

	// Capitalize the words
	words := strings.Fields(location)
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	return strings.Join(words, " ")
}

func getCoords(location string) (*Coordinate, error) {
	url := geoCoderUrl + fmt.Sprintf("?q=%s", location)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var coordinates []*Coordinate
	if err := json.Unmarshal(data, &coordinates); err != nil {
		return nil, err
	}
	sort.Slice(coordinates, func(i, j int) bool {
		return coordinates[i].Importance > coordinates[j].Importance
	})
	if len(coordinates) > 0 {
		return coordinates[0], nil
	}
	return nil, nil
}

func BandHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		Error(w, http.StatusBadRequest)
		return
	}

	firstChar := idStr[0]
	if firstChar == '0' {
		Error(w, http.StatusBadRequest)
		return
	}

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
	locationCoordinates := make(map[string]*Coordinate)
	for location, dates := range band.DatesLocations {
		locationCleaned := cleanLocation(location)
		coords, err := getCoords(locationCleaned)
		coords.Dates = dates
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		if coords == nil {
			Error(w, http.StatusNotFound)
			return
		}
		locationCoordinates[location] = coords
	}
	band.LocationsDatesCoords = locationCoordinates
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
