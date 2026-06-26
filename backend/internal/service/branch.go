package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PDK1744/gomap/internal/models"
	"github.com/PDK1744/gomap/internal/store"
)

type BranchService struct {
	branchStore *store.BranchStore
}

const BASE_DIR = "./branches_local"

func NewBranchService(store *store.BranchStore) *BranchService {
	return &BranchService{branchStore: store}
}

func (s *BranchService) FetchBranchLayout(ctx context.Context, branchName string) (*models.BranchResponse, error) {
	objects, err := loadXMLObjects(branchName)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]: loading xml file: %v", err)
	}
	truIDs := make([]string, 0, len(objects))
	for _, obj := range objects {
		if obj.TruID != "" {
			truIDs = append(truIDs, obj.TruID)
		}
	}

	metaData, err := s.branchStore.FetchMetadataBatch(ctx, truIDs)
	if err != nil {
		return nil, err
	}

	var elements []models.MapElement
	for _, obj := range objects {
		meta, ok := metaData.Map[obj.TruID]
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
