package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

type Reservation struct {
	ID     string
	Name   string
	Guests int
	Table  int
}

func NewID() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", bytes)
}

var reservations []Reservation

func main() {
	http.HandleFunc("GET /reservations/{id}", GetReservation)
	http.HandleFunc("PUT /reservations/{id}", PutReservation)
	http.HandleFunc("DELETE /reservations/{id}", DeleteReservation)
	http.HandleFunc("POST /reservations", PostReservation)
	http.ListenAndServe(":8090", nil)
}

func GetReservation(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	for _, reservation := range reservations {
		if reservation.ID != id {
			continue
		}

		if bytes, err := json.Marshal(reservation); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.WriteHeader(http.StatusOK)
			response.Write(bytes)
		}

		break
	}
}

func PostReservation(response http.ResponseWriter, request *http.Request) {
	/*var body struct{
		Entries []struct {
			Name string
			Guests int
			Table int
		}
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, entry := range body.Entries {
		reservations = append(reservations, Reservation{})
	}*/
}

func PutReservation(response http.ResponseWriter, request *http.Request) {

}

func DeleteReservation(response http.ResponseWriter, request *http.Request) {

}
