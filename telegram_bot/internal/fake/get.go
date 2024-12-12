package fake

import (
	"log/slog"
)

// Генерирует запрос получения пакета
func (f *Faker) eGet() {
	if f.counter <= 2 { // отбрасываем запросы, пока не создано хотя-бы 2 пакета
		return
	}

	log := slog.With("func", "fake.eGet")
	id := genInt(f.counter) + 1

	_, err := f.pkgService.Get(id)
	if err != nil {
		log.Error("FAKE fail to get package", "id", id)
	}
}
