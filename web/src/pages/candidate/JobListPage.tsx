import React, { useEffect, useState } from 'react';
import { List, Input, Card, Tag, Button, Typography, Space } from 'antd';
import { SearchOutlined, EnvironmentOutlined, DollarOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { listJobs, searchJobs, type Job } from '../../api/job';

const { Text, Title } = Typography;

export default function JobListPage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState('');
  const [city, setCity] = useState('');
  const navigate = useNavigate();

  const fetchJobs = async (q?: string, c?: string) => {
    setLoading(true);
    try {
      const result = q ? await searchJobs(q, c || undefined) : await listJobs();
      setJobs(result.jobs);
      setTotal(result.total);
    } finally { setLoading(false); }
  };

  useEffect(() => { fetchJobs(); }, []);

  const parseSalary = (s: string) => {
    try {
      const r = JSON.parse(s);
      return `${r.min / 1000}k-${r.max / 1000}k`;
    } catch { return s; }
  };

  const parseLocation = (s: string) => {
    try { return JSON.parse(s).city || JSON.parse(s).province || ''; } catch { return s; }
  };

  const parseStatus = (s: string) => {
    const map: Record<string, { color: string; text: string }> = {
      active: { color: 'green', text: '招聘中' },
      paused: { color: 'orange', text: '暂停' },
      closed: { color: 'red', text: '关闭' },
    };
    return map[s] || { color: 'default', text: s };
  };

  return (
    <div>
      <Title level={4}>职位搜索</Title>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search placeholder="搜索职位" value={query} onChange={e => setQuery(e.target.value)} onSearch={() => fetchJobs(query, city)} style={{ width: 300 }} />
        <Input placeholder="城市" value={city} onChange={e => setCity(e.target.value)} onPressEnter={() => fetchJobs(query, city)} style={{ width: 150 }} />
      </Space>
      <Text type="secondary">共 {total} 个职位</Text>
      <List loading={loading} dataSource={jobs}
        renderItem={job => (
          <Card hoverable style={{ marginTop: 12 }} onClick={() => navigate(`/jobs/${job.id}`)}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong style={{ fontSize: 16 }}>{job.title}</Text>
                <div style={{ marginTop: 8 }}>
                  <Tag color="blue"><DollarOutlined /> {parseSalary(job.salary_range)}</Tag>
                  <Tag icon={<EnvironmentOutlined />}>{parseLocation(job.location)}</Tag>
                  <Tag color={parseStatus(job.status).color}>{parseStatus(job.status).text}</Tag>
                </div>
              </div>
              <Button type="primary">查看详情</Button>
            </div>
          </Card>
        )}
        locale={{ emptyText: '暂无职位' }}
      />
    </div>
  );
}
