import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    {
      transport: Transport.GRPC,
      options: {
        url: 'localhost:5000',
        package: 'app',
        protoPath: join(__dirname, './app.proto'),
      },
    },
  );
  await app.listen().then(() => console.log('listening on port 5000'));
}
bootstrap();
