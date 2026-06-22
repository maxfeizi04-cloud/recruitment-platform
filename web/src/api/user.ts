import client from './client';
import type { User } from './auth';

export async function getProfile(): Promise<User> {
  const { data } = await client.get('/users/me');
  return data;
}

export async function updateProfile(name: string, avatar_url?: string): Promise<User> {
  const { data } = await client.put('/users/me', { name, avatar_url });
  return data;
}

export async function submitCertification(company_name: string, position: string): Promise<void> {
  await client.post('/hrs/certify', { company_name, position });
}

export async function getCertification(): Promise<any> {
  const { data } = await client.get('/hrs/certification');
  return data;
}
