package main

import (
	"net/http"

	"log"

	"github.com/syrlramadhan/etika-reservation-api/config"
	"github.com/syrlramadhan/etika-reservation-api/controller"
	"github.com/syrlramadhan/etika-reservation-api/repository"
	"github.com/syrlramadhan/etika-reservation-api/service"
)

func main() {
	log.Println("Etika Reservation API")
	db := config.ConnectDB()
	repo := repository.NewReservationRepository(db)
	svc := service.NewReservationService(repo)
	ctrl := controller.NewReservationController(svc)

	http.HandleFunc("/api/reservations", ctrl.CreateReservation)
	http.HandleFunc("/api/reservations/by-date", ctrl.GetReservationsByDate)

	http.ListenAndServe(":8080", nil)
}
