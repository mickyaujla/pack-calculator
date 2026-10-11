package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/mickyaujla/pack-calculator/server/internal/packing"
)

func TestCalculateEndpoint(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"items": 12001}`))
	recorder := httptest.NewRecorder()
	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected 200", recorder.Code)
	}

	var response CalculateResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	expectedPacks := []packing.Pack{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}}

	if response.Items != 12001 || response.TotalItems != 12250 || response.TotalPacks != 4 {
		t.Errorf("got %+v", response)
	}

	if !slices.Equal(response.Packs, expectedPacks) {
		t.Errorf("packs = %v, expected %v", response.Packs, expectedPacks)
	}
}

func TestCalculateEndpointWithCustomPackSizes(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	body := `{"items": 500000, "packSizes": [23, 31, 53]}`
	request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected 200", recorder.Code)
	}

	var response CalculateResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	expectedPacks := []packing.Pack{{Size: 53, Count: 9429}, {Size: 31, Count: 7}, {Size: 23, Count: 2}}
	if !slices.Equal(response.Packs, expectedPacks) {
		t.Errorf("packs = %v, expected %v", response.Packs, expectedPacks)
	}
}

func TestCalculateEndpointWithInvalidPackSizes(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"items": 10, "packSizes": [0, -5]}`))
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, expected 400", recorder.Code)
	}
}

func TestCalculateEndpointBadRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"not JSON", "not json"},
		{"empty body", ""},
		{"negative items", `{"items": -5}`},
		{"zero items", `{"items": 0}`},
		{"items is a string", `{"items": "lots"}`},
	}

	server := NewServer(packing.DefaultSizes())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(test.body))

			recorder := httptest.NewRecorder()

			server.Routes().ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Errorf("status = %d, expected 400", recorder.Code)
			}
		})
	}
}

func TestWrongMethod(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	request := httptest.NewRequest("GET", "/api/calculate", nil)
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, expected 405", recorder.Code)
	}
}

func TestWriteJSONEncodeError(t *testing.T) {
	recorder := httptest.NewRecorder()

	writeJSON(recorder, http.StatusOK, make(chan int))

	if recorder.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", recorder.Body.String())
	}
}

func TestGetPackSizes(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	request := httptest.NewRequest("GET", "/api/pack-sizes", nil)
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, request)

	var response map[string][]int

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if recorder.Code != http.StatusOK || !slices.Equal(response["packSizes"], packing.DefaultSizes()) {
		t.Errorf("got %d %v", recorder.Code, response)
	}
}

func TestHealth(t *testing.T) {
	server := NewServer(packing.DefaultSizes())
	request := httptest.NewRequest("GET", "/healthz", nil)
	recorder := httptest.NewRecorder()

	server.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || recorder.Body.String() != "ok" {
		t.Errorf("got %d %q, expected 200 \"ok\"", recorder.Code, recorder.Body.String())
	}
}
