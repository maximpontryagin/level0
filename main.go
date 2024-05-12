package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func main() {
	cfg := Config{
		Host:     "localhost",
		Port:     "5436",
		Username: "postgres",
		Password: "qwe",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ connected to books db")

	query := `
		drop table if exists books;
		create table if not exists books(
			id integer primary key,
			title text,
			author text,
			num_pages integer,
			rating real
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ created books table")
}
