import { Controller, Logger } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { Payload } from 'src/interfaces/payload.interface';
import { ValidatePayload } from 'src/utilites/payload/validate';
import { TwofaService } from './twofa.service';

/**
 * Grpc entrypoint for 2-factor authentication operations.
 */
@Controller()
export class TwofaController {
  private logger = new Logger('TwofaService');

  constructor(private readonly twofaService: TwofaService) {}

  /**
   * Send a verification email to the user with 6 digit code.
   *
   * @param data Payload
   */
  @GrpcMethod('TwoFAService', 'GenerateTwoFACode')
  async generateTwoFACode(data: Payload): Promise<void> {
    ValidatePayload(data);
    await this.twofaService.sendVerifactionEmail(data);
    this.logger.log(
      `Verification Email sent and Auth Code cached successfully`,
    );
  }
}
