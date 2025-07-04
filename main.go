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
        w.Header().Set("Access-Control-Allow-Origin", "https://etika.studio")
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

    // Tambahkan handler untuk melayani file statis (gambar)
    http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

    http.HandleFunc("/api/reservations", withCORS(ctrl.CreateReservation))
    http.HandleFunc("/api/reservations/by-date", withCORS(ctrl.GetReservationsByDate))
    http.HandleFunc("/api/reservations/range", withCORS(ctrl.GetReservationsByDateRange))
    http.HandleFunc("/api/login", withCORS(ctrl.Login))

    log.Fatal(http.ListenAndServe(":8080", nil))
}