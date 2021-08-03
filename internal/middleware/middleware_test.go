package middleware

import (
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

	if len(hook.Entries) != 1 {
		t.Fatalf("error: len = 1")
	}
	if hook.LastEntry().Data["correlation_id"] == nil {
		t.Fatalf("correalation id is nil")
	}
}

func TestTimeLogMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	logger, hook := test.NewNullLogger()
	entry := logger.WithField("entry", "exists")

	testRequest(t, req, TimeLogMiddleware(entry))

	if len(hook.Entries) != 1 {
		t.Fatalf("error: len = 1")
	}
	if hook.LastEntry().Data["work_time"] == nil {
		t.Fatalf("correalation id is nil")
	}
}

func TestAccessLogMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	logger, hook := test.NewNullLogger()
	entry := logger.WithField("entry", "exists")
	testRequest(t, req, AccessLogMiddleware(entry))

	if len(hook.Entries) != 1 {
		t.Fatalf("error: len = 1")
	}
	if http.MethodGet != hook.LastEntry().Data["method"] &&
		req.Host != hook.LastEntry().Data["host"] {
		t.Fatalf("error: expected %v; get %v", "GET", hook.LastEntry().Data["method"])
	}
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
	if http.StatusOK != w.Code {
		t.Fatalf("Expected: %v; get: %v", http.StatusOK, w.Code)
	}
}
