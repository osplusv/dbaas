package database

import (
	"errors"

	"github.com/osplusv/dbaas/util"
)

type (
	Database struct {
		ID               string
		Name             string
		Type             string
		ConnectionString string
		EnvCredential    util.Credential
		UserID           string
		ContainerID      string
	}

	DatabaseSpecification struct {
		Type string
	}
)

type Repository interface {
	Store(cargo *Database) error
	Find(id string) (*Database, error)
	FindAll() []*Database
	Delete(id string) error
}

var ErrUnknown = errors.New("unknown database")
