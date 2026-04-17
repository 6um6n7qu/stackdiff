package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendSlack_EmptyURL(t *testing.T) {
	if err := sendSlack("", "hello"); err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestSendSlack_Success(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if err := sendSlack(ts.URL, "test message"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text != "test message" {
		t.Errorf("expected %q, got %q", "test message", received.Text)
	}
}

func TestSendSlack_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	if err := sendSlack(ts.URL, "msg"); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
