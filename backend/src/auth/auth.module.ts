import { Module } from '@nestjs/common';
import { AuthService } from './auth.service';
import { AuthController } from './auth.controller';
import { ConfigModule } from '../config/config.module';
import { KeycloakConnectModule } from 'nest-keycloak-connect';
import { KeycloakConfigService } from './strategies/keycloak/keycloak-config.service';
import { KeycloakController } from './keycloak.controller';

@Module({
  imports: [
    ConfigModule,
    KeycloakConnectModule.registerAsync({
      useClass: KeycloakConfigService,
      imports: [ConfigModule],
    })
  ],
  providers: [AuthService, KeycloakConfigService],
  controllers: [AuthController, KeycloakController],
  exports: [KeycloakConfigService]
})
export class AuthModule {}
