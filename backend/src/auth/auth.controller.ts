import { Controller, Get, Req } from '@nestjs/common';
import { Public, Roles } from 'nest-keycloak-connect';

@Controller('auth')
export class AuthController {
  @Get('/public')
  @Public()
  getPublicData() {
    return { message: 'Public allowed' };
  }

  @Get('/protected')
  @Roles({ roles: ['user']})
  getProtectedData(@Req() req: Request) {
    console.log('User roles: ', req)
    return { message: 'Protected allowed' };
  }

  @Get('/admin')
  @Roles({ roles: ['admin']})
  getAdminData() {
    return { message: 'Admin allowed' };
  }
}
