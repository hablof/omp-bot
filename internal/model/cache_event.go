package model

import (
	"encoding/json"
	"time"
)

type CacheEventType uint8

const (
	_ CacheEventType = iota
	SetDescription
	ReadDescription
	RemoveDescription
)

func (e CacheEventType) MarshalJSON() ([]byte, error) {
	var s string
	switch e {
	case SetDescription:
		s = "Set description"
	case ReadDescription:
		s = "Read description"
	case RemoveDescription:
		s = "Remove description"
	}

	return json.Marshal(s)
}

type CacheEvent struct {
	PackageID uint64         `json:"id"`
	EventType CacheEventType `json:"cacheEventType"`
	Success   bool           `json:"success"`
	Timestamp time.Time      `json:"timestamp"`
}
