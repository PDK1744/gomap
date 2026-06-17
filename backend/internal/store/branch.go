package store

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

func (bs *BranchStore) FetchBranchLayout(ctx context.Context, branchName string) (*models.BranchResponse, error) {
	if branchName == "" {
		return nil, fmt.Errorf("branch name is required: %v", branchName)
	}
	objects, err := loadXMLObjects(branchName)
	if err != nil {
		return nil, fmt.Errorf("ERROR: loading xmlFile: %v", err)
	}

	truIDs := make([]string, 0, len(objects))
	for _, obj := range objects {
		if obj.TruID != "" {
			truIDs = append(truIDs, obj.TruID)
		}
	}
	metaMap, err := bs.fetchMetadataBatch(ctx, truIDs)
	if err != nil {
		return nil, err
	}
	var elements []models.MapElement
	for _, obj := range objects {
		meta, ok := metaMap[obj.TruID]
		if !ok {
			continue
		}
		element := models.MapElement{
			TruID:      obj.TruID,
			Name:       meta.Name,
			Type:       meta.Type,
			Assignable: meta.Assignable,
			Generator:  meta.Generator,
			X:          obj.MxCell.MxGeometry.X,
			Y:          obj.MxCell.MxGeometry.Y,
			Width:      obj.MxCell.MxGeometry.Width,
			Height:     obj.MxCell.MxGeometry.Height,
		}
		elements = append(elements, element)
	}

	return &models.BranchResponse{
		BranchName: branchName,
		Elements:   elements,
	}, nil

}

// need to load XML layout but join the metadata from the Database
// Im looking at LoadAndNormalizeBranch to refactor it for this usage
// should loadLayout just return the MXFile?
// or should it return MapElement and then i can append the metadata after that is fetched?
// But the MXFile struct is in /branch/parser.go so I'm not sure where to define the loadLayout???
func loadXMLObjects(branchName string) ([]models.XMLObject, error) {
	xmlPath := filepath.Join(BASE_DIR, branchName, "layout.xml")

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		return []models.XMLObject{}, fmt.Errorf("failed to read layout: %w", err)
	}

	var mxFile models.MXFile
	if err := xml.Unmarshal(xmlData, &mxFile); err != nil {
		return []models.XMLObject{}, fmt.Errorf("failed to unmarshal XML: %v", err)
	}
	xmlObject := mxFile.Diagram.MxGraphModel.Root.Objects

	return xmlObject, nil
}

func (bs *BranchStore) fetchMetadataBatch(ctx context.Context, truIDs []string) (map[string]models.MetadataItem, error) {
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
	return metaMap, nil
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
