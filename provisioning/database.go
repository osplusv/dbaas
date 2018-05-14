package provisioning

import "github.com/osplusv/dbaas/database"

type DatabaseProvisioner interface {
	ProvisionNewDatabase(id string, specification database.DatabaseSpecification) (*database.Database, error)
	ProvisionedDatabases(id string) ([]*database.Database, error)
	UnprovisionDatabase(id string, databaseid string) error
}
