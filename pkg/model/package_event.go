package model

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// EventType тип события
type EventType uint8

// EventStatus статус события
type EventStatus uint8

const (
	_       EventType = iota
	Created           // Created - события создания пакета
	Updated           // Updated - события изменения пакета
	Removed           // Removed - события удаления пакета
)

const (
	_        EventStatus = iota
	Locked               // Locked - событие заблокировано (отправляется в кафку)
	Unlocked             // Unlocked - событие разблокировано (находится в очереди для отправки в кафку)
)

// PackageEvent структура события
type PackageEvent struct {
	ID        int64       `db:"id" json:"ID,omitempty"`
	PackageID int64       `db:"package_id" json:"packageID,omitempty"`
	Type      EventType   `db:"type" json:"type,omitempty"`
	Status    EventStatus `db:"status" json:"status,omitempty"`
	Payload   []byte      `db:"payload" json:"payload,omitempty"`
	TraceID   *string     `db:"traceid" json:"-"`
}

func (e EventType) String() string {
	switch e {
	case Created:
		return "created"
	case Updated:
		return "updated"
	case Removed:
		return "removed"
	default:
		return "unknown"
	}
}

func (p *PackageEvent) String() string {
	if p.Type == Removed {
		return fmt.Sprintf("Package № %d %s", p.PackageID, p.Type)
	}

	log := slog.With("func", "model.PackageEvent.String")

	pkg := &Package{}
	err := json.Unmarshal(p.Payload, pkg)
	if err != nil {
		log.Error("Failed to unmarshal payload", slog.String("error", err.Error()))
	}

	return fmt.Sprintf("Package № %d %s. Package: %s", p.PackageID, p.Type, pkg)
}
