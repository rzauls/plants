package httpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"plants/plants"
	"plants/store"
	"testing"
)

func TestListPlants(t *testing.T) {
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
			service := newHttpService(tc.store, nil)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			service.handleListPlants(w, request)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantCode {
				t.Errorf("status code mismatch, expected: %v, got: %v", tc.wantCode, res.StatusCode)
			}

			gotBody, gotErr := io.ReadAll(res.Body)
			if gotErr != nil {
				t.Errorf("failed to read response body: %v", gotErr)
			}
			if string(gotBody) != tc.wantResponse {
				t.Errorf(gotWantJson(string(gotBody), tc.wantResponse))
			}
		})
	}
}

func TestGetPlantByID(t *testing.T) {
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
			// nil logger is noop
			service := newHttpService(tc.store, nil)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.id != "" {
				request.SetPathValue("id", tc.id)
			}
			w := httptest.NewRecorder()

			service.handleGetPlant(w, request)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantCode {
				t.Errorf("status code mismatch, expected: %v, got: %v", tc.wantCode, res.StatusCode)
			}

			gotBody, gotErr := io.ReadAll(res.Body)
			if gotErr != nil {
				t.Errorf("failed to read response body: %v", gotErr)
			}
			if string(gotBody) != tc.wantResponse {
				t.Errorf(gotWantJson(string(gotBody), tc.wantResponse))
			}
		})
	}
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

// gotWantJson - formats json strings into somewhat more human-comparable look
func gotWantJson(got any, want any) string {
	spacing := "  "
	gotJson, _ := json.MarshalIndent(got, "", spacing)
	wantJson, _ := json.MarshalIndent(want, "", spacing)

	return fmt.Sprintf("\ngot:\n%s\nwant:\n%s", string(gotJson), string(wantJson))
}
