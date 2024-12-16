// Package model - DTO структуры
package model

import (
	"database/sql"
	"encoding/json"
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
	ID      uint64        `db:"id" json:"id"`
	Title   string        `db:"title" json:"title"`
	Weight  sql.NullInt64 `db:"weight" json:"weight,omitempty"`
	Created time.Time     `db:"created" json:"created"`
	Updated sql.NullTime  `db:"updated" json:"updated,omitempty"`
}

// String implements fmt.Stringer
func (p *Package) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("ID: %d, Title: %s", p.ID, p.Title))
	if p.Weight.Valid {
		str.WriteString(fmt.Sprintf(", Weight: %d", p.Weight.Int64))
	}
	str.WriteString(fmt.Sprintf(", Created: %s", p.Created))
	if p.Updated.Valid {
		str.WriteString(fmt.Sprintf(", Updated: %s", p.Updated.Time))
	}
	return str.String()
}

// LogValue implements slog.LogValuer interface
// вывод Package в лог
func (p *Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("ID", p.ID),
		slog.String("Title", p.Title),
		slog.Any("Weight", p.Weight),
		slog.Time("Created", p.Created),
		slog.Any("Updated", p.Updated),
	)
}

// ToProto converts model.Package to pb.Package
func (p *Package) ToProto() *pb.Package {
	// проверка опциональных полей
	var weight *int64
	if p.Weight.Valid { // Если указан вес
		t := p.Weight.Int64
		weight = &t
	}
	var updated *timestamp.Timestamp // если есть updated time
	if p.Updated.Valid {
		updated = timestamppb.New(p.Updated.Time)
	}
	created := timestamppb.New(p.Created)

	return &pb.Package{
		Id:      p.ID,
		Title:   p.Title,
		Weight:  weight,
		Created: created,
		Updated: updated,
	}
}

// FromProto converts pb.Package to model.Package
func (p *Package) FromProto(pkg *pb.Package) {
	p.ID = pkg.Id
	p.Title = pkg.Title

	if pkg.Weight != nil {
		p.Weight = sql.NullInt64{Int64: *pkg.Weight, Valid: true}
	} else {
		p.Weight = sql.NullInt64{}
	}

	p.Created = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))

	if pkg.Updated != nil {
		p.Updated = sql.NullTime{Time: time.Unix(pkg.Updated.Seconds, int64(pkg.Updated.Nanos)), Valid: true}
	} else {
		p.Updated = sql.NullTime{}
	}
}

// MarshalBinary реализует encoding.BinaryMarshaller, для работы с Redis
func (p *Package) MarshalBinary() ([]byte, error) {
	return json.Marshal(*p)
}

// UnmarshalBinary реализует encoding.BinaryUnmarshaler, для работы с Redis
func (p *Package) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
