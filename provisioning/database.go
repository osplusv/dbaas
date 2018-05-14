package provisioning

import (
	"errors"

	"github.com/osplusv/dbaas/container"
	"github.com/osplusv/dbaas/database"
)

type DatabaseProvisioner interface {
	ProvisionNewDatabase(id string, specification database.DatabaseSpecification) (*database.Database, error)
	ProvisionedDatabases(id string) ([]*database.Database, error)
	UnprovisionDatabase(id string, databaseid string) error
}

type databaseService struct {
	database  database.Repository
	container container.Repository
}

var ErrInvalidArgument = errors.New("invalid argument")

func (s *databaseService) ProvisionNewDatabase(id string, specification database.DatabaseSpecification) (*database.Database, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}

	var dbContainer *container.DatabaseContainer
	for _, c := range s.container.FindAll() {
		if c.Image == specification.Type && c.DatabaseServices < 5 {
			dbContainer = c
			break
		}
	}

	// There are no availables containers in order to provision new databases
	if dbContainer == nil {
		var err error
		dbContainer, err = container.NewDatabaseContainer(container.ContainerSpecification{Image: specification.Type})
		if err != nil {
			return nil, err
		}

		if err = dbContainer.StartContainer(); err != nil {
			return nil, err
		}
	}
	database, err := dbContainer.CreateNewDatabase()
	if err != nil {
		return nil, err
	}
	database.UserID = id

	s.container.Store(dbContainer)
	s.database.Store(database)
	return database, err
}

func (s *databaseService) ProvisionedDatabases(id string) ([]*database.Database, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}

	databases := make([]*database.Database, 0)
	for _, d := range s.database.FindAll() {
		if d.UserID == id {
			databases = append(databases, d)
		}
	}

	return databases, nil
}

func (s *databaseService) UnprovisionDatabase(userID string, databaseID string) error {
	if userID == "" || databaseID == "" {
		return ErrInvalidArgument
	}

	database, err := s.database.Find(databaseID)
	if err != nil {
		return err
	}

	if database.UserID != userID {
		return errors.New("No database found for user " + string(userID))
	}

	container, err := s.container.Find(database.ContainerID)
	if err != nil {
		return errors.New("Error trying to delete database for user " + string(userID))
	}

	if err := container.DeleteDatabase(database); err != nil {
		return err
	}
	s.database.Delete(database.ID)

	return nil
}
