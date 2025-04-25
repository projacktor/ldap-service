// src/logger/kafka-logger.service.ts
import { Injectable, LoggerService, OnModuleInit } from '@nestjs/common';
import { Kafka } from 'kafkajs';

@Injectable()
export class KafkaLoggerService implements LoggerService, OnModuleInit {
  private readonly TOPIC = 'nestjs-logs';
  private kafka = new Kafka({
    clientId: 'nestjs-logger',
    brokers: [process.env.KAFKA_BROKERS || 'kafka:29092'],
  });

  private producer = this.kafka.producer();
  private admin = this.kafka.admin();

  constructor() {
    this.initialize();
  }

  async onModuleInit() {
    await this.createTopic();
  }

  private async initialize() {
    try {
      await this.producer.connect();
      console.log('Kafka producer connected successfully');
    } catch (error) {
      console.error('Failed to connect to Kafka:', error);
    }
  }

  private async createTopic() {
    try {
      await this.admin.connect();
      const existingTopics = await this.admin.listTopics();
      
      if (!existingTopics.includes(this.TOPIC)) {
        await this.admin.createTopics({
          topics: [{
            topic: this.TOPIC,
            numPartitions: 1,
            replicationFactor: 1
          }]
        });
        console.log(`Topic ${this.TOPIC} created successfully`);
      }
      await this.admin.disconnect();
    } catch (error) {
      console.error('Failed to create topic:', error);
    }
  }

  async log(message: string) {
    try {
      await this.producer.send({
        topic: this.TOPIC,
        messages: [{ value: `LOG: ${message}` }],
      });
      console.log(`Logged: ${message}`);
    } catch (error) {
      console.error('Failed to send log message:', error);
    }
  }

  async error(message: string, trace?: string) {
    try {
      await this.producer.send({
        topic: this.TOPIC,
        messages: [{ value: `ERROR: ${message}, TRACE: ${trace}` }],
      });
      console.error(`Error: ${message}`);
    } catch (error) {
      console.error('Failed to send error message:', error);
    }
  }

  async warn(message: string) {
    try {
      await this.producer.send({
        topic: this.TOPIC,
        messages: [{ value: `WARN: ${message}` }],
      });
      console.warn(`Warn: ${message}`);
    } catch (error) {
      console.error('Failed to send warn message:', error);
    }
  }
}
