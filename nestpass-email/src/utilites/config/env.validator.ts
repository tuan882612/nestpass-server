import * as dotenv from 'dotenv';
import { Config } from 'src/interfaces/config.interface';

/**
 * Validate environment variables
 *
 * @returns {boolean} true if all environment variables are set, false otherwise
 */
export default function ValidateEnv(): boolean {
  dotenv.config();
  const config: Config = {
    MAIL_API_KEY: process.env.MAIL_API_KEY,
    MAIL_SENDER: process.env.MAIL_SENDER,
    REDIS_URL: process.env.REDIS_URL,
    EMAIL_KEY: process.env.EMAIL_KEY,
    PORT: process.env.PORT,
    HOST: process.env.HOST,
  };

  const emptyKeys: string[] = Object.keys(config).filter(
    (key: string) => !Object.values(config)[Object.keys(config).indexOf(key)],
  );

  if (emptyKeys.length) {
    console.error(`Environment variables are missing: ${emptyKeys.join(', ')}`);
    return false;
  }

  return true;
}
