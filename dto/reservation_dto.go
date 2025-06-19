package dto

type CreateReservationRequest struct {
	ReservedDate string `json:"reserved_date"`
	CustomerName string `json:"customer_name"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Notes        string `json:"notes"`
}