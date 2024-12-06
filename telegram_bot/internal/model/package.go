// Package model работа с моделью данных
package model

import (
	"fmt"
	pb "github.com/arslanovdi/logistic-package/pkg/logistic-package-api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log/slog"
	"strings"
	"time"
)

// Package - модель данных - пакет
type Package struct {
	ID      uint64
	Title   string
	Weight  *uint64
	Created time.Time
	Updated *time.Time
}

// String - строковое представление данных о пакет
func (c *Package) String() string {
	str := strings.Builder{}

	str.WriteString(fmt.Sprintf("ID: %d, Title: %s", c.ID, c.Title))
	if c.Weight != nil {
		str.WriteString(fmt.Sprintf(", Weight: %d", *c.Weight))
	}
	str.WriteString(fmt.Sprintf(", Created: %s", c.Created))
	if c.Updated != nil {
		str.WriteString(fmt.Sprintf(", Updated: %s", c.Updated))
	}

	return str.String()
}

// LogValue implements slog.LogValuer interface. Структурированный вывод данных Package в лог
func (c *Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("ID", c.ID),
		slog.String("Title", c.Title),
		slog.Any("Weight", *c.Weight),
		slog.Time("Created", c.Created),
		slog.Any("Updated", *c.Updated),
	)
}

// ToProto - конвертация данных о пакете в протобуф (DTO)
func (c *Package) ToProto() *pb.Package {

	var update *timestamp.Timestamp
	if c.Updated != nil {
		update = &timestamp.Timestamp{
			Seconds: c.Updated.Unix(),
			Nanos:   int32(c.Updated.Nanosecond()),
		}
	}

	return &pb.Package{
		Id:     c.ID,
		Title:  c.Title,
		Weight: c.Weight,
		Created: &timestamp.Timestamp{
			Seconds: c.Created.Unix(),
			Nanos:   int32(c.Created.Nanosecond()),
		},
		Updated: update,
	}
}

// FromProto - конвертация протобуф (DTO) в данные о пакете
func (c *Package) FromProto(pkg *pb.Package) {
	c.ID = pkg.Id
	c.Title = pkg.Title
	c.Weight = pkg.Weight
	c.Created = time.Unix(pkg.Created.Seconds, int64(pkg.Created.Nanos))
	if pkg.Updated != nil {
		c.Updated = &time.Time{}
		*c.Updated = time.Unix(pkg.Updated.Seconds, int64(pkg.Updated.Nanos))
	}
}
