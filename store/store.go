package store

import (
	"context"
	"fmt"
	"log/slog"
	"plants/log"
	"plants/plants"

	"github.com/google/uuid"
)

type Store interface {
	Find(ctx context.Context, id string) (*plants.Plant, error)
	List(ctx context.Context) ([]plants.Plant, error)
	Create(ctx context.Context, plant plants.Plant) (*plants.Plant, error)
}

func NewMemoryStore(items []plants.Plant) *MemoryStore {
	return &MemoryStore{
		items: items,
	}
}

type MemoryStore struct {
	items []plants.Plant
}

type ErrorResourceDoesNotExist struct {
	Err error
}

func (e ErrorResourceDoesNotExist) Error() string {
	return e.Err.Error()
}

func (s *MemoryStore) Find(ctx context.Context, id string) (*plants.Plant, error) {
	// fancy slices index generic function
	// index := slices.IndexFunc(s.items, func(p plants.Plant) bool { return p.ID == id })

	logger := log.LoggerFromCtx(ctx, slog.Default())
	logger.Debug("some kind of debug message from store package", slog.Int("additionalField", 42))
	for _, p := range s.items {
		if p.ID == id {
			return &p, nil

		}
	}

	// NOTE: realistically there would be more than 1 way of this find failing, so we could return typed errors and handle them in different ways
	return nil, ErrorResourceDoesNotExist{Err: fmt.Errorf("plant with ID '%s' does not exist", id)}
}

func (s *MemoryStore) List(ctx context.Context) ([]plants.Plant, error) {
	return s.items, nil
}

func (s *MemoryStore) Create(ctx context.Context, plant plants.Plant) (*plants.Plant, error) {
	plant.ID = uuid.New().String()
	s.items = append(s.items, plant)
	return &plant, nil
}
