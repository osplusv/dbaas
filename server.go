package main

import (
	"os"

	"github.com/osplusv/dbaas/api"
)

func main() {
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "1323"
	}

	server := api.NewServer()
	server.Logger.Fatal(server.Start(":" + serverPort))
}
