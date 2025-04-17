// src/main.ts
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { KafkaLoggerService } from './logger/kafka-logger.service';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const logger = app.get(KafkaLoggerService);
  app.useLogger(logger);
  await app.listen(3000);
}
bootstrap();
