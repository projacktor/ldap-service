import { Module } from '@nestjs/common';
import { ConfigModule as NestConfigModule } from '@nestjs/config';
import { KeycloakConfig } from './keycloak.config';

@Module({
  imports: [NestConfigModule.forRoot()],
  providers: [KeycloakConfig],
  exports: [KeycloakConfig],
})
export class ConfigModule {}
