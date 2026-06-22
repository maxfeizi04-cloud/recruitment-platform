import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Typography, Tag, Button, Space, Modal, Select, message, Divider } from 'antd';
import { EnvironmentOutlined, DollarOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { getJob, type Job } from '../../api/job';
import { listResumes, type Resume } from '../../api/resume';
import { apply } from '../../api/application';
import { useAuth } from '../../stores/auth';

const { Title, Text } = Typography;

export default function JobDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [job, setJob] = useState<Job | null>(null);
  const [resumes, setResumes] = useState<Resume[]>([]);
  const [applyOpen, setApplyOpen] = useState(false);
  const [selectedResume, setSelectedResume] = useState<string>('');
  const [applying, setApplying] = useState(false);
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (id) getJob(id).then(setJob);
    if (auth.isAuthenticated && !auth.isHR) listResumes().then(setResumes);
  }, [id]);

  const handleApply = async () => {
    if (!selectedResume) { message.error('请选择简历'); return; }
    setApplying(true);
    try {
      await apply(id!, selectedResume);
      message.success('投递成功');
      setApplyOpen(false);
    } finally { setApplying(false); }
  };

  if (!job) return <div>加载中...</div>;

  const parseJSON = (s: string) => { try { return JSON.parse(s); } catch { return s; } };

  return (
    <div>
      <Button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>← 返回</Button>
      <Card>
        <Title level={3}>{job.title}</Title>
        <Space style={{ marginBottom: 16 }}>
          <Tag color="blue"><DollarOutlined /> {parseJSON(job.salary_range).min / 1000}k - {parseJSON(job.salary_range).max / 1000}k</Tag>
          <Tag icon={<EnvironmentOutlined />}>{parseJSON(job.location).city || parseJSON(job.location).province}</Tag>
          <Tag icon={<ClockCircleOutlined />}>{job.created_at?.slice(0, 10)}</Tag>
        </Space>
        <Divider />
        <Title level={5}>职位描述</Title>
        <Text>{job.description || '暂无描述'}</Text>
        {job.requirements && (
          <>
            <Divider />
            <Title level={5}>技能要求</Title>
            {Array.isArray(parseJSON(job.requirements).skills) && parseJSON(job.requirements).skills.map((s: string) => <Tag key={s}>{s}</Tag>)}
          </>
        )}
        {!auth.isHR && (
          <div style={{ marginTop: 24 }}>
            <Button type="primary" size="large" onClick={() => setApplyOpen(true)}>立即投递</Button>
          </div>
        )}
      </Card>
      <Modal open={applyOpen} title="选择简历投递" onOk={handleApply} onCancel={() => setApplyOpen(false)} confirmLoading={applying}>
        <Select style={{ width: '100%' }} placeholder="选择简历" value={selectedResume} onChange={setSelectedResume}>
          {resumes.map(r => <Select.Option key={r.id} value={r.id}>{r.title}</Select.Option>)}
        </Select>
      </Modal>
    </div>
  );
}
