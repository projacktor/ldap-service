"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.KeycloakConfigService = void 0;
const common_1 = require("@nestjs/common");
const keycloak_config_1 = require("../../../config/keycloak.config");
let KeycloakConfigService = class KeycloakConfigService {
    keycloakConfig;
    constructor(keycloakConfig) {
        this.keycloakConfig = keycloakConfig;
    }
    createKeycloakConnectOptions() {
        return {
            authServerUrl: this.keycloakConfig.authServerUrl,
            realm: this.keycloakConfig.realm,
            clientId: this.keycloakConfig.clientId,
            secret: this.keycloakConfig.secret
        };
    }
};
exports.KeycloakConfigService = KeycloakConfigService;
exports.KeycloakConfigService = KeycloakConfigService = __decorate([
    (0, common_1.Injectable)(),
    __metadata("design:paramtypes", [keycloak_config_1.KeycloakConfig])
], KeycloakConfigService);
//# sourceMappingURL=keycloak-config.service.js.map