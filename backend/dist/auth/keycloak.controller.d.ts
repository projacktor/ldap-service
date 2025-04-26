import { KeycloakConfig } from '../config/keycloak.config';
export declare class KeycloakController {
    private readonly keycloakConfig;
    constructor(keycloakConfig: KeycloakConfig);
    checkConnection(): {
        authServerUrl: string;
        realm: string;
        clientId: string;
        secret: string;
        error?: undefined;
    } | {
        error: any;
        authServerUrl?: undefined;
        realm?: undefined;
        clientId?: undefined;
        secret?: undefined;
    };
}
