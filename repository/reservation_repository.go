package repository

import (
    "database/sql"
    "github.com/syrlramadhan/etika-reservation-api/model"
)

type ReservationRepository interface {
    Save(reservation model.Reservation) error
    FindByDate(date string) ([]model.Reservation, error)
    FindByDateRange(startDate, endDate string) ([]model.Reservation, error)
}

type reservationRepository struct {
    db *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
    return &reservationRepository{db}
}

func (r *reservationRepository) Save(res model.Reservation) error {
    _, err := r.db.Exec(
        `INSERT INTO reservations (id, reserved_date, event_name, image_url)
        VALUES (?, ?, ?, ?)`,
        res.ID, res.ReservedDate, res.EventName, res.ImageURL,
    )
    return err
}

func (r *reservationRepository) FindByDate(date string) ([]model.Reservation, error) {
    rows, err := r.db.Query(
        `SELECT id, reserved_date, event_name, image_url, created_at 
        FROM reservations WHERE reserved_date = ?`, date)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reservations []model.Reservation
    for rows.Next() {
        var res model.Reservation
        if err := rows.Scan(&res.ID, &res.ReservedDate, &res.EventName, 
                           &res.ImageURL, &res.CreatedAt); err != nil {
            return nil, err
        }
        reservations = append(reservations, res)
    }
    return reservations, nil
}

func (r *reservationRepository) FindByDateRange(startDate, endDate string) ([]model.Reservation, error) {
    rows, err := r.db.Query(
        `SELECT id, reserved_date, event_name, image_url, created_at 
        FROM reservations 
        WHERE reserved_date BETWEEN ? AND ?`,
        startDate, endDate,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reservations []model.Reservation
    for rows.Next() {
        var res model.Reservation
        if err := rows.Scan(&res.ID, &res.ReservedDate, &res.EventName, 
                           &res.ImageURL, &res.CreatedAt); err != nil {
            return nil, err
        }
        reservations = append(reservations, res)
    }
    return reservations, nil
}