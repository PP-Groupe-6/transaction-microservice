package transfer_microservice

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DbConnexionInfo struct {
	DbUrl    string
	DbPort   string
	DbName   string
	Username string
	Password string
}

func GetDbConnexion(info DbConnexionInfo) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "port="+info.DbPort+" user="+info.Username+" password="+info.Password+" dbname="+info.DbName+" sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
