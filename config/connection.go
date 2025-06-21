package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dsn := "root:Sementara123!@tcp(103.84.207.100:3306)/etika_reservation"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
	return db
}

var JwtSecret = []byte(getSecret())

func getSecret() string {
	// Bisa juga dari env atau config file
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret" // fallback
	}
	return secret
}