package store

import (
	"fmt"
	"plants/plants"
)

func NewMemoryStore(items []plants.Plant) *MemoryStore {
	return &MemoryStore{
		items: items,
	}
}

type MemoryStore struct {
	items []plants.Plant
}

func (s *MemoryStore) Find(id string) (*plants.Plant, error) {
	// fancy slices index generic function
	// index := slices.IndexFunc(s.items, func(p plants.Plant) bool { return p.ID == id })

	for _, p := range s.items {
		if p.ID == id {
			return &p, nil

		}
	}

	// NOTE: realistically there would be more than 1 way of this find failing, so we could return typed errors and handle them in different ways
	return nil, fmt.Errorf("plant with ID '%s' does not exist", id)
}

func (s *MemoryStore) List() ([]plants.Plant, error) {
	return s.items, nil
}

type Store interface {
	Find(id string) (*plants.Plant, error)
	List() ([]plants.Plant, error)
}
