package models

import "encoding/xml"

// MetadataItem reflects the values inside your JSON file
type MetadataItem struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Assignable bool    `json:"assignable"`
	Parent     string  `json:"parent,omitempty"`
	Generator  bool    `json:"generator,omitempty"`
	Assets     []Asset `json:"assets,omitempty"`
}

// MapElement is what our frontend will consume.
// It combines coordinates from XML and metadata from JSON.
type MapElement struct {
	TruID      string  `json:"tru_id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Assignable bool    `json:"assignable"`
	Generator  bool    `json:"generator,omitempty"`
	Assets     []Asset `json:"assets,omitempty"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
}

type BranchResponse struct {
	BranchName string       `json:"branch_name"`
	Elements   []MapElement `json:"elements"`
}

type Asset struct {
	Name string
}

// XML
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
