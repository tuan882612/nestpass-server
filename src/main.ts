import { NestFactory } from '@nestjs/core';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { AppModule } from './app.module';
import ValidateEnv from './utilites/config/env.validator';

async function bootstrap() {
  if (!ValidateEnv()) {
    return;
  }

  try {
    const app = await NestFactory.createMicroservice<MicroserviceOptions>(
      AppModule,
      {
        transport: Transport.GRPC,
        options: {
          url: process.env.HOST + ':' + process.env.PORT,
          package: ['twofa', 'ping'],
          protoPath: [
            join(__dirname, '../src/twofa/twofa.proto'),
            join(__dirname, '../src/ping/ping.proto'),
          ],
          channelOptions: {
            'grpc.keepalive_time_ms': 1800000, // 30 minutes in milliseconds
            'grpc.keepalive_timeout_ms': 5000, // Timeout after waiting 5 seconds for a response
            'grpc.keepalive_permit_without_calls': 1, // 1 = true, allows pinging without active calls
            'grpc.http2.min_time_between_pings_ms': 1800000, // 30 minutes in milliseconds
            'grpc.http2.max_pings_without_data': 0, // 0 = unlimited pings when there's no data/stream
            'grpc.http2.min_ping_interval_without_data_ms': 300000,
          },
        },
      },
    );
    app.listen();
  } catch (error) {
    process.exit(1);
  }
}
bootstrap();
