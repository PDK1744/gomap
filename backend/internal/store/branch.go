package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PDK1744/gomap/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchStore struct {
	db *pgxpool.Pool
}

func NewBranchStore(db *pgxpool.Pool) *BranchStore {
	return &BranchStore{db: db}
}

const BASE_DIR = "./branches_local"

type MetaMap struct {
	Map map[string]models.MetadataItem
}

func (bs *BranchStore) FetchMetadataBatch(ctx context.Context, truIDs []string) (*MetaMap, error) {
	rows, err := bs.db.Query(ctx,
		`SELECT tru_id, name, object_type, assignable, is_generator, parent FROM object_metadata WHERE tru_id = ANY($1)`,
		truIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metaMap := make(map[string]models.MetadataItem, len(truIDs))
	for rows.Next() {
		var id string
		var item models.MetadataItem
		rows.Scan(&id, &item.Name, &item.Type, &item.Assignable, &item.Generator, &item.Parent)
		metaMap[id] = item
	}
	return &MetaMap{Map: metaMap}, nil
}

type Metadata struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Assignable bool   `json:"assignable"`
	Parent     string `json:"parent"`
	Generator  bool   `json:"generator"`
}

func (bs *BranchStore) SeedMetadata(ctx context.Context, branchName string, branchTruID string) error {
	jsonPath := filepath.Join(BASE_DIR, branchName, "metadata.json")

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	var data map[string]Metadata
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var batch pgx.Batch

	for truID, meta := range data {
		batch.Queue(`
			INSERT INTO object_metadata (
				tru_id,
				branch_tru_id,
				name,
				object_type,
				assignable,
				is_generator,
				parent
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT (tru_id)
			DO UPDATE SET
				branch_tru_id = EXCLUDED.branch_tru_id,
				name = EXCLUDED.name,
				object_type = EXCLUDED.object_type,
				assignable = EXCLUDED.assignable,
				is_generator = EXCLUDED.is_generator,
				parent = EXCLUDED.parent
		`,
			truID,
			branchTruID,
			meta.Name,
			meta.Type,
			meta.Assignable,
			meta.Generator,
			meta.Parent,
		)
	}
	br := bs.db.SendBatch(ctx, &batch)
	defer br.Close()

	for range len(data) {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("batch execution failed: %w", err)
		}
	}
	return nil

}
