export default interface Email {
  to: string;
  subject: string;
  template: string;
}

export enum EmailType {
  TWOFA = '2fa',
  BASE = 'base',
}
