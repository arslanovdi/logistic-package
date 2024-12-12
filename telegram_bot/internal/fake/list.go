package fake

import (
	"log/slog"
)

// Генерирует запрос получения списка пакетов
func (f *Faker) eList() {
	if f.counter < listCount { // отбрасываем запросы, пока не создано хотя-бы listCount пакетов (10)
		return
	}

	log := slog.With("func", "fake.eList")

	offset := genInt(f.counter - listCount) // случайное смещение, до половины созданных пакетов
	limit := genInt(listCount)

	_, err := f.pkgService.List(offset, limit)

	if err != nil {
		log.Error("FAKE fail to list package", "error", err)
	}
}
