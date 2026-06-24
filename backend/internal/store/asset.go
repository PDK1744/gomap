package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetStore struct {
	db *pgxpool.Pool
}

func NewAssetStore(db *pgxpool.Pool) *AssetStore {
	return &AssetStore{
		db: db,
	}
}

type Asset struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	BranchID string `json:"branch_id"`
}

type ObjectAsset struct {
	RoomID  string //`json:"room_id"`
	AssetID string //`json:"asset_id"`
	//BranchID string //`json:"branch_id,omitempty"`
}

func (a *AssetStore) GetAssetsByBranch(ctx context.Context, branchName string) ([]Asset, error) {

	branchID, err := a.fetchBranchIdByName(ctx, branchName)
	if err != nil {
		return nil, err
	}

	var assets []Asset
	assetQuery := "SELECT asset_id, name, asset_type, branch_id from assets WHERE is_active = true AND branch_id = $1"

	rows, err := a.db.Query(ctx, assetQuery, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.BranchID); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, rows.Err()

}

func (a *AssetStore) GetRoomAssignments(ctx context.Context, branchName string) ([]ObjectAsset, error) {
	branchID, err := a.fetchBranchIdByName(ctx, branchName)
	if err != nil {
		return nil, err
	}
	if branchID == "" {
		return nil, fmt.Errorf("branch_id cannot be blank")
	}

	var objectAsset []ObjectAsset

	query := "SELECT room_id, asset_id from object_assets WHERE branch_id = $1"

	rows, err := a.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oa ObjectAsset
		if err := rows.Scan(&oa.RoomID, &oa.AssetID); err != nil {
			return nil, err
		}
		objectAsset = append(objectAsset, oa)
	}
	return objectAsset, rows.Err()

}

func (a *AssetStore) fetchBranchIdByName(ctx context.Context, branchName string) (string, error) {
	branchIDQuery := "SELECT branch_id FROM branch_aliases WHERE alias ILIKE $1"
	var branchID string
	err := a.db.QueryRow(ctx, branchIDQuery, branchName).Scan(&branchID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("No branch found with that ID: %v", err)
		}
		return "", err
	}
	return branchID, nil
}
