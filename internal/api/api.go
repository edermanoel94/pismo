package api

import (
	"github.com/edermanoel94/pismo/internal/api/account"
	ad "github.com/edermanoel94/pismo/internal/api/account/data"
	at "github.com/edermanoel94/pismo/internal/api/account/transport"
	"github.com/edermanoel94/pismo/internal/infra/database"
	"github.com/edermanoel94/pismo/internal/infra/server"
)

func Start() error {

	s := server.New()

	db, err := database.New()

	if err != nil {
		return err
	}

	accRepository := ad.NewAccountRepository(db)

	accService := account.New(accRepository)

	at.NewHTTP(accService, s)

	server.Start(s)

	return nil
}
