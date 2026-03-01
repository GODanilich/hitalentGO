package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterPathParams(t *testing.T) {
	r := New()
	called := false
	r.Handle(http.MethodGet, "/departments/{id}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		called = true
		if Param(req, "id") != "42" {
			t.Fatalf("expected id param=42, got %q", Param(req, "id"))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/departments/42", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if !called {
		t.Fatal("handler was not called")
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestSplitPathTrimsSlashes(t *testing.T) {
	parts := splitPath("/departments/42/")
	if len(parts) != 2 || parts[0] != "departments" || parts[1] != "42" {
		t.Fatalf("unexpected parts: %#v", parts)
	}
}
