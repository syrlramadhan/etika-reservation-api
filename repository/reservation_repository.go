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
		`INSERT INTO reservations (id, reserved_date, customer_name, phone_number, email, notes)
		VALUES (?, ?, ?, ?, ?, ?)`,
		res.ID, res.ReservedDate, res.CustomerName, res.PhoneNumber, res.Email, res.Notes,
	)
	return err
}

func (r *reservationRepository) FindByDate(date string) ([]model.Reservation, error) {
	rows, err := r.db.Query(`SELECT id, reserved_date, customer_name, phone_number, email, notes FROM reservations WHERE reserved_date = ?`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		if err := rows.Scan(&res.ID, &res.ReservedDate, &res.CustomerName, &res.PhoneNumber, &res.Email, &res.Notes); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	return reservations, nil
}

func (r *reservationRepository) FindByDateRange(startDate, endDate string) ([]model.Reservation, error) {
	rows, err := r.db.Query(`
		SELECT id, reserved_date, customer_name, phone_number, email, notes 
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
		if err := rows.Scan(&res.ID, &res.ReservedDate, &res.CustomerName, &res.PhoneNumber, &res.Email, &res.Notes); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	return reservations, nil
}