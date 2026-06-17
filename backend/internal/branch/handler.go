package branch

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

type Handler struct {
	BaseDir string
}

func NewHandler(baseDir string) *Handler {
	return &Handler{BaseDir: baseDir}
}

func (h *Handler) GetBranchDetails(w http.ResponseWriter, r *http.Request) {
	branchName := r.PathValue("branch_name")
	if branchName == "" {
		http.Error(w, "Branch name is required", http.StatusBadRequest)
		return
	}

	response, err := LoadAndNormalizeBranch(h.BaseDir, branchName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Branch map data not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatalf("ENCODING FAILURE: %v", err)
		return
	}
}
