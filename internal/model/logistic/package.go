package logistic

import (
	"strconv"
)

const (
	Title         = "название"
	Material      = "материал"
	MaximumVolume = "объём"
	Reusable      = "многоразовая"
)

type Package struct {
	ID            uint64  `json:"id"`
	Title         string  `json:"title"`
	Material      string  `json:"material"`
	MaximumVolume float32 `json:"volume"`
	Reusable      bool    `json:"reusable"`
}

var PackageFieldsCount int = 5

func (s *Package) String() string {
	return s.Title + " [id: " + strconv.FormatUint(s.ID, 10) + "]"
}
