"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthModule = void 0;
const common_1 = require("@nestjs/common");
const auth_service_1 = require("./auth.service");
const auth_controller_1 = require("./auth.controller");
const config_module_1 = require("../config/config.module");
const nest_keycloak_connect_1 = require("nest-keycloak-connect");
const keycloak_config_service_1 = require("./strategies/keycloak/keycloak-config.service");
const keycloak_controller_1 = require("./keycloak.controller");
let AuthModule = class AuthModule {
};
exports.AuthModule = AuthModule;
exports.AuthModule = AuthModule = __decorate([
    (0, common_1.Module)({
        imports: [
            config_module_1.ConfigModule,
            nest_keycloak_connect_1.KeycloakConnectModule.registerAsync({
                useClass: keycloak_config_service_1.KeycloakConfigService,
                imports: [config_module_1.ConfigModule],
            })
        ],
        providers: [auth_service_1.AuthService, keycloak_config_service_1.KeycloakConfigService],
        controllers: [auth_controller_1.AuthController, keycloak_controller_1.KeycloakController],
        exports: [keycloak_config_service_1.KeycloakConfigService]
    })
], AuthModule);
//# sourceMappingURL=auth.module.js.map