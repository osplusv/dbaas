package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/osplusv/dbaas/handlers"
	"github.com/osplusv/dbaas/inmem"
	"github.com/osplusv/dbaas/provisioning"
)

type Handlers struct {
	Database *handlers.DatabaseHandler
}

var (
	containerRep = inmem.NewContainerRepository()
	databaseRep  = inmem.NewDatabaseRepository()
)

func NewServer() *echo.Echo {
	databaseService := provisioning.NewDatabaseService(containerRep, databaseRep)
	handlers := &Handlers{
		Database: handlers.NewDatabaseHandler(databaseService),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("api/v1")

	provisionAPI := api.Group("/provision")
	provisionAPI.GET("/database/:userid", handlers.Database.ListDatabases)
	provisionAPI.POST("/database", handlers.Database.ProvisionNewDatabase)
	provisionAPI.DELETE("/database/:userid/:databaseid", handlers.Database.UnprovisionDatabase)

	return e
}
