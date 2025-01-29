// Package consumer чтение событий из kafka
package consumer

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/arslanovdi/logistic-package/events/internal/config"
	"github.com/arslanovdi/logistic-package/events/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/model"
	oteltracer "github.com/arslanovdi/otel-kafka-go"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

const (
	readTimeout = 1 * time.Second
)

// KafkaConsumer - чтение из кафки
type KafkaConsumer struct {
	consumer     *kafka.Consumer
	deserializer *jsonschema.Deserializer
	stop         bool                    // consumer stop sign
	readers      sync.WaitGroup          // wait group для корректного закрытия consumer
	tracer       oteltracer.OtelConsumer // Интерфейс для отправки трассировки
	group        string
}

// NewKafkaConsumer конструктор
func NewKafkaConsumer() (*KafkaConsumer, error) {
	log := slog.With("func", "consumer.MustNewKafkaConsumer")

	cfg := config.GetConfigInstance()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(cfg.Kafka.Brokers, ","),
		"group.id":           cfg.Kafka.GroupID,
		"enable.auto.commit": "false",
		"auto.offset.reset":  "earliest",
		// "auto.commit.interval.ms": 1000,
	})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry))
	if err != nil {
		return nil, err
	}

	deserializerConfig := jsonschema.NewDeserializerConfig()
	deserializerConfig.UseLatestVersion = true

	deserializer, err := jsonschema.NewDeserializer(sr, serde.ValueSerde, deserializerConfig)
	if err != nil {
		return nil, err
	}

	tracer := oteltracer.NewOtelConsumer(cfg.Project.Instance) // интерфейс отправки трассировки

	log.Info("KafkaConsumer created")

	return &KafkaConsumer{
		consumer:     consumer,
		deserializer: deserializer,
		stop:         false,
		tracer:       tracer,
		group:        cfg.Kafka.GroupID,
	}, nil
}

// Run read and process messages from kafka
func (k *KafkaConsumer) Run(topic string, handler func(ctx context.Context, key string, msg model.PackageEvent, offset int64)) error {
	if k.stop { // check if consumer is closed
		return general.ErrConsumerClosed
	}
	k.readers.Add(1)
	defer k.readers.Done()

	err := k.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}

	for !k.stop {
		kafkaMsg, err := k.consumer.ReadMessage(readTimeout) // read message with timeout
		if err != nil {
			var e kafka.Error
			ok := errors.As(err, &e)
			if ok {
				if e.IsTimeout() { // `err.(kafka.Error).IsTimeout() == true`
					continue
				}
			}

			return err
		}

		k.tracer.OnPoll(kafkaMsg, k.group) // send message to tracer

		event := model.PackageEvent{}
		err1 := k.deserializer.DeserializeInto(
			topic,
			kafkaMsg.Value,
			&event)
		if err1 != nil {
			return err1
		}

		k.tracer.OnProcess(kafkaMsg, k.group) // send message to tracer

		ctxWithTrace := oteltracer.Context(kafkaMsg) // traceid from message to context

		handler(ctxWithTrace, string(kafkaMsg.Key), event, int64(kafkaMsg.TopicPartition.Offset)) // call handler

		_, err2 := k.consumer.CommitMessage(kafkaMsg)
		if err2 != nil {
			return err2
		}

		k.tracer.OnCommit(kafkaMsg, k.group) // send message to tracer
	}
	return nil
}

// Close kafka consumer with waiting for processing to complete
func (k *KafkaConsumer) Close() {
	log := slog.With("func", "consumer.Close")

	k.stop = true    // command to stop reading messages
	k.readers.Wait() // waiting for processing to complete

	err := k.deserializer.Close()
	if err != nil {
		log.Error("Failed to close deserializer: ", slog.String("error", err.Error()))
	}
	err = k.consumer.Close()
	if err != nil {
		log.Error("Failed to close consumer: ", slog.String("error", err.Error()))
	}

	log.Info("KafkaConsumer closed")
}
