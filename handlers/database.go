package handlers

import "github.com/osplusv/dbaas/provisioning"

type DatabaseHandler struct {
	service provisioning.DatabaseProvisioner
}

func NewDatabaseHandler(service provisioning.DatabaseProvisioner) *DatabaseHandler {
	return &DatabaseHandler{service: service}
}
