import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class KeycloakConfig {
  constructor(private configService: ConfigService) {}

  private getOrThrow(key: string): string {
    const value = this.configService.get<string>(key);
    if (!value) {
      throw new Error(`Configuration key "${key}" is missing or undefined`);
    }
    return value;
  }

  get authServerUrl(): string {
    return this.getOrThrow('KEYCLOAK_AUTH_SERVER_URL');
  }

  get realm(): string {
    return this.getOrThrow('KEYCLOAK_REALM');
  }

  get clientId(): string {
    return this.getOrThrow('KEYCLOAK_CLIENT_ID');
  }

  get secret(): string {
    return this.getOrThrow('KEYCLOAK_CLIENT_SECRET');
  }
}