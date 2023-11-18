import { Module } from '@nestjs/common';
import { NotificationModule } from './notification/notification.module';
import { TwofaModule } from './twofa/twofa.module';
import { PingModule } from './ping/ping.module';

@Module({
  imports: [TwofaModule, NotificationModule, PingModule],
  providers: [],
})
export class AppModule {}
