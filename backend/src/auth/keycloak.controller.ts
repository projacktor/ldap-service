import { Controller, Get } from '@nestjs/common';
import { KeycloakConfig } from '../config/keycloak.config';

@Controller('keycloak')
export class KeycloakController {
  constructor(private readonly keycloakConfig: KeycloakConfig) {}

  @Get('/check')
  checkConnection() {
    try {
      return {
        authServerUrl: this.keycloakConfig.authServerUrl,
        realm: this.keycloakConfig.realm,
        clientId: this.keycloakConfig.clientId,
        secret: this.keycloakConfig.secret,
      };
    } catch (err) {
      return { error: err.message };
    }
  }
}
