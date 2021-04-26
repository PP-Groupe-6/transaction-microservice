package main

import (
	"net/http"
	"os"

	transferService "github.com/PP-Groupe-6/transfer-microservice/transfer_microservice"
	"github.com/go-kit/kit/log"
)

func main() {
	info := transferService.DbConnexionInfo{
		DbUrl:    "postgre://",
		DbPort:   "5432",
		DbName:   "prix_banque_test",
		Username: "dev",
		Password: "dev",
	}

	service := transferService.NewTransferService(info)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err := http.ListenAndServe(":8001", transferService.MakeHTTPHandler(service, logger))
	if err != nil {
		panic(err)
	}
}
