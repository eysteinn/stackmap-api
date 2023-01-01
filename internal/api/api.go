package api

import (
	"fmt"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/database"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/router"

	"log"
	"net/http"
)

func Run(host string) error {

	r := router.Setup()
	if database.GetDB() == nil {
		//return errors.New("Unable to connect to Database")
		fmt.Println("Unable to connect to Database")
	}

	log.Println("Listening at " + host)
	log.Fatal(http.ListenAndServe(host, r))
	return nil
}
