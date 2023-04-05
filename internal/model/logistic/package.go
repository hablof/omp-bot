package logistic

import (
	"reflect"
	"strconv"
)

const (
	Title         = "название="
	Material      = "материал="
	MaximumVolume = "объём="
	Reusable      = "многоразовая="
)

type Package struct {
	ID            uint64
	Title         string
	Material      string
	MaximumVolume float32 //cm^3
	Reusable      bool
}

var PackageFieldsCount int = reflect.TypeOf(Package{}).NumField()

func (s *Package) String() string {
	return s.Title + " id:" + strconv.FormatUint(s.ID, 10)
}
