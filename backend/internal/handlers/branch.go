package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/PDK1744/gomap/internal/store"
)

type Handler struct {
	branchStore *store.BranchStore
	assetStore  *store.AssetStore
}

func NewBranchHandler(branchStore *store.BranchStore, assetStore *store.AssetStore) *Handler {
	return &Handler{branchStore: branchStore, assetStore: assetStore}
}

func (h *Handler) GetBranchLayout(w http.ResponseWriter, r *http.Request) {
	branchName := r.PathValue("branch_name")
	layout, err := h.branchStore.FetchBranchLayout(r.Context(), branchName)
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

func (h *Handler) GetAssetsByBranch(w http.ResponseWriter, r *http.Request) {
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
