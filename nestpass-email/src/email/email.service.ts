import { status } from '@grpc/grpc-js';
import { Injectable, Logger } from '@nestjs/common';
import { RpcException } from '@nestjs/microservices';
import * as sendgrid from '@sendgrid/mail';
import Email from 'src/interfaces/email.interface';
import { EmailType } from '../interfaces/email.interface';

/**
 * Email Service for sending emails using SendGrid
 */
@Injectable()
export class EmailService {
  private logger = new Logger('EmailService');
  private sender: string = process.env.MAIL_SENDER;
  constructor() {
    sendgrid.setApiKey(process.env.MAIL_API_KEY);
  }

  /**
   * raw email sending method
   *
   * @param email Email object
   */
  async sendEmail(email: Email): Promise<void> {
    try {
      await sendgrid.send({
        to: email.to,
        from: this.sender,
        subject: email.subject,
        html: email.template,
      });
    } catch (error) {
      this.logger.error(error);
      throw new RpcException({ details: error.message, code: status.INTERNAL });
    }
  }

  /**
   * Returns the email template based on the email type
   *
   * @param emailType EmailType
   * @returns string
   * - email template based on the email type
   */
  getTemplate(emailType: EmailType): string {
    switch (emailType) {
      case EmailType.TWOFA:
        return `
        <table cellpadding="0" cellspacing="0" style="vertical-align: -webkit-baseline-middle; font-size: medium; font-family: Arial;">
          <tbody>
            <span>Your auth code expires in <b>3 min</b></span>
            <tr>
              <td width="190">
                <img src="">
              </td>
              <td>
                <table cellpadding="20" cellspacing="0" style="vertical-align: -webkit-baseline-middle; font-size: medium; font-family: Arial; width: 100%;">
                  <tbody>
                    <tr>
                      <td>
                        <p style="margin: 0px; font-size: 15px; font-weight:bold; color: #111; line-height: 20px;">
                          <span>MAIL Service</span>
                        </p>
                        <p style="margin: 0px; color: #687087; font-size: 14px; line-height: 20px;">
                          <span>domain nestpass.tech</span>
                        </p>
                        <p style="margin: 0px; color: #687087; font-size: 14px; line-height: 20px;"></p>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </td>
            </tr>
          </tbody>
        </table>
        `;
      default:
        return '';
    }
  }
}
