import { status } from '@grpc/grpc-js';
import { RpcException } from '@nestjs/microservices';
import { Payload } from 'src/interfaces/payload.interface';

export function ValidatePayload(payload: Payload) {
  const ref: string[] = ['userId', 'email'];

  const missingKeys = ref.filter((key) => {
    return !payload[key];
  });

  if (missingKeys.length > 0) {
    const errMsg = `Payload param is missing: ${missingKeys.join(', ')}`;
    console.error(errMsg);
    throw new RpcException({ details: errMsg, code: status.INVALID_ARGUMENT });
  }

  const emptyKeys = ref.filter((key) => {
    return payload[key] === '';
  });

  if (!emptyKeys || emptyKeys.includes(undefined)) {
    const errMsg = `Payload param is empty: ${emptyKeys.join(', ')}`;
    console.error(errMsg);
    throw new RpcException({ details: errMsg, code: status.INVALID_ARGUMENT });
  }
}
