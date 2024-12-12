// Package model - DTO структуры
package model

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	var weight *int64
	if c.Weight.Valid { // Если указан вес
		t := c.Weight.Int64
		weight = &t
	}
	var updated *timestamp.Timestamp // если есть updated time
	if c.Updated.Valid {
		updated = timestamppb.New(c.Updated.Time)
	}
	created := timestamppb.New(c.Created)

	return &pb.Package{
		Id:      c.ID,
		Title:   c.Title,
		Weight:  weight,
		Created: created,
		Updated: updated,
	}
}

// FromProto converts pb.Package to model.Package
func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id

	c.Title = pkg.Title

	if pkg.Weight != nil {
		c.Weight = sql.NullInt64{Int64: *pkg.Weight, Valid: true}
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
