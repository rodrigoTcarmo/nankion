package nankion

import (
	"github.com/dstotijn/go-notion"
)

type Select struct {
	SelectOptions notion.SelectOptions `json:"select_options"`
}

type Date struct {
	Date notion.Date `json:"date"`
}

type Amount struct {
	NumberMetadata notion.NumberMetadata `json:"number_metadata"`
}
