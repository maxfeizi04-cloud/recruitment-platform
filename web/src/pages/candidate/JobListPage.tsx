import React, { useEffect, useState } from 'react';
import { List, Input, Card, Tag, Button, Typography, Space, Pagination, InputNumber, Select } from 'antd';
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
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);

  // Filter state
  const [salaryMin, setSalaryMin] = useState<number | null>(null);
  const [salaryMax, setSalaryMax] = useState<number | null>(null);
  const [sortBy, setSortBy] = useState<string>('latest');

  const navigate = useNavigate();

  const fetchJobs = async (q?: string, c?: string, p?: number, ps?: number) => {
    setLoading(true);
    const limit = ps ?? pageSize;
    const offset = ((p ?? page) - 1) * limit;
    try {
      const result = q ? await searchJobs(q, c || undefined, limit, offset) : await listJobs(limit, offset);
      setJobs(result.jobs);
      setTotal(result.total);
    } finally { setLoading(false); }
  };

  useEffect(() => { fetchJobs(); }, []);

  const handleSearch = () => {
    setPage(1);
    fetchJobs(query, city, 1, pageSize);
  };

  const handlePageChange = (p: number) => {
    setPage(p);
    fetchJobs(query, city, p, pageSize);
  };

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

  // Client-side filter + sort logic
  const filtered = jobs.filter(j => {
    try {
      const s = JSON.parse(j.salary_range);
      if (salaryMin != null && s.min < salaryMin * 1000) return false;
      if (salaryMax != null && s.max > salaryMax * 1000) return false;
    } catch { /* ignore parse errors */ }
    return true;
  }).sort((a, b) => {
    if (sortBy === 'salary_high') {
      const sa = JSON.parse(a.salary_range).max || 0;
      const sb = JSON.parse(b.salary_range).max || 0;
      return sb - sa;
    }
    if (sortBy === 'salary_low') {
      const sa = JSON.parse(a.salary_range).min || 0;
      const sb = JSON.parse(b.salary_range).min || 0;
      return sa - sb;
    }
    return 0; // latest = keep original order
  });

  return (
    <div>
      <Title level={4}>职位搜索</Title>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search placeholder="搜索职位" value={query} onChange={e => setQuery(e.target.value)} onSearch={handleSearch} style={{ width: 300 }} />
        <Input placeholder="城市" value={city} onChange={e => setCity(e.target.value)} onPressEnter={handleSearch} style={{ width: 150 }} />
      </Space>

      <Space wrap style={{ marginBottom: 12 }}>
        <InputNumber placeholder="最低薪资(K)" min={0} max={100} value={salaryMin} onChange={setSalaryMin} style={{ width: 120 }} />
        <span>-</span>
        <InputNumber placeholder="最高薪资(K)" min={0} max={100} value={salaryMax} onChange={setSalaryMax} style={{ width: 120 }} />
        <Select value={sortBy} onChange={setSortBy} style={{ width: 130 }}
          options={[
            { value: 'latest', label: '最新发布' },
            { value: 'salary_high', label: '薪资最高' },
            { value: 'salary_low', label: '薪资最低' },
          ]}
        />
      </Space>

      <Text type="secondary">共 {total} 个职位</Text>
      <List loading={loading} dataSource={filtered}
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
      <Pagination
        current={page}
        pageSize={pageSize}
        total={total}
        onChange={handlePageChange}
        style={{ marginTop: 16, textAlign: 'center' }}
        showSizeChanger={false}
      />
    </div>
  );
}
