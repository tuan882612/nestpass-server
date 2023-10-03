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
          package: ['twofa'],
          protoPath: join(__dirname, '../src/twofa/twofa.proto'),
        },
      },
    );
    app.listen();
  } catch (error) {
    process.exit(1);
  }
}
bootstrap();
