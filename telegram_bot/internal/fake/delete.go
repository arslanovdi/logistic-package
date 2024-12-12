package fake

import (
	"log/slog"
)

// Генерирует запрос удаления пакета
func (f *Faker) eDelete() {
	if f.counter <= 2 { // отбрасываем запросы, пока не создано хотя-бы 2 пакета
		return
	}

	log := slog.With("func", "fake.eDelete")

	id := genInt(f.counter) + 1

	err := f.pkgService.Delete(id)
	if err != nil {
		log.Error("FAKE fail to delete package", "error", err)
	}

	f.counter-- // счетчик пакетов
}
