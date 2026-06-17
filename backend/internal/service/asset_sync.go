package service

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetSyncService struct {
	mapDB *pgxpool.Pool
	invDB *pgxpool.Pool
}

type Asset struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Branch   string `json:"branch"`
	BranchID string `json:"branch_id"`
	Status   string `json:"status"`
}

type LocalAsset struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Branch   string `json:"branch"`
	BranchID string `json:"branch_id"`
	Status   bool   `json:"status"`
}

type AssetKey struct {
	Type string
	ID   int
}

func NewAssetSyncService(mapDB, invDB *pgxpool.Pool) *AssetSyncService {
	return &AssetSyncService{mapDB: mapDB, invDB: invDB}
}

func (s *AssetSyncService) SyncAssets(ctx context.Context) error {
	// load branch aliases
	branchAliases, err := s.loadBranchAliases(ctx)
	if err != nil {
		return err
	}

	// load external assets
	externalAssets := make(map[AssetKey]Asset)

	pcs, err := s.loadExternalPCs(ctx)
	if err != nil {
		return err
	}
	for _, pc := range pcs {
		key := AssetKey{
			Type: pc.Type,
			ID:   pc.ID,
		}
		id, ok := branchAliases[strings.ToLower(pc.Branch)]
		if !ok {
			// TODO: possibly store the not matched branch names for review???
			log.Printf("WARNING: no branch alias found for %q", pc.Branch)
		} else {
			pc.BranchID = id
		}
		externalAssets[key] = pc
	}

	printers, err := s.loadExternalPrinters(ctx)
	for _, pr := range printers {
		key := AssetKey{
			Type: pr.Type,
			ID:   pr.ID,
		}
		id, ok := branchAliases[strings.ToLower(pr.Branch)]
		if !ok {
			// TODO: possibly store the not matched branch names for review???
			log.Printf("WARNING: no branch alias found for %q", pr.Branch)
		} else {
			pr.BranchID = id
		}
		externalAssets[key] = pr
	}

	// load local assets
	localSlice, err := s.loadLocalAssets(ctx)
	if err != nil {
		return err
	}

	local := make(map[AssetKey]LocalAsset)
	for _, a := range localSlice {
		key := AssetKey{
			Type: a.Type,
			ID:   a.ID,
		}
		local[key] = a
	}

	// build diff lists

	toUpsert := make([]Asset, 0, len(externalAssets))
	toDeactive := make([]int, 0)

	// insert or update detection
	for key, ext := range externalAssets {
		loc, exists := local[key]

		if !exists {
			toUpsert = append(toUpsert, ext)
			continue
		}

		if ext.Name != loc.Name ||
			ext.Type != loc.Type ||
			ext.Branch != loc.Branch ||
			loc.Status != true {
			toUpsert = append(toUpsert, ext)
		}
	}

	for key, loc := range local {
		if _, exists := externalAssets[key]; !exists {
			if loc.Status == true {
				toDeactive = append(toDeactive, key.ID)
			}
		}
	}

	// db transaction
	tx, err := s.mapDB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// upsert batch
	var batch pgx.Batch

	for _, a := range toUpsert {
		batch.Queue(`
		INSERT INTO assets (
		asset_id,
		name,
		asset_type,
		branch,
		branch_id,
		is_active
		)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (asset_id)
		DO UPDATE SET
			name = EXCLUDED.name,
			asset_type = EXCLUDED.asset_type,
			branch = EXCLUDED.branch,
			branch_id = EXCLUDED.branch_id,
			is_active = TRUE
			`,
			a.ID,
			a.Name,
			a.Type,
			a.Branch,
			a.BranchID,
		)
	}
	br := tx.SendBatch(ctx, &batch)
	for range toUpsert {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return err
		}
	}
	br.Close()

	// bulk deactivate missing
	if len(toDeactive) > 0 {
		_, err = tx.Exec(ctx, `
		UPDATE assets
		SET is_active = FALSE
		WHERE asset_id = ANY($1)
		`, toDeactive)
		if err != nil {
			return err
		}
	}
	log.Printf("Asset Sync Completed")
	return tx.Commit(ctx)

}

func (s *AssetSyncService) loadExternalPCs(ctx context.Context) ([]Asset, error) {

	rows, err := s.invDB.Query(ctx, `
	SELECT id, pc_number, branch, status
FROM pc_inventory
WHERE status ILIKE 'active'
  AND branch NOT ILIKE 'n/a'
  AND pc_number NOT ILIKE 'laptop'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []Asset

	for rows.Next() {
		var a Asset
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Branch,
			&a.Status,
		); err != nil {
			return nil, err
		}
		a.Type = "pc"

		assets = append(assets, a)

	}
	return assets, nil
}

func (s *AssetSyncService) loadExternalPrinters(ctx context.Context) ([]Asset, error) {

	rows, err := s.invDB.Query(ctx, `
	SELECT id, printer_name, branch, status
	FROM printer_inventory
	WHERE status ILIKE 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []Asset

	for rows.Next() {
		var a Asset
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Branch,
			&a.Status,
		); err != nil {
			return nil, err
		}
		a.Type = "printer"

		if a.Name == "" || a.Branch == "" {
			continue
		}

		assets = append(assets, a)

	}
	return assets, nil
}

// will load current assets table from the gomap DB
func (s *AssetSyncService) loadLocalAssets(ctx context.Context) ([]LocalAsset, error) {

	rows, err := s.mapDB.Query(ctx, `
	SELECT asset_id, name, asset_type, branch, is_active
	FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []LocalAsset

	for rows.Next() {
		var a LocalAsset
		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Type,
			&a.Branch,
			&a.Status,
		); err != nil {
			return nil, err
		}

		if a.Name == "" || a.Branch == "" {
			continue
		}

		assets = append(assets, a)

	}
	return assets, rows.Err()
}

// Branch Name is not also uniform in the external database
// Will match defined aliases
func (s *AssetSyncService) loadBranchAliases(ctx context.Context) (map[string]string, error) {
	query := "SELECT LOWER(alias), branch_id FROM branch_aliases"
	rows, err := s.mapDB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	aliases := make(map[string]string)

	for rows.Next() {
		var alias string
		var branchID string
		if err := rows.Scan(&alias, &branchID); err != nil {
			return nil, err
		}

		aliases[alias] = branchID
	}
	return aliases, rows.Err()
}
