// src/logger/kafka-logger.service.ts
import { Injectable, LoggerService } from '@nestjs/common';
import { Kafka } from 'kafkajs';

@Injectable()
export class KafkaLoggerService implements LoggerService {
  private kafka = new Kafka({
    clientId: 'nestjs-logger',
    brokers: ['localhost:9092'],
  });

  private producer = this.kafka.producer();

  constructor() {
    this.producer.connect();
  }

  async log(message: string) {
    await this.producer.send({
      topic: 'nestjs-logs',
      messages: [{ value: `LOG: ${message}` }],
    });
    console.log(`Logged: ${message}`);
  }

  async error(message: string, trace?: string) {
    await this.producer.send({
      topic: 'nestjs-logs',
      messages: [{ value: `ERROR: ${message}, TRACE: ${trace}` }],
    });
    console.error(`Error: ${message}`);
  }

  async warn(message: string) {
    await this.producer.send({
      topic: 'nestjs-logs',
      messages: [{ value: `WARN: ${message}` }],
    });
    console.warn(`Warn: ${message}`);
  }
}
