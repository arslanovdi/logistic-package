// Package cache - кэширование запросов.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/general"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/extra/redisprometheus/v9"

	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/service"
	"github.com/arslanovdi/logistic-package/pkg/model"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	startTimeout = 5 * time.Second // Время ожидания подключения к Redis
)

/*
Redis - реализация кэша.
Если проблемы с консистентностью, кэш реинициализируется, ретраями с шагом cfg.Redis.RetryDuration.
При недоступности redis, данные идут минуя кэш в/из БД.
*/
type Redis struct {
	rdb     *redis.Client
	repo    service.Repo
	healthy *atomic.Bool // Кэш консистентен
	ttl     time.Duration
	stop    chan struct{}
}

// Использовать, если кэш не синхронизирован с базой данных
func (c *Redis) restart() {
	log := slog.With("func", "cache.Restart")

	err := c.rdb.FlushAll(context.Background()).Err()
	if err != nil {
		log.Error("Redis restart error", slog.String("error", err.Error()))
		return
	}
	log.Info("Redis restarted")
	c.healthy.Store(true)
}

// IsConsistent - кэш консистентен
func (c *Redis) IsConsistent() bool {
	return c.healthy.Load()
}

// Приведем []interface{} к []model.Package
// так как redis.MGet возвращает []interface{}
func castToPackages(list []any) (pkgCache []model.Package, err error) {
	pkgCache = make([]model.Package, len(list))

	for i, record := range list {
		if record == nil {
			return nil, general.ErrNotFound
		}

		str, ok := record.(string)
		if !ok {
			return nil, errors.New("error on cast interface{} to string")
		}

		err = json.Unmarshal([]byte(str), &pkgCache[i]) // Развернуть в структуру
		if err != nil {
			return nil, err
		}
	}

	return pkgCache, err
}

// Create - создание нового пакета в БД -> кэше
func (c *Redis) Create(ctx context.Context, pkg *model.Package) (*uint64, error) {
	log := slog.With("func", "cache.Create")

	id, err := c.repo.Create(ctx, pkg)
	if err != nil {
		return id, err
	}

	if !c.healthy.Load() { // нарушена консистентность кэша, пропускаем операции с кэшем
		return id, nil
	}

	// добавить запись в кэш
	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	err1 := c.rdb.Set(ctx, strconv.FormatUint(pkg.ID, 10), pkg, c.ttl).Err()
	if err1 != nil {
		c.healthy.Store(false) // нарушена консистентность кэша

		log.Error("Cache consistency is broken", slog.String("error", err1.Error()))
	}

	log.Debug("package created in cache", slog.String("package", pkg.String()))

	return id, err
}

// Delete - удаление пакета из БД -> кэша
func (c *Redis) Delete(ctx context.Context, id uint64) error {
	log := slog.With("func", "cache.Delete")

	err := c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if !c.healthy.Load() { // нарушена консистентность кэша, пропускаем операции с кэшем
		return err
	}

	// удалить запись в кэше
	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	err1 := c.rdb.Del(ctx, strconv.FormatUint(id, 10)).Err()
	if err1 != nil {
		c.healthy.Store(false) // нарушена консистентность кэша
		log.Error("Cache consistency is broken", slog.String("error", err1.Error()))
		return err
	}

	log.Debug("package deleted from cache", slog.Uint64("package id", id))

	return err
}

// Get - получение пакета из кэша -> БД
func (c *Redis) Get(ctx context.Context, id uint64) (*model.Package, error) {
	log := slog.With("func", "cache.Get")

	if c.healthy.Load() { // кэш консистентен, проверить наличие в кэше
		pkgCache := model.Package{}

		// получить запись из кэша
		ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

		err := c.rdb.Get(ctx, strconv.FormatUint(id, 10)).Scan(&pkgCache)

		if err == nil {
			// запись найдена в кэше
			log.Debug("package found in cache", slog.Uint64("package id", id))

			return &pkgCache, nil
		}
	}

	pkgDB, err1 := c.repo.Get(ctx, id)
	if err1 == nil && pkgDB != nil && c.healthy.Load() { // нет ошибки, есть запись, кэш консистентен
		// добавить запись в кэш
		ctx = context.WithoutCancel(ctx)                               // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился
		c.rdb.Set(ctx, strconv.FormatUint(pkgDB.ID, 10), pkgDB, c.ttl) // ошибку не проверяем, операция не может нарушить консистентность кэша
	}

	return pkgDB, err1
}

// List - получение пакетов из кэша -> БД
func (c *Redis) List(ctx context.Context, offset, limit uint64) ([]model.Package, error) {
	log := slog.With("func", "cache.List")

	if c.healthy.Load() { // кэш консистентен, проверить наличие в кэше
		keys := make([]string, limit)
		for i := range limit {
			keys[i] = strconv.FormatUint(offset+i, 10)
		}

		// получить записи из кэша
		list, err := c.rdb.MGet(ctx, keys...).Result() // Возвращает []interface{}.

		if err == nil {
			pkgCache, err1 := castToPackages(list)
			if err1 == nil { // все записи нашлись в кэше
				log.Debug("packages found in cache", slog.Uint64("offset", offset), slog.Uint64("limit", limit))
				return pkgCache, nil
			}
			if !errors.Is(err1, general.ErrNotFound) {
				log.Error("cast to struct", slog.String("error", err1.Error()))
			}
		}
	}

	// получить записи из базы
	pkgDB, err1 := c.repo.List(ctx, offset, limit)

	if err1 == nil && pkgDB != nil && c.healthy.Load() { // нет ошибки, есть записи, кэш консистентен
		// добавить запись в кэш
		ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

		pipe := c.rdb.Pipeline() // создать пайплайн
		for _, pkg := range pkgDB {
			pipe.Set(ctx, strconv.FormatUint(pkg.ID, 10), &pkg, c.ttl) // добавить команду в пайплайн
		}
		_, err2 := pipe.Exec(ctx) // выполнить пайплайн
		if err2 != nil {
			log.Error("error executing pipeline", slog.String("error", err2.Error()))
		}
	}
	return pkgDB, err1
}

// Update - изменение пакета в БД -> кэше
func (c *Redis) Update(ctx context.Context, pkg *model.Package) error {
	log := slog.With("func", "cache.Update")

	err := c.repo.Update(ctx, pkg)
	if err != nil {
		return err
	}

	if !c.healthy.Load() { // нарушена консистентность кэша, пропускаем операции с кэшем
		return err
	}

	// изменить запись в кэше
	ctx = context.WithoutCancel(ctx) // Отвязать контекст, нужно завершить операцию, даже если клиент отвалился

	err1 := c.rdb.Set(ctx, strconv.FormatUint(pkg.ID, 10), pkg, c.ttl).Err()
	if err1 != nil {
		c.healthy.Store(false) // нарушена консистентность кэша

		log.Error("Cache consistency is broken", slog.String("error", err1.Error()))
		return err
	}

	log.Debug("package updated in cache", slog.String("package", pkg.String()))

	return err
}

// NewRedis - конструктор
func NewRedis(repo service.Repo) *Redis {
	log := slog.With("func", "cache.NewRedis")

	cfg := config.GetConfigInstance()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password, // No password set ""
		DB:       cfg.Redis.DB,       // Use default DB (0)
		Protocol: 3,                  // Connection protocol

	})

	ctx, cancel := context.WithTimeout(context.Background(), startTimeout)
	defer cancel()

	var healthy atomic.Bool
	healthy.Store(true)

	_, err := client.Ping(ctx).Result()
	if err != nil { // connection error, don't panic
		log.Error("failed to connect to redis", slog.String("error", err.Error()))
		healthy.Store(false)
	}

	collector := redisprometheus.NewCollector("logistic_package", "redis", client) // enable redis metrics to Prometheus
	prometheus.MustRegister(collector)
	/*
		| `pool_hit_total`          | Counter metric | number of times a connection was found in the pool                          |
		| `pool_miss_total`         | Counter metric | number of times a connection was not found in the pool                      |
		| `pool_timeout_total`      | Counter metric | number of times a timeout occurred when getting a connection from the pool  |
		| `pool_conn_total_current` | Gauge metric   | current number of connections in the pool                                   |
		| `pool_conn_idle_current`  | Gauge metric   | current number of idle connections in the pool                              |
		| `pool_conn_stale_total`   | Counter metric | number of times a connection was removed from the pool because it was stale |
	*/

	err = redisotel.InstrumentTracing(client) // enable redis OpenTelemetry tracing
	if err != nil {
		log.Error("failed enable redis tracing", slog.String("error", err.Error()))
	}

	cache := &Redis{
		rdb:     client,
		repo:    repo,
		healthy: &healthy,
		ttl:     time.Second * time.Duration(cfg.Redis.TTL),
		stop:    make(chan struct{}),
	}

	// Очистка кэша перед работой, возможно он сохранялся на диск.
	cache.restart()

	// redis consistency monitor
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Redis.RetryDuration) * time.Second)
		for {
			select {
			case <-ticker.C:
				if !cache.healthy.Load() {
					cache.restart()
				}
			case <-cache.stop:
				return
			}
		}
	}()

	return cache
}

// Close - закрытие кэша
func (c *Redis) Close() {
	log := slog.With("func", "cache.Close")

	close(c.stop)

	err := c.rdb.Close()
	if err != nil {
		log.Error("failed to close redis", slog.String("error", err.Error()))
		return
	}

	log.Info("redis closed correctly")
}
