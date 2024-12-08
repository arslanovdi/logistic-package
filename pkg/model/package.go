// Package model - структуры для работы с пакетами
package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log/slog"
	"strings"
	"time"
)

// Package структура пакета
type Package struct {
	ID      uint64        `db:"id" json:"ID"`
	Title   string        `db:"title" json:"title"`
	Weight  sql.NullInt64 `db:"weight" json:"weight,omitempty"`
	Created time.Time     `db:"created" json:"created"`
	Updated sql.NullTime  `db:"updated" json:"updated,omitempty"`
	Removed sql.NullBool  `db:"removed" json:"removed,omitempty"`
}

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

	//str := fmt.Sprintf("Package № %s %s. Package: %s", p.PackageID, p.Type, p.Payload)

	pkg := &Package{}
	json.Unmarshal(p.Payload, pkg)

	str := fmt.Sprintf("Package № %d %s. Package: %s", p.PackageID, p.Type, pkg)

	return str
}

// String implements fmt.Stringer
func (c *Package) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("ID: %d, Title: %s", c.ID, c.Title))
	if c.Weight.Valid {
		str.WriteString(fmt.Sprintf(", Weight: %d", c.Weight.Int64))
	}
	str.WriteString(fmt.Sprintf(", Created: %s", c.Created))
	if c.Updated.Valid {
		str.WriteString(fmt.Sprintf(", Updated: %s", c.Updated.Time))
	}
	if c.Removed.Valid {
		str.WriteString(fmt.Sprintf(", Removed: %t", c.Removed.Bool))
	}
	return str.String()
}

// LogValue implements slog.LogValuer interface
// вывод Package в лог
func (c *Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("ID", c.ID),
		slog.String("Title", c.Title),
		slog.Any("Weight", c.Weight),
		slog.Time("Created", c.Created),
		slog.Any("Updated", c.Updated),
		slog.Any("Removed", c.Removed),
	)
}

// ToProto converts model.Package to pb.Package
func (c *Package) ToProto() *pb.Package {
	// проверка опциональных полей
	var weight *uint64
	if c.Weight.Valid { // Если указан вес
		t := uint64(c.Weight.Int64)
		weight = &t
	}
	var updated *timestamp.Timestamp // если есть updated time
	if c.Updated.Valid {
		updated = &timestamp.Timestamp{
			Seconds: c.Updated.Time.Unix(),
			Nanos:   int32(c.Updated.Time.Nanosecond()),
		}
	}

	return &pb.Package{
		Id:     c.ID,
		Title:  c.Title,
		Weight: weight,
		Created: &timestamp.Timestamp{
			Seconds: c.Created.Unix(),
			Nanos:   int32(c.Created.Nanosecond()),
		},
		Updated: updated,
	}
}

// FromProto converts pb.Package to model.Package
func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id

	c.Title = pkg.Title

	if pkg.Weight != nil {
		c.Weight = sql.NullInt64{Int64: int64(*pkg.Weight), Valid: true}
	} else {
		c.Weight = sql.NullInt64{}
	}

	c.Created = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))

	if pkg.Updated != nil {
		c.Updated = sql.NullTime{Time: time.Unix(pkg.Updated.Seconds, int64(pkg.Updated.Nanos)), Valid: true}
	} else {
		c.Updated = sql.NullTime{}
	}
}
