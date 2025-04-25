import { Controller, Get, Logger } from '@nestjs/common';
import { KafkaLoggerService } from './kafka-logger.service';

@Controller('logger')
export class LoggerController {
  constructor(private readonly kafkaLogger: KafkaLoggerService) {}

  @Get('test')
  async testLogging() {
    await this.kafkaLogger.log('Hello World from Kafka Logger!');
    await this.kafkaLogger.warn('This is a warning message');
    await this.kafkaLogger.error('This is an error message', 'Test error trace');
    
    return {
      message: 'Logs have been sent to Kafka',
      timestamp: new Date().toISOString()
    };
  }
} 