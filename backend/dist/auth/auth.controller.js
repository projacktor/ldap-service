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
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthController = void 0;
const common_1 = require("@nestjs/common");
const nest_keycloak_connect_1 = require("nest-keycloak-connect");
let AuthController = class AuthController {
    getPublicData() {
        return { message: 'Public allowed' };
    }
    getProtectedData(req) {
        console.log('User roles: ', req);
        return { message: 'Protected allowed' };
    }
    getAdminData() {
        return { message: 'Admin allowed' };
    }
};
exports.AuthController = AuthController;
__decorate([
    (0, common_1.Get)('/public'),
    (0, nest_keycloak_connect_1.Public)(),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", void 0)
], AuthController.prototype, "getPublicData", null);
__decorate([
    (0, common_1.Get)('/protected'),
    (0, nest_keycloak_connect_1.Roles)({ roles: ['user'] }),
    __param(0, (0, common_1.Req)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Request]),
    __metadata("design:returntype", void 0)
], AuthController.prototype, "getProtectedData", null);
__decorate([
    (0, common_1.Get)('/admin'),
    (0, nest_keycloak_connect_1.Roles)({ roles: ['admin'] }),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", void 0)
], AuthController.prototype, "getAdminData", null);
exports.AuthController = AuthController = __decorate([
    (0, common_1.Controller)('auth')
], AuthController);
//# sourceMappingURL=auth.controller.js.map