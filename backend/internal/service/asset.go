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
	Assets      map[int]AssetInfo `json:"assets"`
	Assignments map[string][]int  `json:"assignments"`
}

type AssetInfo struct {
	//ID       string    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	BranchID string `json:"branch_id"`
}

func (a *AssetService) FetchBranchAssetsData(ctx context.Context, branchName string) (*BranchData, error) {
	branchID, err := a.assetStore.FetchBranchIdByName(ctx, branchName)
	if err != nil {
		return nil, err
	}

	assets, err := a.fetchAssetsByBranch(ctx, branchID)
	if err != nil {
		return nil, err
	}

	assignmentData, err := a.fetchAssetAssignments(ctx, branchID)
	if err != nil {
		return nil, err
	}

	return &BranchData{
		Assets:      assets,
		Assignments: assignmentData,
	}, nil

}

func (a *AssetService) fetchAssetsByBranch(ctx context.Context, branchID string) (map[int]AssetInfo, error) {
	assetList, err := a.assetStore.GetAssetsByBranch(ctx, branchID)
	if err != nil {
		return nil, err
	}
	assets := make(map[int]AssetInfo, len(assetList))
	for _, a := range assetList {
		assets[a.ID] = AssetInfo{
			Name:     a.Name,
			Type:     a.Type,
			BranchID: a.BranchID,
		}
	}
	return assets, nil
}

func (a *AssetService) fetchAssetAssignments(ctx context.Context, branchID string) (map[string][]int, error) {
	assetAssignments, err := a.assetStore.GetRoomAssignments(ctx, branchID)
	if err != nil {
		return nil, err
	}
	assignmentData := make(map[string][]int, len(assetAssignments))
	for _, a := range assetAssignments {
		assignmentData[a.RoomID] = append(assignmentData[a.RoomID], a.AssetID)
	}

	return assignmentData, nil
}

// func (a *AssetService) GetAssetsByBranch(ctx context.Context, branchName string) ([]Asset, error) {
// 	// branchID, err := a.assetStore.FetchBranchIdByName(ctx, branchName)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	return nil, nil

// }
