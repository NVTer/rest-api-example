package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestIDMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	logger, hook := test.NewNullLogger()
	entry := logger.WithField("entry", "exists")

	testRequest(t, req, IDMiddleware(entry))

	assert.NotNil(t, hook.LastEntry().Data["correlation_id"])
}

func TestTimeLogMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	logger, hook := test.NewNullLogger()
	entry := logger.WithField("entry", "exists")

	testRequest(t, req, TimeLogMiddleware(entry))
	assert.Equal(t, 1, len(hook.Entries))
	assert.NotNil(t, hook.LastEntry().Data["work_time"])
}

func TestAccessLogMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	logger, hook := test.NewNullLogger()
	entry := logger.WithField("entry", "exists")
	testRequest(t, req, AccessLogMiddleware(entry))

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, http.MethodGet, hook.LastEntry().Data["method"])
	assert.Equal(t, req.Host, hook.LastEntry().Data["host"])
}

func testRequest(t *testing.T, req *http.Request, middleware func(next http.Handler) http.Handler) {
	w := httptest.NewRecorder()
	r := mux.NewRouter()
	r.Use(middleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			t.Error(err)
		}
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK , w.Code)
}
