import React, { useEffect, useState } from 'react';
import { Table, Button, Tag, Space, Popconfirm, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { listMyJobs, updateJobStatus, type Job } from '../../api/job';

const { Title } = Typography;

export default function JobManagePage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const navigate = useNavigate();

  useEffect(() => { listMyJobs().then(setJobs); }, []);

  const handleStatus = async (id: string, status: string) => {
    await updateJobStatus(id, status);
    message.success('状态已更新');
    listMyJobs().then(setJobs);
  };

  const statusMap: Record<string, { color: string; text: string }> = {
    active: { color: 'green', text: '招聘中' },
    paused: { color: 'orange', text: '已暂停' },
    closed: { color: 'red', text: '已关闭' },
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>职位管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/jobs/new/edit')}>发布新职位</Button>
      </div>
      <Table dataSource={jobs} rowKey="id" pagination={false}
        columns={[
          { title: '职位', dataIndex: 'title', key: 'title' },
          { title: '状态', dataIndex: 'status', key: 'status', render: (s: string) => {
            const m = statusMap[s] || { color: 'default', text: s };
            return <Tag color={m.color}>{m.text}</Tag>;
          }},
          { title: '发布时间', dataIndex: 'created_at', key: 'time', render: (t: string) => t?.slice(0, 10) },
          { title: '操作', key: 'actions', render: (_: any, record: Job) => (
            <Space>
              <Button size="small" icon={<EditOutlined />} onClick={() => navigate(`/jobs/${record.id}/edit`)}>编辑</Button>
              {record.status === 'active' && <Button size="small" onClick={() => handleStatus(record.id, 'paused')}>暂停</Button>}
              {record.status === 'paused' && <Button size="small" onClick={() => handleStatus(record.id, 'active')}>恢复</Button>}
              {record.status !== 'closed' && <Popconfirm title="确定关闭？" onConfirm={() => handleStatus(record.id, 'closed')}><Button size="small" danger>关闭</Button></Popconfirm>}
            </Space>
          )},
        ]}
      />
    </div>
  );
}
