package httpd

import (
	"net/http"
	"testing"
)

func TestRoutePathGenerator(t *testing.T) {
	testCases := []struct {
		name     string
		rpg      routePathGenerator
		method   string
		path     string
		wantPath string
	}{
		{
			name:     "valid path/method combo",
			rpg:      routePathGenerator{root: "/test"},
			method:   http.MethodGet,
			path:     "/foo",
			wantPath: "GET /test/foo",
		},
		{
			name:     "valid path/method combo with path variable",
			rpg:      routePathGenerator{root: "/bar"},
			method:   http.MethodPost,
			path:     "/foo/{id}",
			wantPath: "POST /bar/foo/{id}",
		},
		{
			name:     "valid path/method combo when no root",
			rpg:      routePathGenerator{},
			method:   http.MethodPost,
			path:     "/foo/{id}",
			wantPath: "POST /foo/{id}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.rpg.route(tc.method, tc.path)
			if got != tc.wantPath {
				t.Errorf("generate path\ngot: %s, want %s", got, tc.wantPath)
			}
		})
	}
}
