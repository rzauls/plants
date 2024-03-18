package store

import (
	"plants/plants"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStoreList(t *testing.T) {
	testPlants := []plants.Plant{
		{ID: "1", Name: "foo", Height: 4},
		{ID: "2", Name: "bar", Height: 3},
		{ID: "3", Name: "baz", Height: 2},
	}

	// NOTE: in memory implementation of .List cant return an error, so there is nothing to test here
	// NOTE: this test is quite arbitrary since the implementation is extremely basic
	tests := map[string]struct {
		store Store

		want    []plants.Plant
		wantErr *error
	}{
		"empty store returns empty slice": {
			store: &MemoryStore{},

			want:    []plants.Plant{},
			wantErr: nil,
		},
		"store returns its items": {
			store: &MemoryStore{items: testPlants},

			want:    testPlants,
			wantErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := tc.store.List()
			if gotErr != nil && tc.wantErr != nil {
				t.Fatalf("unexpected branch, no errors should occur in this test: %v", gotErr)
			}

			if len(tc.want) > 0 || len(got) > 0 {
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("expected: %v, got: %v", tc.want, got)
				}
			}
		})
	}
}

func TestMemoryStoreFind(t *testing.T) {
	testPlant := plants.Plant{ID: "2", Name: "bar", Height: 3}
	testPlants := []plants.Plant{
		{ID: "1", Name: "foo", Height: 4},
		testPlant,
		{ID: "3", Name: "baz", Height: 2},
	}

	tests := map[string]struct {
		store Store
		id    string

		want *plants.Plant
		// we only test if there was an error, since the error types dont matter for this implementation
		wantErr bool
	}{
		"store finds item by ID": {
			store: &MemoryStore{items: testPlants},
			id:    "2",

			want:    &testPlant,
			wantErr: false,
		},
		"store returns error if item not found": {
			store: &MemoryStore{},

			want:    nil,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := tc.store.Find(tc.id)
			if gotErr != nil && !tc.wantErr {
				t.Errorf("got error when didnt expect one: %v", gotErr)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestMemoryStoreCreate(t *testing.T) {
	testPlant := plants.Plant{Name: "foo", Height: 1}

	tests := map[string]struct {
		store Store
		plant plants.Plant

		want *plants.Plant
	}{
		"creates an item": {
			store: &MemoryStore{},
			plant: testPlant,

			want: &testPlant,
		},
		// NOTE: memoryStore.Create cannot fail so there is no error cases to test
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.store.Create(tc.plant)
			if err != nil {
				t.Errorf("got error when didnt expect one: %v", err)
			}

			gotPlant, err := tc.store.Find(got.ID)
			if err != nil {
				t.Errorf("got error when didnt expect one while reading from store: %v", err)
			}

			assert.Equal(t, tc.want.Name, gotPlant.Name)
			assert.Equal(t, tc.want.Name, gotPlant.Name)
			// check if the ID is initialized after returning from store
			assert.NotZero(t, gotPlant.ID)
		})
	}
}
