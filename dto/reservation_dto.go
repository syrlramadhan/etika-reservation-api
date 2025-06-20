package dto

type CreateReservationRequest struct {
    ReservedDate string `json:"reserved_date"`
    EventName    string `json:"event_name"`
    Participants int    `json:"participants"`
    ImageURL     string `json:"image_url"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}