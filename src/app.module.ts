import { Module } from '@nestjs/common';
import { NotificationModule } from './notification/notification.module';
import { TwofaModule } from './twofa/twofa.module';

@Module({
  imports: [TwofaModule, NotificationModule],
  providers: [],
})
export class AppModule {}
