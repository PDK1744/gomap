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
	Branch   string `json:"branch"`
	BranchID string `json:"branch_id"`
	Status   bool   `json:"status"`
}

func (a *AssetStore) GetAssetsByBranch(ctx context.Context, branchName string) ([]Asset, error) {
	branchIDQuery := "SELECT branch_id FROM branch_alias WHERE name ILIKE $1"
	var branchID string
	err := a.db.QueryRow(ctx, branchIDQuery, branchName).Scan(&branchID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("No branch found with that ID: %v", err)
		}
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

	return assets, nil

}
