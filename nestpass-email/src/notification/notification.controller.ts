import { Controller, Logger } from '@nestjs/common';
import { NotificationService } from './notification.service';

@Controller()
export class NotificationController {
  private loggger = new Logger('EmailController');

  constructor(private readonly notificationService: NotificationService) {}
}
