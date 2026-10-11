package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mickyaujla/pack-calculator/server/internal/packing"
)

type Server struct {
	packSizes []int
}

func NewServer(packSizes []int) *Server {
	return &Server{packSizes: packSizes}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /api/pack-sizes", s.getPackSizes)
	mux.HandleFunc("POST /api/calculate", s.calculate)
	return mux
}

type CalculateRequest struct {
	Items     int   `json:"items"`
	PackSizes []int `json:"packSizes"`
}

type CalculateResponse struct {
	Items      int            `json:"items"`
	TotalItems int            `json:"totalItems"`
	TotalPacks int            `json:"totalPacks"`
	Packs      []packing.Pack `json:"packs"`
}

func (s *Server) getPackSizes(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string][]int{"packSizes": s.packSizes})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (s *Server) calculate(w http.ResponseWriter, r *http.Request) {
	var request CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return
	}

	packSizes := s.packSizes
	if len(request.PackSizes) > 0 {
		packSizes = request.PackSizes
	}

	packs, err := packing.Calculate(request.Items, packSizes)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := CalculateResponse{Items: request.Items, Packs: packs}

	for _, pack := range packs {
		response.TotalItems += pack.Size * pack.Count
		response.TotalPacks += pack.Count
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
