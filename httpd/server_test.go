package httpd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunHttpDaemon(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	go Run(ctx, []string{}, os.Getenv, os.Stdin, io.Discard, io.Discard)

	err := waitForReady(ctx, 10*time.Second, "/health")
	assert.NoError(t, err)
}

func TestEncode(t *testing.T) {
	data := struct {
		Greeting string `json:"greeting"`
	}{
		Greeting: "hello",
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	err := encode(w, r, http.StatusOK, data)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)
	body, _ := io.ReadAll(w.Body)
	if !assert.JSONEq(t, string(body), `{"greeting":"hello"}`) {
		t.Errorf("json doesnt match")
	}
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
}

func TestDecode(t *testing.T) {
	type obj struct {
		Name string
	}
	j := `[
		{"name": "Mat"},
		{"name": "David"},
		{"name": "Aaron"}
	]`
	req, err := http.NewRequest(http.MethodPost, "/service/method", strings.NewReader(j))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	requestObjects, err := decode[[]obj](req)
	assert.NoError(t, err)
	assert.Equal(t, len(requestObjects), 3)
	assert.Equal(t, requestObjects[0].Name, "Mat")
	assert.Equal(t, requestObjects[1].Name, "David")
	assert.Equal(t, requestObjects[2].Name, "Aaron")
}

func waitForReady(ctx context.Context, timeout time.Duration, path string) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			fmt.Printf("error while waiting for ready: %s\n", err.Error())
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error while making request: %s\n", err.Error())
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("endpoint is ready")
			resp.Body.Close()
			return nil
		}

		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait between polling again
			time.Sleep(250 * time.Millisecond)
		}
	}

}
