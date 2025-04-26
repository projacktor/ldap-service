package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

)

// kafkaSyncer implements zapcore.WriteSyncer by sending each log entry to Kafka.
type kafkaSyncer struct {
	writer *kafka.Writer
}

func (k *kafkaSyncer) Write(p []byte) (n int, err error) {
	msg := kafka.Message{
		Key:   []byte("log"),
		Value: p,
        Time: time.Now().In(time.FixedZone("UTC+3", 3*60*60)),
	}

	// Пробуем записать сообщение несколько раз
	maxRetries := 3
	for i := range maxRetries {
		ctx := context.Background()
		if err := k.writer.WriteMessages(ctx, msg); err != nil {
			log.Printf("attempt %d/%d: failed to write to kafka: %v", i+1, maxRetries, err)
			if i == maxRetries-1 {
				return 0, err
			}
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return len(p), nil
	}
	return len(p), nil
}

func (k *kafkaSyncer) Sync() error {
	// kafka-go handles flushing internally
	return nil
}

func (k *kafkaSyncer) Close() error {
	return k.writer.Close()
}

// NewKafkaLogger constructs a zap.Logger that pushes Info+ logs to Kafka.
func NewKafkaLogger(brokers []string, topic string) (*zap.Logger, error) {
	// 1) Kafka writer
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		BatchSize:    1, // Отправляем каждое сообщение отдельно
		BatchTimeout: 10 * time.Millisecond,
		Async:        false, // Отключаем асинхронную отправку
	}

	// Проверяем соединение с Kafka
	if err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
        Time: time.Now().In(time.FixedZone("UTC+3", 3*60*60)),
	}); err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %v", err)
	}

	// 2) JSON encoder
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewJSONEncoder(encCfg)

	// 3) Wrap writer in WriteSyncer
	syncer := &kafkaSyncer{writer: w}
	ws := zapcore.AddSync(syncer)

	// 4) Core + logger
	core := zapcore.NewCore(enc, ws, zap.InfoLevel)
	logger := zap.New(core)

	return logger, nil
}
