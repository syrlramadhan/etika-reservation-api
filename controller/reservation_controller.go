package controller

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/syrlramadhan/etika-reservation-api/config"
    "github.com/syrlramadhan/etika-reservation-api/dto"
    "github.com/syrlramadhan/etika-reservation-api/service"
)

type ReservationController struct {
    service service.ReservationService
}

func NewReservationController(service service.ReservationService) *ReservationController {
    return &ReservationController{service}
}

// saveImage menyimpan file gambar ke folder uploads dan mengembalikan path relatif
func saveImage(file io.Reader, filename string) (string, error) {
    // Validasi ekstensi file
    ext := strings.ToLower(filepath.Ext(filename))
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
        return "", fmt.Errorf("only JPG and PNG files are allowed")
    }

    // Tentukan direktori penyimpanan
    uploadDir := "./uploads"
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        return "", fmt.Errorf("failed to create upload directory: %v", err)
    }

    // Buat nama file unik menggunakan timestamp
    newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    filePath := filepath.Join(uploadDir, newFilename)

    // Simpan file
    out, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %v", err)
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return "", fmt.Errorf("failed to save file: %v", err)
    }

    // Kembalikan path relatif
    return filePath, nil
}

func (c *ReservationController) CreateReservation(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form dengan batas ukuran 10MB
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Ambil data dari form
    reservedDate := r.FormValue("reserved_date")
    eventName := r.FormValue("event_name")

    // Validasi input
    if reservedDate == "" || eventName == "" {
        http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
        return
    }

    // Ambil file gambar (opsional)
    var imageURL string
    file, handler, err := r.FormFile("image")
    if err == nil {
        defer file.Close()
        // Simpan gambar dan dapatkan URL
        imageURL, err = saveImage(file, handler.Filename)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    } else if err != http.ErrMissingFile {
        http.Error(w, "Failed to get image", http.StatusBadRequest)
        return
    }

    // Buat request DTO
    req := dto.CreateReservationRequest{
        ReservedDate: reservedDate,
        EventName:    eventName,
        ImageURL:     imageURL,
    }

    // Simpan reservasi
    id, err := c.service.CreateReservation(req)
    if err != nil {
        http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"message": "Reservation created", "id": id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (c *ReservationController) GetReservationsByDate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    date := r.URL.Query().Get("date")
    if date == "" {
        http.Error(w, "Date is required", http.StatusBadRequest)
        return
    }

    reservations, err := c.service.GetReservationsByDate(date)
    if err != nil {
        http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservations)
}

func (c *ReservationController) GetReservationsByDateRange(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    startDate := r.URL.Query().Get("start")
    endDate := r.URL.Query().Get("end")

    if startDate == "" || endDate == "" {
        http.Error(w, "Start and end dates are required", http.StatusBadRequest)
        return
    }

    reservations, err := c.service.GetReservationsByDateRange(startDate, endDate)
    if err != nil {
        http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reservations)
}

func (c *ReservationController) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Dummy user check
    if req.Username != "etika" || req.Password != "etika123" {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // JWT Payload
    claims := jwt.MapClaims{
        "username": req.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
        "iat":      time.Now().Unix(),
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(config.JwtSecret)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    resp := map[string]string{
        "message": "Login successful",
        "token":   signedToken,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}