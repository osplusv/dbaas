package inmem

import (
	"sync"

	"github.com/osplusv/dbaas/container"
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
