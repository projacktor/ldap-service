import { Injectable } from '@nestjs/common';

@Injectable()
export class AuthService {
  uses(): void {
  console.log("Hi")
  }
}
