package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode int
	ErrorText string
}

func Error(w http.ResponseWriter, code int) {
	var response ErrorResponse
	response.ErrorCode = code
	response.ErrorText = http.StatusText(code)

	w.WriteHeader(code)
	fmt.Println(code)
	fmt.Println(http.StatusText(code))
	tmpl, err := template.ParseFiles("./assets/errorPage.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, response)
}
