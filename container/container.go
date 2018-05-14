package container

import (
	"errors"
	"strconv"
	"strings"

	"github.com/osplusv/dbaas/docker"
	"github.com/osplusv/dbaas/util"
	"github.com/pborman/uuid"
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

type ContainerSpecification struct {
	Image string
}

var currentPort = 55000
var ErrUnknown = errors.New("unknown database")

func NewDatabaseContainer(specifications ContainerSpecification) (*DatabaseContainer, error) {
	container, err := allocateNewContainer(specifications.Image)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func allocateNewContainer(image string) (*DatabaseContainer, error) {
	d, err := docker.New()
	if err != nil {
		return nil, err
	}

	secret := strings.Split(strings.ToUpper(uuid.New()), "-")[0]
	currentPort++
	newHostPort := strconv.Itoa(currentPort)

	containerID, err := d.CreateContainer(image, []string{
		"MYSQL_ROOT_PASSWORD=" + secret,
	}, newHostPort, secret)

	return &DatabaseContainer{
		ContainerID:      containerID,
		Image:            image,
		HostPort:         newHostPort,
		DatabaseServices: 0,
		EnvCredential: util.Credential{
			Username: "root",
			Password: secret,
		},
	}, nil
}
