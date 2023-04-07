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
	ID            uint64
	Title         string
	Material      string
	MaximumVolume float32 //cm^3
	Reusable      bool
}

var PackageFieldsCount int = 5

func (s *Package) String() string {
	return s.Title + " [id: " + strconv.FormatUint(s.ID, 10) + "]"
}
