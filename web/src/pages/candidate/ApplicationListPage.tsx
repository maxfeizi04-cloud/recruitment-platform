import React, { useEffect, useState } from 'react';
import { Table, Tag, Typography } from 'antd';
import { listMyApplications, type Application } from '../../api/application';

const { Title } = Typography;

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'blue', text: '待查看' },
  viewed: { color: 'orange', text: '已查看' },
  accepted: { color: 'green', text: '已通过' },
  rejected: { color: 'red', text: '未通过' },
};

export default function ApplicationListPage() {
  const [apps, setApps] = useState<Application[]>([]);

  useEffect(() => { listMyApplications().then(setApps); }, []);

  return (
    <div>
      <Title level={4}>投递记录</Title>
      <Table dataSource={apps} rowKey="id" pagination={false}
        columns={[
          { title: '职位', dataIndex: 'job_title', key: 'job' },
          { title: '状态', dataIndex: 'status', key: 'status', render: (s: string) => {
            const m = statusMap[s] || { color: 'default', text: s };
            return <Tag color={m.color}>{m.text}</Tag>;
          }},
          { title: '投递时间', dataIndex: 'created_at', key: 'time', render: (t: string) => t?.slice(0, 10) },
        ]}
        locale={{ emptyText: '暂无投递记录' }}
      />
    </div>
  );
}
