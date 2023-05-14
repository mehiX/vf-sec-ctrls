package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

func TestAuthenticatorDenyAccess(t *testing.T) {

	s := New()

	genericResp := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	srvr := httptest.NewTLSServer(jwtauth.Verifier(s.tokenAuth)(Authenticator(genericResp)))
	defer srvr.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srvr.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// disable redirects since we want to check the redirect headers are returned correctly
	srvr.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := srvr.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	status := resp.StatusCode
	if status != http.StatusFound {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusFound, status)
	}
	loc := resp.Header.Get("Location")
	if loc != "/login" {
		t.Fatalf("wrong header Location. expected: %s, got: %s", "/login", loc)
	}
}

func TestAuthenticatorAllowAccess(t *testing.T) {

	s := New()
	claims := map[string]any{
		"user_id": "test user",
	}
	jwtauth.SetExpiryIn(claims, time.Hour)

	_, tknStr, err := s.tokenAuth.Encode(claims)
	if err != nil {
		t.Fatal(err)
	}

	handleSecure := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	srvr := httptest.NewTLSServer(jwtauth.Verifier(s.tokenAuth)(Authenticator(handleSecure)))
	defer srvr.Close()

	srvr.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srvr.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tknStr)

	resp, err := srvr.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}
}
