import { Injectable } from '@nestjs/common';
import { KeycloakConnectOptions, KeycloakConnectOptionsFactory } from 'nest-keycloak-connect';
import { KeycloakConfig } from '../../../config/keycloak.config';

@Injectable()
export class KeycloakConfigService implements KeycloakConnectOptionsFactory {
  constructor(private readonly keycloakConfig: KeycloakConfig) {}

  createKeycloakConnectOptions(): KeycloakConnectOptions {
    return {
      authServerUrl: this.keycloakConfig.authServerUrl,
      realm: this.keycloakConfig.realm,
      clientId: this.keycloakConfig.clientId,
      secret: this.keycloakConfig.secret
    }
  }
}