package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	DB, err = sqlx.Connect("postgres", "host=localhost port=5432 user=mikaeloganesan dbname=payments sslmode=disable")
	if err != nil {
		log.Fatalln("DB connection failed:", err)
	}
}
