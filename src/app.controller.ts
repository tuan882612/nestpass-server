import { Controller } from '@nestjs/common';
import { AppService } from './app.service';
import { GrpcMethod } from '@nestjs/microservices';
import { HelloRequest, HelloResponse } from './app.interface';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @GrpcMethod('AppService', 'GetHello')
  getHello(data: HelloRequest): HelloResponse {
    return { message: this.appService.getHello(data.name) };
  }
}
