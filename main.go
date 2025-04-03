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
		log.Fatal("Fout bij laden van databaseconfiguratie:", err)
	}

	db, err := connectDB(config)
	if err != nil {
		log.Fatal("Fout bij verbinden met database:", err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal("Fout bij maken van tabel:", err)
	}

	password, err := generatePassword(12)
	if err != nil {
		log.Fatal("Fout bij genereren van wachtwoord:", err)
	}

	if _, err := db.Exec("INSERT INTO passwords (password) VALUES ($1)", password); err != nil {
		log.Fatal("Fout bij opslaan in database:", err)
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
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	return string(password), nil
}
