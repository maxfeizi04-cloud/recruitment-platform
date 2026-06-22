import client from './client';

export interface Interview {
  id: string;
  job_application_id: string;
  hr_user_id: string;
  candidate_user_id: string;
  scheduled_at: string;
  company_address: string;
  contact_name: string;
  contact_phone: string;
  notes: string;
  status: string;
  job_title?: string;
  candidate_name?: string;
  hr_name?: string;
  created_at: string;
  updated_at: string;
}

export async function createInterview(inv: {
  application_id: string;
  scheduled_at: string;
  company_address: string;
  contact_name?: string;
  contact_phone?: string;
  notes?: string;
}): Promise<Interview> {
  const { data } = await client.post('/interviews', inv);
  return data;
}

export async function listMyInterviews(): Promise<Interview[]> {
  const { data } = await client.get('/interviews/my');
  return data.invitations;
}

export async function updateInterviewStatus(id: string, status: string): Promise<void> {
  await client.patch(`/interviews/${id}/status`, { status });
}

export async function getNavigateURL(id: string): Promise<string> {
  const { data } = await client.get(`/interviews/${id}/navigate`);
  return data.navigation_url;
}

export async function placeSearch(q: string): Promise<any[]> {
  const { data } = await client.get('/maps/place-search', { params: { q } });
  return data.suggestions;
}
