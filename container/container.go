package container

import (
	"errors"

	"github.com/osplusv/dbaas/util"
)

type DatabaseContainer struct {
	ContainerID      string
	Image            string
	HostPort         string
	DatabaseServices int
	EnvCredential    util.Credential
}

type Repository interface {
	Store(container *DatabaseContainer) error
	Find(id string) (*DatabaseContainer, error)
	FindAll() []*DatabaseContainer
}

var ErrUnknown = errors.New("unknown database")
