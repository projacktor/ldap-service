// src/logger/logger.module.ts
import { Module } from '@nestjs/common';
import { KafkaLoggerService } from './kafka-logger.service';
import { LoggerController } from './logger.controller';

@Module({
  controllers: [LoggerController],
  providers: [KafkaLoggerService],
  exports: [KafkaLoggerService],
})
export class LoggerModule {}
