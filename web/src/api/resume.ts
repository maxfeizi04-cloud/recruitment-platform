import client from './client';

export interface Resume {
  id: string;
  user_id: string;
  title: string;
  content: string;
  is_default: boolean;
  attachment_urls: string[];
  created_at: string;
  updated_at: string;
}

export async function listResumes(): Promise<Resume[]> {
  const { data } = await client.get('/resumes');
  return data.resumes;
}

export async function createResume(title: string, content: string): Promise<Resume> {
  const { data } = await client.post('/resumes', { title, content });
  return data;
}

export async function updateResume(id: string, title: string, content: string): Promise<void> {
  await client.put(`/resumes/${id}`, { title, content });
}

export async function deleteResume(id: string): Promise<void> {
  await client.delete(`/resumes/${id}`);
}

export async function setDefaultResume(id: string): Promise<void> {
  await client.post(`/resumes/${id}/set-default`);
}

export async function uploadAttachment(id: string, file: File): Promise<string> {
  const form = new FormData();
  form.append('file', file);
  const { data } = await client.post(`/resumes/${id}/upload-attachment`, form);
  return data.url;
}
