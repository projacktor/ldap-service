import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
// import { LoggerModule } from './logger/logger.module';
import { AuthModule } from './auth/auth.module';
import { ConfigModule } from '@nestjs/config';
import { APP_GUARD } from '@nestjs/core';
import { AuthGuard, KeycloakConnectModule, RoleGuard } from 'nest-keycloak-connect';
import { KeycloakConfigService } from './auth/strategies/keycloak/keycloak-config.service';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    AuthModule,
    KeycloakConnectModule.registerAsync({
      useExisting: KeycloakConfigService,
      imports: [AuthModule]
    })
  ],
  controllers: [AppController],
  providers: [AppService,
    {
      provide: APP_GUARD,
      useClass: AuthGuard
    },
    {
      provide: APP_GUARD,
      useClass: RoleGuard
    }
  ],
})

export class AppModule {}
