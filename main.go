package main

import (
	"log"
	"net/http"

	"github.com/syrlramadhan/etika-reservation-api/config"
	"github.com/syrlramadhan/etika-reservation-api/controller"
	"github.com/syrlramadhan/etika-reservation-api/repository"
	"github.com/syrlramadhan/etika-reservation-api/service"
)

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	log.Println("Etika Reservation API")

	db := config.ConnectDB()
	repo := repository.NewReservationRepository(db)
	svc := service.NewReservationService(repo)
	ctrl := controller.NewReservationController(svc)

	http.HandleFunc("/api/reservations", withCORS(ctrl.CreateReservation))
	http.HandleFunc("/api/reservations/by-date", withCORS(ctrl.GetReservationsByDate))
	http.HandleFunc("/api/reservations/range", withCORS(ctrl.GetReservationsByDateRange))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
