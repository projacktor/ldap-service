import { ConfigService } from '@nestjs/config';
export declare class KeycloakConfig {
    private configService;
    constructor(configService: ConfigService);
    private getOrThrow;
    get authServerUrl(): string;
    get realm(): string;
    get clientId(): string;
    get secret(): string;
}
