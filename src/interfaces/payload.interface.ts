export interface Payload {
  userId: string;
  email: string;
}

export interface PayloadDTO {
  user_id: string;
  email: string;
}

export interface CachePayload {
  Code: string;
  Retries: number;
}
