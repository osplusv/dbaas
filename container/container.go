package container

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/osplusv/dbaas/database"
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

func (c *DatabaseContainer) StartContainer() error {
	d, err := docker.New()
	if err != nil {
		return err
	}

	if err := d.StartContainer(c.ContainerID); err != nil {
		return err
	}

	return nil
}

func (c *DatabaseContainer) CreateNewDatabase() (*database.Database, error) {
	d, err := docker.New()
	if err != nil {
		return nil, err
	}
	fmt.Println("Creating database")

	dbName := strings.Split(strings.ToUpper(uuid.New()), "-")[0]
	user := strings.Split(strings.ToUpper(uuid.New()), "-")[0]
	pwd := strings.Split(strings.ToUpper(uuid.New()), "-")[0]

	cmd := []string{
		"/bin/bash",
		"-c",
		"mysql -u " + c.EnvCredential.Username + " -p" + c.EnvCredential.Password + " -e \"CREATE DATABASE " + dbName + ";CREATE USER '" + user + "'@'%' IDENTIFIED BY '" + pwd + "';GRANT ALL ON " + dbName + ".* TO '" + user + "'@'%';FLUSH PRIVILEGES;\"",
	}
	fmt.Println(cmd[2])
	if err := d.ExecCommand(c.ContainerID, cmd); err != nil {
		return nil, err
	}

	c.DatabaseServices++

	return &database.Database{
		ID:               strings.Split(strings.ToUpper(uuid.New()), "-")[0],
		Type:             c.Image,
		Name:             dbName,
		ContainerID:      c.ContainerID,
		EnvCredential:    util.Credential{Username: user, Password: pwd},
		ConnectionString: c.Image + "://" + user + ":" + pwd + "@" + "localhost:" + c.HostPort + "/" + dbName,
	}, nil
}

func (c *DatabaseContainer) DeleteDatabase(database *database.Database) error {
	d, err := docker.New()
	if err != nil {
		return err
	}

	fmt.Println("Deleting database")
	cmd := []string{
		"/bin/bash",
		"-c",
		"mysql -u " + c.EnvCredential.Username + " -p" + c.EnvCredential.Password + " -e \"DROP DATABASE " + database.Name + ";DROP USER '" + string(database.EnvCredential.Username) + "'@'%';\"",
	}
	fmt.Println(cmd[2])
	if err := d.ExecCommand(c.ContainerID, cmd); err != nil {
		return err
	}

	c.DatabaseServices--

	return nil
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
