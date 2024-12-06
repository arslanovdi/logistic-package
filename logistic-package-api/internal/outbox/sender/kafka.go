// Package sender send events to kafka
package sender

import (
	"context"
	"errors"
	"fmt"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/config"
	"github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	oteltracer "github.com/arslanovdi/otel-kafka-go"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

// EventSender - интерфейс для отправки событий в кафку
type EventSender interface {
	// Send отправка сообщения в kafka. key = pkg.id (string); value = pkg.
	Send(ctx context.Context, pkg *model.PackageEvent, topic string) error
	Close() error
}

type kafkaSender struct {
	producer     *kafka.Producer
	serializer   serde.Serializer
	flushTimeout int // milliseconds
	initialized  bool
	stop         chan struct{}  // Команда на принудительное завершение produce
	senders      sync.WaitGroup // Количество одновременных produce в kafka, их может быть несколько горутин
	tracer       oteltracer.OtelProducer
}

// Send отправить событие в kafka, avro сериализация.
// ctx - контекст с трассировкой от OpenTelemetry
func (k *kafkaSender) Send(ctx context.Context, pkg *model.PackageEvent, topic string) error {
	if !k.initialized {
		return model.ErrProducerClosed
	}

	log := slog.With("func", "kafkaSender.Send")

	k.senders.Add(1)
	defer k.senders.Done()

	payload, err := k.serializer.Serialize(topic, *pkg)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(fmt.Sprintf("%d", pkg.PackageID)), // TODO ID ?
		Value: payload,
	}

	k.tracer.OnSend(ctx, &msg) // OpenTelemtry trace

	err = k.producer.Produce(&msg, deliveryChan)
	if err != nil {
		return err
	}

	log.Debug("send event", slog.Any("event", pkg))

	select {
	case <-k.stop:
		return model.ErrNoDeliveredMessage

	case event := <-deliveryChan:
		switch ev := event.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return ev.TopicPartition.Error
			}
			return nil
		case kafka.Error:
			return ev
		}
		return errors.New(event.String()) // Неизвестная ошибка
	}
}

// MustNewKafkaSender конструктор.
// Глобальный провайдер OpenTelemetry должен быть инициализирован.
func MustNewKafkaSender() EventSender {
	log := slog.With("func", "MustNewKafkaSender")

	cfg := config.GetConfigInstance()
	brokers := ""
	for _, b := range cfg.Kafka.Brokers {
		brokers = brokers + b + ","
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[:len(brokers)-1],
		"acks":              "all",
	})
	if err != nil {
		log.Warn("Failed to create kafka producer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry))
	if err != nil {
		log.Warn("Failed to create kafka schema registry", slog.String("error", err.Error()))
		os.Exit(1)
	}

	jsoncfg := jsonschema.NewSerializerConfig()
	jsoncfg.AutoRegisterSchemas = true

	serializer, err := jsonschema.NewSerializer(
		sr,
		serde.ValueSerde,
		jsoncfg)

	if err != nil {
		log.Warn("Failed to create avro serializer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tracer := oteltracer.NewOtelProducer(cfg.Project.Instance)

	stop := make(chan struct{})

	return &kafkaSender{
		producer:     producer,
		serializer:   serializer,
		flushTimeout: cfg.Kafka.FlushTimeout,
		initialized:  true,
		senders:      sync.WaitGroup{},
		tracer:       tracer,
		stop:         stop,
	}
}

// Close - останавливает отправку сообщений, с таймаутом cfg.Kafka.flushTimeout ms.
func (k *kafkaSender) Close() (err error) {
	err = nil

	log := slog.With("func", "kafkaSender.Close")

	k.initialized = false // command to block future sends

	outEvents := k.producer.Flush(k.flushTimeout) // Wait for message deliveries before shutting down
	if outEvents > 0 {
		err = errors.New("Kafka producer outstanding events count: " + strconv.Itoa(outEvents))
	}

	close(k.stop)

	k.senders.Wait() // waiting for the completion of the launched message sending

	err2 := k.serializer.Close()

	if err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = errors.Join(err, err2)
		}
	}

	k.producer.Close()

	log.Info("Kafka producer closed")

	return err
}
