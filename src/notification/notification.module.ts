import { Module } from '@nestjs/common';
import { NotificationService } from './notification.service';
import { NotificationController } from './notification.controller';
import { EmailService } from 'src/email/email.service';

@Module({
  imports: [],
  controllers: [NotificationController],
  providers: [NotificationService, EmailService],
})
export class NotificationModule {}
