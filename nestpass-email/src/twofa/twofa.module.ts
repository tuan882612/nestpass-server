import { Module } from '@nestjs/common';
import { EmailService } from 'src/email/email.service';
import { TwofaController } from './twofa.controller';
import { TwofaService } from './twofa.service';

@Module({
  imports: [],
  controllers: [TwofaController],
  providers: [TwofaService, EmailService],
})
export class TwofaModule {}
