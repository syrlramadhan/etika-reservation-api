package service

import (
	"github.com/google/uuid"
	"github.com/syrlramadhan/etika-reservation-api/dto"
	"github.com/syrlramadhan/etika-reservation-api/model"
	"github.com/syrlramadhan/etika-reservation-api/repository"
)

type ReservationService interface {
	CreateReservation(req dto.CreateReservationRequest) (string, error)
	GetReservationsByDate(date string) ([]model.Reservation, error)
}

type reservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(repo repository.ReservationRepository) ReservationService {
	return &reservationService{repo}
}

func (s *reservationService) CreateReservation(req dto.CreateReservationRequest) (string, error) {
	id := uuid.New().String()
	reservation := model.Reservation{
		ID:           id,
		ReservedDate: req.ReservedDate,
		CustomerName: req.CustomerName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Notes:        req.Notes,
	}
	return id, s.repo.Save(reservation)
}

func (s *reservationService) GetReservationsByDate(date string) ([]model.Reservation, error) {
	return s.repo.FindByDate(date)
}
