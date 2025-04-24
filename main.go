package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	_ "github.com/lib/pq"
)

const (
	errLoadConfig     = "Fout bij laden van databaseconfiguratie:"
	errDBConnection   = "Fout bij verbinden met database:"
	errCreateTable    = "Fout bij maken van tabel:"
	errGeneratePass   = "Fout bij genereren van wachtwoord:"
	errInsertPassword = "Fout bij opslaan in database:"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func main() {
	config, err := loadDBConfig("config.json")
	if err != nil {
		log.Fatal(errLoadConfig, err)
	}

	db, err := connectDB(config)
	if err != nil {
		log.Fatal(errDBConnection, err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal(errCreateTable, err)
	}

	password, err := generatePassword(12)
	if err != nil {
		log.Fatal(errGeneratePass, err)
	}

	if _, err := db.Exec("INSERT INTO passwords (password) VALUES ($1)", password); err != nil {
		log.Fatal(errInsertPassword, err)
	}

	fmt.Println("Gegenereerd wachtwoord:", password)
}

func loadDBConfig(filename string) (*DBConfig, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config DBConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func connectDB(config *DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	return sql.Open("postgres", connStr)
}

func createTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS passwords (id SERIAL PRIMARY KEY, password TEXT NOT NULL);")
	return err
}

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, length)
	for i := range password { // maakt wachtwoord aan vanuit de charset
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset)))) // kiest getal van 0 -charset
		if err != nil {                                                    // loopt die door elke positie in het wachtwoord
			return "", err
		} // controleert of er iets misgaat
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}
