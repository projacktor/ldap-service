// src/logger/logger.module.ts
import { Module } from '@nestjs/common';
import { KafkaLoggerService } from './kafka-logger.service';

@Module({
  providers: [KafkaLoggerService],
  exports: [KafkaLoggerService],
})
export class LoggerModule {}
