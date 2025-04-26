import { KeycloakConnectOptions, KeycloakConnectOptionsFactory } from 'nest-keycloak-connect';
import { KeycloakConfig } from '../../../config/keycloak.config';
export declare class KeycloakConfigService implements KeycloakConnectOptionsFactory {
    private readonly keycloakConfig;
    constructor(keycloakConfig: KeycloakConfig);
    createKeycloakConnectOptions(): KeycloakConnectOptions;
}
