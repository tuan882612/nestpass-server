export interface Payload {
  userId: string;
  email: string;
  userStatus: string;
}

export interface PayloadDTO {
  user_id: string;
  email: string;
  user_status: string;
}

export interface CachePayload {
  Code: string;
  Retries: number;
  UserStatus: string;
}
