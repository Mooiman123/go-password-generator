package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Database verbinden
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Fout bij verbinden met database:", err)
	}
	defer db.Close()

	// Tabel maken als die nog niet bestaat
	err = createTable()
	if err != nil {
		log.Fatal("Fout bij maken van tabel:", err)
	}

	// HTTP server starten
	http.HandleFunc("/generate", generatePasswordHandler)
	fmt.Println("Server draait op poort 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Database verbinding maken
func connectDB() (*sql.DB, error) {
	dbhost := getEnv("DB_HOST", "localhost")
	dbport := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "postgres")
	dbuser := getEnv("DB_USER", "postgres")
	dbpass := getEnv("DB_PASS", "mysecretpassword")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)
	return sql.Open("postgres", connStr)
}

// Omgevingsvariabelen ophalen
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Database tabel maken
func createTable() error {
	query := `CREATE TABLE IF NOT EXISTS passwords (
		id SERIAL PRIMARY KEY,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

// Wachtwoord genereren
func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}

// HTTP handler om wachtwoord te genereren en op te slaan
func generatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password, err := generatePassword(12)
	if err != nil {
		http.Error(w, "Fout bij genereren van wachtwoord", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO passwords (password) VALUES ($1)", password)
	if err != nil {
		http.Error(w, "Fout bij opslaan in database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Gegenereerd wachtwoord: %s\n", password)
}
