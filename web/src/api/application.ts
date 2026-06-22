import client from './client';

export interface Application {
  id: string;
  job_id: string;
  user_id: string;
  resume_id: string;
  status: string;
  job_title?: string;
  user_name?: string;
  created_at: string;
  updated_at: string;
}

export async function apply(job_id: string, resume_id: string): Promise<Application> {
  const { data } = await client.post('/applications', { job_id, resume_id });
  return data;
}

export async function listMyApplications(): Promise<Application[]> {
  const { data } = await client.get('/applications/my');
  return data.applications;
}

export async function listReceivedApplications(): Promise<Application[]> {
  const { data } = await client.get('/applications/received');
  return data.applications;
}

export async function updateApplicationStatus(id: string, status: string): Promise<void> {
  await client.patch(`/applications/${id}/status`, { status });
}
