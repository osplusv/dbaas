package inmem

import (
	"sync"

	"github.com/osplusv/dbaas/container"
	"github.com/osplusv/dbaas/database"
)

type containerRepository struct {
	mtx        sync.RWMutex
	containers map[string]*container.DatabaseContainer
}

func (r *containerRepository) Store(c *container.DatabaseContainer) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.containers[c.ContainerID] = c
	return nil
}

func (r *containerRepository) Find(id string) (*container.DatabaseContainer, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if val, ok := r.containers[id]; ok {
		return val, nil
	}
	return nil, container.ErrUnknown
}

func (r *containerRepository) FindAll() []*container.DatabaseContainer {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	c := make([]*container.DatabaseContainer, 0, len(r.containers))
	for _, val := range r.containers {
		c = append(c, val)
	}
	return c
}

// NewContainerRepository returns a new instance of a in-memory container repository.
func NewContainerRepository() container.Repository {
	return &containerRepository{
		containers: make(map[string]*container.DatabaseContainer),
	}
}

type databaseRepository struct {
	mtx       sync.RWMutex
	databases map[string]*database.Database
}

func (r *databaseRepository) Store(d *database.Database) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.databases[d.ID] = d
	return nil
}

func (r *databaseRepository) Find(id string) (*database.Database, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if val, ok := r.databases[id]; ok {
		return val, nil
	}
	return nil, database.ErrUnknown
}

func (r *databaseRepository) FindAll() []*database.Database {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	c := make([]*database.Database, 0, len(r.databases))
	for _, val := range r.databases {
		c = append(c, val)
	}
	return c
}

func (r *databaseRepository) Delete(id string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.databases, id)
	return nil
}

// NewDatabaseRepository returns a new instance of a in-memory database repository.
func NewDatabaseRepository() database.Repository {
	return &databaseRepository{
		databases: make(map[string]*database.Database),
	}
}
