package transfer_microservice

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dbConnexionInfo struct {
	dbUrl    string
	dbPort   string
	dbName   string
	username string
	password string
}

func GetDbConnexion(info dbConnexionInfo) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "port="+info.dbPort+" user="+info.username+" password="+info.password+" dbname="+info.dbName+" sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
