package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PDK1744/gomap/internal/service"
)

type BranchHandler struct {
	branchService *service.BranchService
	assetService  *service.AssetService
}

func NewHandler(brService *service.BranchService, asService *service.AssetService) *BranchHandler {
	return &BranchHandler{branchService: brService, assetService: asService}
}

func (h *BranchHandler) GetBranchLayout(w http.ResponseWriter, r *http.Request) {
	// TODO: Centralize Logs
	fmt.Println("GET BRANCH Request Received")
	branchName := r.PathValue("branch_name")
	if branchName == "" {
		http.Error(w, "branch name is empty", http.StatusBadRequest)
		return
	}
	layout, err := h.branchService.FetchBranchLayout(r.Context(), branchName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Branch map data not found", http.StatusNotFound)
			return
		}
		log.Printf("Error with request: %v", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(layout); err != nil {
		log.Fatalf("ENCODING FAILURE: %v", err)
		return
	}
}

func (h *BranchHandler) GetAssetsByBranch(w http.ResponseWriter, r *http.Request) {
	// TODO: Centralize Logs
	fmt.Println("GET ASSETS BY BRANCH Request Received")
	branchName := r.PathValue("branch_name")
	assets, err := h.assetStore.GetAssetsByBranch(r.Context(), branchName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "assets not found", http.StatusNotFound)
			return
		}
		log.Printf("Error with request: %v", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(assets); err != nil {
		log.Fatalf("ENCODING FAILURE: %v", err)
		return
	}
}

func (h *BranchHandler) GetAssetRoomAssignments(w http.ResponseWriter, r *http.Request) {
	// TODO: Centralize Logs
	branchName := r.PathValue("branch_name")
	assign, err := h.assetStore.GetRoomAssignments(r.Context(), branchName)
	if err != nil {
		log.Printf("[REQUEST ERROR]: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	assingments := make(map[string][]string, len(assign))
	for _, a := range assign {
		assingments[a.RoomID] = append(assingments[a.RoomID], a.AssetID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(assingments); err != nil {
		log.Fatalf("ENCODING FAILURE: %v", err)
		return
	}

}

func (h *BranchHandler) PutAssetRoomAssignments(w http.ResponseWriter, r *http.Request) {
	var assignments map[string][]string

	err := json.NewDecoder(r.Body).Decode(&assignments)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// TODO: call PutRoomAssignments

}
