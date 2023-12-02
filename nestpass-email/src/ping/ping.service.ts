import { Injectable, Logger } from '@nestjs/common';
import { PingData } from 'src/interfaces/ping.interface';

/**
 * Used to ping the server
 */
@Injectable()
export class PingService {
  private logger = new Logger('PingService');

  /**
   * Returns a string message
   * @returns ping message data
   */
  public Ping(data: string): PingData {
    const msg = 'Origin: ' + data + ' - successfull ping';
    this.logger.log(msg);
    return { message: msg };
  }
}
