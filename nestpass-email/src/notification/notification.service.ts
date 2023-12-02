import { Injectable } from '@nestjs/common';
import { EmailService } from 'src/email/email.service';

@Injectable()
export class NotificationService {
  constructor(private readonly emailService: EmailService) {}
}
