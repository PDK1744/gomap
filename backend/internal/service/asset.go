package service

import (
	"context"

	"github.com/PDK1744/gomap/internal/store"
)

type AssetService struct {
	assetStore *store.AssetStore
}

func NewAssetService(store *store.AssetStore) *AssetService {
	return &AssetService{assetStore: store}
}

// TODO: Define a BranchData struct that contains Assets and Assignments
type BranchData struct {
	Assets      map[string]AssetInfo `json:"assets"`
	Assignments map[string][]int     `json:"assignments`
}

type AssetInfo struct {
	//ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	BranchID string `json:"branch_id"`
}

func (a *AssetService) GetAssetsByBranch(ctx context.Context, branchName string) ([]Asset, error) {
	branchID, err := a.assetStore.FetchBranchIdByName(ctx, branchName)
	if err != nil {
		return nil, err
	}

}
