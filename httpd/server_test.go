package httpd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: this test is sort of pointless and flaky
// the idea is to test if cancelling the parent context stops the Run function
// and this graceful shutdown isnt treated as an error by the application (because it isnt)
func TestGracefulShutdownHttpDaemon(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := Run(ctx, []string{}, func(string) string { return "" }, os.Stdin, io.Discard, io.Discard)
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

// test struct to check json decoding
type obj struct {
	Name string `json:"name"`
}

func TestDecodeSingle(t *testing.T) {
	j := `{"name": "Bob"}`
	req, err := http.NewRequest(http.MethodPost, "/service/method", strings.NewReader(j))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	requestObject, err := decode[obj](req)
	assert.NoError(t, err)
	assert.Equal(t, requestObject.Name, "Bob")
}

func TestDecodeArray(t *testing.T) {
	j := `[
		{"name": "Bob"},
		{"name": "Bobbe"},
		{"name": "Foo"}
	]`
	req, err := http.NewRequest(http.MethodPost, "/service/method", strings.NewReader(j))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	requestObjects, err := decode[[]obj](req)
	assert.NoError(t, err)
	assert.Equal(t, len(requestObjects), 3)
	assert.Equal(t, requestObjects[0].Name, "Bob")
	assert.Equal(t, requestObjects[1].Name, "Bobbe")
	assert.Equal(t, requestObjects[2].Name, "Foo")
}
