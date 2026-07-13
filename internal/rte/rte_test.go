package rte

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// patchBody is the JSON PATCH /gpio/{id} body the client sends.
type patchBody struct {
	State     uint   `json:"state"`
	Direction string `json:"direction"`
	Time      uint   `json:"time"`
}

// captureServer records the last request method, path, and decoded PATCH body,
// and replies with the given status and JSON body.
func captureServer(t *testing.T, status int, reply string) (*httptest.Server, *struct {
	method string
	path   string
	body   patchBody
}) {
	t.Helper()
	last := &struct {
		method string
		path   string
		body   patchBody
	}{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last.method = r.Method
		last.path = r.URL.Path
		if r.Method == http.MethodPatch {
			data, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(data, &last.body)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, reply)
	}))
	t.Cleanup(srv.Close)
	return srv, last
}

// testClient points a Client at a test server.
func testClient(srv *httptest.Server) *Client {
	return &Client{base: srv.URL, hc: srv.Client()}
}

func TestNewBuildsBaseURL(t *testing.T) {
	c := New("bench.local")
	if c.base != "http://bench.local:8000/api/v1" {
		t.Errorf("base = %q, want http://bench.local:8000/api/v1", c.base)
	}
}

func TestNewBracketsIPv6Host(t *testing.T) {
	c := New("fd00::1")
	if c.base != "http://[fd00::1]:8000/api/v1" {
		t.Errorf("base = %q, want the IPv6 literal bracketed", c.base)
	}
}

func TestPinAsserted(t *testing.T) {
	if !(Pin{State: 1}).Asserted() || (Pin{State: 0}).Asserted() {
		t.Error("Asserted should be true for odd State, false for even")
	}
}

func TestGetParsesPin(t *testing.T) {
	srv, last := captureServer(t, 200, `{"id":13,"state":1,"direction":"in","description":"led1"}`)
	pin, err := testClient(srv).Get(13)
	if err != nil {
		t.Fatal(err)
	}
	if last.method != "GET" || last.path != "/gpio/13" {
		t.Errorf("request = %s %s, want GET /gpio/13", last.method, last.path)
	}
	want := Pin{ID: 13, State: 1, Direction: "in", Description: "led1"}
	if pin != want {
		t.Errorf("Get = %+v, want %+v", pin, want)
	}
}

func TestSetOpenCollectorMapsState(t *testing.T) {
	cases := []struct {
		state string
		want  uint
	}{
		{"low", 1},    // OC asserted
		{"high-z", 0}, // OC released
	}
	for _, tc := range cases {
		srv, last := captureServer(t, 200, `{}`)
		if err := testClient(srv).Set(3, tc.state, 0); err != nil {
			t.Fatalf("Set(3, %q): %v", tc.state, err)
		}
		if last.method != "PATCH" || last.path != "/gpio/3" {
			t.Errorf("request = %s %s, want PATCH /gpio/3", last.method, last.path)
		}
		if last.body.State != tc.want || last.body.Direction != "out" {
			t.Errorf("Set(3, %q) body = %+v, want state %d direction out", tc.state, last.body, tc.want)
		}
	}
}

func TestSetPushPullMapsState(t *testing.T) {
	cases := []struct {
		state string
		want  uint
	}{
		{"high", 1},
		{"low", 0},
	}
	for _, tc := range cases {
		srv, last := captureServer(t, 200, `{}`)
		if err := testClient(srv).Set(14, tc.state, 0); err != nil {
			t.Fatalf("Set(14, %q): %v", tc.state, err)
		}
		if last.body.State != tc.want || last.body.Direction != "out" {
			t.Errorf("Set(14, %q) body = %+v, want state %d direction out", tc.state, last.body, tc.want)
		}
	}
}

func TestSetSendsHoldTime(t *testing.T) {
	srv, last := captureServer(t, 200, `{}`)
	if err := testClient(srv).Set(9, "low", 6); err != nil {
		t.Fatal(err)
	}
	if last.body.Time != 6 {
		t.Errorf("Set hold time = %d, want 6", last.body.Time)
	}
}

func TestSetRejectsStateOutOfRange(t *testing.T) {
	srv, _ := captureServer(t, 200, `{}`)
	// An open-collector pin cannot be driven high.
	if err := testClient(srv).Set(3, "high", 0); err == nil {
		t.Error("Set(3, high) should error: open-collector pins cannot be driven high")
	}
	// A push-pull pin has no high-z state.
	if err := testClient(srv).Set(14, "high-z", 0); err == nil {
		t.Error("Set(14, high-z) should error: push-pull pins have no high-z state")
	}
}

func TestSetErrorsOnNon200(t *testing.T) {
	srv, _ := captureServer(t, 404, `{"error":"id not found"}`)
	err := testClient(srv).Set(3, "low", 0)
	if err == nil {
		t.Fatal("Set should error on a non-200 response")
	}
	if !strings.Contains(err.Error(), "id not found") {
		t.Errorf("error %q should include the server message", err)
	}
}

func TestGetErrorsOnNon200(t *testing.T) {
	srv, _ := captureServer(t, 404, `{"error":"id not found"}`)
	if _, err := testClient(srv).Get(99); err == nil {
		t.Fatal("Get should error on a non-200 response")
	}
}
