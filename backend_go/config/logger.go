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

// kafkaSyncer implements zapcore.WriteSyncer interface for Kafka integration
// It handles writing log messages to Kafka with retry logic
type kafkaSyncer struct {
	writer *kafka.Writer // Kafka writer instance
}

// Write sends log messages to Kafka with retry mechanism
// Implements io.Writer interface for zap integration
// Parameters:
//   - p: Byte slice containing the log message
//
// Returns:
//   - n: Number of bytes written
//   - err: Error if all retries failed
func (k *kafkaSyncer) Write(p []byte) (n int, err error) {
	msg := kafka.Message{
		Key:   []byte("log"),                                   // Fixed key for all log messages
		Value: p,                                               // Actual log message content
		Time:  time.Now().In(time.FixedZone("UTC+3", 3*60*60)), // Timestamp with timezone
	}

	// Retry mechanism for Kafka writes
	maxRetries := 3
	for i := range maxRetries {
		ctx := context.Background()
		if err := k.writer.WriteMessages(ctx, msg); err != nil {
			log.Printf("attempt %d/%d: failed to write to kafka: %v", i+1, maxRetries, err)
			if i == maxRetries-1 { // Final attempt failed
				return 0, err
			}
			// Exponential backoff would be better here
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return len(p), nil
	}
	return len(p), nil
}

// Sync implements zapcore.WriteSyncer interface
// No-op since kafka-go handles flushing internally
func (k *kafkaSyncer) Sync() error {
	return nil
}

// Close cleans up the Kafka writer resources
func (k *kafkaSyncer) Close() error {
	return k.writer.Close()
}

// NewKafkaLogger creates a zap.Logger that sends logs to Kafka
// Parameters:
//   - brokers: List of Kafka broker addresses
//   - topic: Kafka topic to write logs to
//
// Returns:
//   - *zap.Logger: Configured logger instance
//   - error: If Kafka connection test fails
func NewKafkaLogger(brokers []string, topic string) (*zap.Logger, error) {
	// 1) Configure Kafka writer
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...), // Broker addresses
		Topic:        topic,                 // Target topic
		Balancer:     &kafka.LeastBytes{},   // Partition balancing strategy
		RequiredAcks: kafka.RequireOne,      // Wait for leader acknowledgment
		BatchSize:    1,                     // No batching - immediate send
		BatchTimeout: 10 * time.Millisecond, // Max time to wait for batch
		Async:        false,                 // Synchronous mode for reliability
	}

	// 2) Test Kafka connection with probe message
	if err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: []byte("test"),
		Time:  time.Now().In(time.FixedZone("UTC+3", 3*60*60)),
	}); err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %v", err)
	}

	// 3) Configure JSON encoder for log formatting
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder   // Human-readable timestamps
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder // INFO/WARN/ERROR formatting
	enc := zapcore.NewJSONEncoder(encCfg)

	// 4) Create synchronized writer
	syncer := &kafkaSyncer{writer: w}
	ws := zapcore.AddSync(syncer)

	// 5) Create and return logger core
	core := zapcore.NewCore(
		enc,           // Encoder
		ws,            // WriteSyncer
		zap.InfoLevel, // Minimum log level
	)
	logger := zap.New(core)

	return logger, nil
}
