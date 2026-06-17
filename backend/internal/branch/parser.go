package branch

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	branch "github.com/PDK1744/gomap/internal/models"
)

type MXFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram Diagram  `xml:"diagram"`
}

type Diagram struct {
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

type MxGraphModel struct {
	Root Root `xml:"root"`
}

type Root struct {
	Objects []XMLObject `xml:"object"`
}

type XMLObject struct {
	TruID  string `xml:"tru_id,attr"`
	ID     string `xml:"id,attr"`
	MxCell MxCell `xml:"mxCell"`
}

type MxCell struct {
	Parent     string     `xml:"parent,attr"`
	Style      string     `xml:"style,attr"`
	MxGeometry MxGeometry `xml:"mxGeometry"`
}

type MxGeometry struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
}

type Parser struct {
	BaseDir    string
	BranchName string
}

// func NewParser(baseDir, branchName string)

func LoadAndNormalizeBranch(baseDir string, branchName string) (*branch.BranchResponse, error) {
	xmlPath := filepath.Join(baseDir, branchName, "layout.xml")
	jsonPath := filepath.Join(baseDir, branchName, "metadata.json")

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadataMap map[string]branch.MetadataItem
	if err := json.Unmarshal(jsonData, &metadataMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	xmlData, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read layout: %w", err)
	}

	var mxFile MXFile
	if err := xml.Unmarshal(xmlData, &mxFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	var elements []branch.MapElement
	xmlObjects := mxFile.Diagram.MxGraphModel.Root.Objects
	for _, obj := range xmlObjects {
		if obj.TruID == "" {
			continue
		}

		meta, found := metadataMap[obj.TruID]
		if !found {
			fmt.Printf("METADATA ISSUE X: %v Y: %v", obj.MxCell.MxGeometry.X, obj.MxCell.MxGeometry.Y)
			fmt.Printf("METADATA NOT FOUND for: %v", metadataMap[obj.TruID])
			continue
		}

		element := branch.MapElement{
			TruID:      obj.TruID,
			Name:       meta.Name,
			Type:       meta.Type,
			Assignable: meta.Assignable,
			Generator:  meta.Generator,
			Assets:     meta.Assets,
			X:          obj.MxCell.MxGeometry.X,
			Y:          obj.MxCell.MxGeometry.Y,
			Width:      obj.MxCell.MxGeometry.Width,
			Height:     obj.MxCell.MxGeometry.Height,
		}

		elements = append(elements, element)
	}

	return &branch.BranchResponse{
		BranchName: branchName,
		Elements:   elements,
	}, nil
}
