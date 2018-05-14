package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/osplusv/dbaas/database"
	"github.com/osplusv/dbaas/provisioning"
)

type DatabaseHandler struct {
	service provisioning.DatabaseProvisioner
}

type (
	provisionDatabaseRequest struct {
		UserID       string `json:"user_id" validate:"required"`
		DatabaseType string `json:"database_type" validate:"required"`
	}
	provisionDatabaseResponse struct {
		Database databaseResponse `json:"database"`
	}
)

type databaseResponse struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
}

func NewDatabaseHandler(service provisioning.DatabaseProvisioner) *DatabaseHandler {
	return &DatabaseHandler{service: service}
}

func (r *DatabaseHandler) ProvisionNewDatabase(c echo.Context) error {
	request := provisionDatabaseRequest{}
	if err := bindPayload(&request, c); err != nil {
		return respondWithError(http.StatusBadRequest, err.Error(), c)
	}
	request.DatabaseType = strings.ToLower(request.DatabaseType)

	if !isDatabaseAvailable(request.DatabaseType) {
		return respondWithError(http.StatusBadRequest, errors.New("sorry, we dont support that database yet").Error(), c)
	}

	database, err := r.service.ProvisionNewDatabase(request.UserID, database.DatabaseSpecification{Type: request.DatabaseType})
	if err != nil {
		return respondWithError(http.StatusInternalServerError, err.Error(), c)
	}

	resp := &provisionDatabaseResponse{Database: databaseResponse{ID: database.ID, Type: request.DatabaseType, ConnectionString: database.ConnectionString}}
	return respondWithPayload(http.StatusCreated, resp, c)
}
