package httpd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"plants/plants"
	"plants/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPlants(t *testing.T) {
	t.Skip()
	testPlants := []plants.Plant{
		{ID: "1", Name: "foo", Height: 4},
		{ID: "2", Name: "bar", Height: 3},
		{ID: "3", Name: "baz", Height: 2},
	}

	testError := errors.New("foo bar test error")

	tests := map[string]struct {
		store store.Store

		wantResponse string
		wantCode     int
	}{
		"returns empty json array when no data": {
			store: &mockStore{plants: []plants.Plant{}},

			wantResponse: `[]`,
			wantCode:     http.StatusOK,
		},
		"returns json array": {
			store: &mockStore{plants: testPlants},

			wantResponse: `[{"id":"1","name":"foo","height":4},{"id":"2","name":"bar","height":3},{"id":"3","name":"baz","height":2}]`,
			wantCode:     http.StatusOK,
		},
		"returns error when error": {
			store: &mockStore{err: testError},

			wantResponse: `{"message":"retrieve all plants: foo bar test error"}`,
			wantCode:     http.StatusInternalServerError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// nil logger is noop
			// server := NewServer(nil, tc.store)

			// request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// service.handleListPlants(w, request)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantCode {
				t.Errorf("status code mismatch, expected: %v, got: %v", tc.wantCode, res.StatusCode)
			}

			gotBody, gotErr := io.ReadAll(res.Body)
			if gotErr != nil {
				t.Errorf("failed to read response body: %v", gotErr)
			}

			assert.JSONEq(t, tc.wantResponse, string(gotBody))
		})
	}
}

func TestGetPlantByID(t *testing.T) {
	t.Skip()
	testPlant := plants.Plant{ID: "2", Name: "bar", Height: 3}
	testError := errors.New("foo bar test error")

	tests := map[string]struct {
		store store.Store
		id    string

		wantResponse string
		wantCode     int
	}{
		"returns error when no data": {
			store: &mockStore{},
			id:    "123",

			wantResponse: `{"message":"plant with ID '123' does not exist"}`,
			wantCode:     http.StatusNotFound,
		},
		"returns object json": {
			store: &mockStore{plant: &testPlant},
			id:    "123",

			wantResponse: `{"id":"2","name":"bar","height":3}`,
			wantCode:     http.StatusOK,
		},
		"returns error when store error": {
			store: &mockStore{err: testError},
			id:    "123",

			wantResponse: `{"message":"find plant by id: foo bar test error"}`,
			wantCode:     http.StatusInternalServerError,
		},

		"returns error when invalid request params": {
			store: &mockStore{plant: &testPlant},
			id:    "",

			wantResponse: `{"message":"id is required in path parameters"}`,
			wantCode:     http.StatusUnprocessableEntity,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// service := newHttpService(tc.store, nil)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.id != "" {
				request.SetPathValue("id", tc.id)
			}
			w := httptest.NewRecorder()

			// service.handleGetPlant(w, request)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantCode {
				t.Errorf("status code mismatch, expected: %v, got: %v", tc.wantCode, res.StatusCode)
			}

			gotBody, gotErr := io.ReadAll(res.Body)
			if gotErr != nil {
				t.Errorf("failed to read response body: %v", gotErr)
			}

			assert.JSONEq(t, tc.wantResponse, string(gotBody))
		})
	}
}

// TODO: write tests for create handler

func TestCreatePlant(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	go Run(ctx, []string{}, os.Getenv, os.Stdin, os.Stdout, os.Stderr)

	plant := plants.Plant{ID: "foo", Name: "foo", Height: 5}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(plant)
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/plants", &buf)
	assert.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		assert.Fail(t, "recieved non 200 response")
	}

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

}

type mockStore struct {
	plants []plants.Plant
	plant  *plants.Plant
	err    error
}

func (s *mockStore) List() ([]plants.Plant, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.plants, nil
}

func (s *mockStore) Find(id string) (*plants.Plant, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	if s.err != nil {
		return nil, s.err
	}
	return s.plant, nil
}

func (s *mockStore) Create(plant plants.Plant) (*plants.Plant, error) {
	plant.ID = "new id"
	return &plant, nil
}
