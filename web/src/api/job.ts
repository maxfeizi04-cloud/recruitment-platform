import client from './client';

export interface Job {
  id: string;
  hr_user_id: string;
  title: string;
  description: string;
  requirements: string;
  salary_range: string;
  location: string;
  status: string;
  expire_at: string;
  created_at: string;
  updated_at: string;
}

export interface JobListResult {
  jobs: Job[];
  total: number;
}

export async function listJobs(limit = 20, offset = 0): Promise<JobListResult> {
  const { data } = await client.get('/jobs', { params: { limit, offset } });
  return data;
}

export async function searchJobs(q: string, city?: string, limit = 20, offset = 0): Promise<JobListResult> {
  const { data } = await client.get('/jobs/search', { params: { q, city, limit, offset } });
  return data;
}

export async function getJob(id: string): Promise<Job> {
  const { data } = await client.get(`/jobs/${id}`);
  return data;
}

export async function createJob(job: { title: string; description?: string; requirements?: string; salary_range?: string; location?: string }): Promise<Job> {
  const { data } = await client.post('/jobs', job);
  return data;
}

export async function updateJob(id: string, job: { title: string; description?: string; requirements?: string; salary_range?: string; location?: string }): Promise<void> {
  await client.put(`/jobs/${id}`, job);
}

export async function updateJobStatus(id: string, status: string): Promise<void> {
  await client.patch(`/jobs/${id}/status`, { status });
}

export async function listMyJobs(): Promise<Job[]> {
  const { data } = await client.get('/jobs/my');
  return data.jobs;
}
