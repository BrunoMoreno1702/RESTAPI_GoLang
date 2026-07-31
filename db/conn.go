package db

import (
	"database/sql"
	"fmt"
	"os"
)

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("a variavel de ambiente %s e obrigatoria", key)
	}

	return value, nil
}

func ConectDB() (*sql.DB, error) {
	host, err := getRequiredEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getRequiredEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	user, err := getRequiredEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := getRequiredEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbname, err := getRequiredEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	pqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", pqlInfo)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Conexao com o banco de dados realizada com sucesso!")
	return db, nil
}
