import { Payload, PayloadDTO } from 'src/interfaces/payload.interface';
import { ValidatePayload } from './validate';

export function ToPayloadDTO(payload: Payload): PayloadDTO {
  ValidatePayload(payload);
  const mapped: PayloadDTO = {
    user_id: payload.userId,
    email: payload.email,
    user_status: payload.userStatus,
  };

  return mapped;
}
