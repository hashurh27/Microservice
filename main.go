package main

import (
	"Warranty-Microservice/api/forms"
	"net/http"
)

func main() {
	//http.HandleFunc("/hi", forms.CarSender)

	http.HandleFunc("/hi", forms.CustomerHandler)
	http.ListenAndServe(":8080", nil)

}
