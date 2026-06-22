import client from './client';

export interface User {
  id: string;
  phone: string;
  role: 'candidate' | 'hr';
  name: string;
  avatar_url: string;
}

export interface LoginResult {
  token: string;
  user: User;
}

export async function sendCode(phone: string): Promise<void> {
  await client.post('/auth/send-code', { phone });
}

export async function login(phone: string, code: string, role: string): Promise<LoginResult> {
  const { data } = await client.post('/auth/login', { phone, code, role });
  return data;
}
