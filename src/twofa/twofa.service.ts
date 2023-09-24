import { status } from '@grpc/grpc-js';
import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { RpcException } from '@nestjs/microservices';
import { RedisClientType, createClient } from 'redis';
import { EmailService } from 'src/email/email.service';
import Email, { EmailType } from 'src/interfaces/email.interface';
import { CachePayload, Payload } from 'src/interfaces/payload.interface';

/**
 * Service for handling 2-factor authentication operations.
 */
@Injectable()
export class TwofaService implements OnModuleInit {
  private redisClient: RedisClientType;
  private logger = new Logger('TwofaService');

  constructor(private readonly emailService: EmailService) {}

  // initialize redis client after the service is initialized.
  onModuleInit() {
    try {
      this.redisClient = createClient({ url: process.env.REDIS_URL });
      this.redisClient.connect();
    } catch (error) {
      this.logger.error(error);
      throw new RpcException({ details: error.message, code: status.INTERNAL });
    }
  }

  /**
   * Sends a verification email to the user.
   * - also caches the verification code in redis.
   *
   * @param payload Payload of the user.
   * @returns Promise<void>
   */
  public async sendVerifactionEmail(payload: Payload): Promise<void> {
    const authCode: string = this.generateVerifactionCode();
    const email: Email = {
      to: payload.email,
      subject: 'nestpass - Auth Code: ' + authCode,
      template: this.emailService.getTemplate(EmailType.TWOFA),
    };

    // send email and cache the verification code asynchronously
    const cachePromise = this.cacheVerifactionCode(payload, authCode);
    const emailPromise = this.emailService.sendEmail(email);
    await Promise.all([cachePromise, emailPromise]);
  }

  /**
   * caches the verification code in redis.
   * - sets the expiry time to 3 minutes.
   * - transforms the payload to UserPayloadDTO.
   *
   * @param payload Payload of the user.
   * @param code Verification code.
   * @returns Promise<void>
   */
  public async cacheVerifactionCode(
    payload: Payload,
    code: string,
  ): Promise<void> {
    // check if redis client is connected.
    if (!this.redisClient || !this.redisClient.connect) {
      const errMsg: string = 'Redis client is not connected';
      this.logger.error(errMsg);
      throw new RpcException({ details: errMsg, code: status.INTERNAL });
    }

    const key: string = payload.userId;
    const value: CachePayload = {
      Code: code,
      Retries: 5,
    };

    // cache the verification code and set the expiry time.
    try {
      await this.redisClient.set(key, JSON.stringify(value));
      await this.redisClient.expire(key, 180);
    } catch (error) {
      this.logger.error(error);
      throw new RpcException({ details: error.message, code: status.INTERNAL });
    }
  }

  private generateVerifactionCode(): string {
    const min = 100000;
    const max = 999999;
    const code: number = Math.floor(Math.random() * (max - min + 1) + min);
    return code.toString();
  }
}
