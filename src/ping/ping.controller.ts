import { Controller } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { PingData } from 'src/interfaces/ping.interface';
import { PingService } from './ping.service';

/**
 * GRPC controller for ping
 */
@Controller()
export class PingController {
  /**
   * Injects the ping service
   * @param {PingService} pingService
   */
  constructor(private readonly pingService: PingService) {}

  /**
   * Returns a string message
   * @returns string
   */
  @GrpcMethod('PingService', 'Ping')
  public Ping(data: PingData): PingData {
    return this.pingService.Ping(data.message);
  }
}
