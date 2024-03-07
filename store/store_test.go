package store

import (
	"plants/plants"
	"reflect"
	"testing"
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
				// there is nothing to test here, since MemotyStore cant fail the List method
			}

			if len(tc.want) > 0 || len(got) > 0 {
				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("expected: %v, got: %v", tc.want, got)
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
				t.Fatalf("got error when didnt expect one: %v", gotErr)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
