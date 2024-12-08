package consumer

import (
	"github.com/arslanovdi/logistic-package/events/internal/config"
	"github.com/arslanovdi/logistic-package/events/internal/general"
	"github.com/arslanovdi/logistic-package/pkg/model"
	oteltracer "github.com/arslanovdi/otel-kafka-go"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const (
	readTimeout = 1 * time.Second
)

type KafkaConsumer struct {
	consumer     *kafka.Consumer
	deserializer *jsonschema.Deserializer
	stop         bool
	readers      sync.WaitGroup
	tracer       oteltracer.OtelConsumer
	group        string
}

func NewKafkaConsumer() (*KafkaConsumer, error) {
	log := slog.With("func", "consumer.MustNewKafkaConsumer")

	cfg := config.GetConfigInstance()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       strings.Join(cfg.Kafka.Brokers, ","),
		"group.id":                cfg.Kafka.GroupID,
		"enable.auto.commit":      "false",
		"auto.commit.interval.ms": 1000,
		"auto.offset.reset":       "earliest",
	})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaRegistry))
	if err != nil {
		return nil, err
	}

	deserializercfg := jsonschema.NewDeserializerConfig()
	deserializercfg.UseLatestVersion = true

	deserializer, err := jsonschema.NewDeserializer(sr, serde.ValueSerde, deserializercfg)
	if err != nil {
		return nil, err
	}

	tracer := oteltracer.NewOtelConsumer(cfg.Project.Instance)

	log.Info("KafkaConsumer created")

	return &KafkaConsumer{
		consumer:     consumer,
		deserializer: deserializer,
		stop:         false,
		tracer:       tracer,
		group:        cfg.Kafka.GroupID,
	}, nil
}

func (k *KafkaConsumer) Run(topic string, handler func(key string, msg model.PackageEvent, offset int64)) error {
	if k.stop {
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
			if err.(kafka.Error).IsTimeout() {
				continue
			}
			return err
		}

		k.tracer.OnPoll(kafkaMsg, k.group)

		event := model.PackageEvent{}
		err1 := k.deserializer.DeserializeInto(
			topic,
			kafkaMsg.Value,
			&event)
		if err1 != nil {
			return err1
		} else {
			k.tracer.OnProcess(kafkaMsg, k.group)

			handler(string(kafkaMsg.Key), event, int64(kafkaMsg.TopicPartition.Offset))

			_, err2 := k.consumer.CommitMessage(kafkaMsg)
			if err2 != nil {
				return err2
			}

			k.tracer.OnCommit(kafkaMsg, k.group)
		}

	}
	return nil
}

func (k *KafkaConsumer) Close() {

	log := slog.With("func", "consumer.Close")

	k.stop = true // command to stop reading messages
	k.readers.Wait()

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
