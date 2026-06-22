import React, { useEffect, useState } from 'react';
import { Card, Typography, Tag, Button, Space, message } from 'antd';
import { EnvironmentOutlined, PhoneOutlined, CalendarOutlined } from '@ant-design/icons';
import { listMyInterviews, updateInterviewStatus, getNavigateURL, type Interview } from '../../api/interview';
import { useAuth } from '../../stores/auth';

const { Title, Text } = Typography;

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'blue', text: '待确认' },
  accepted: { color: 'green', text: '已接受' },
  declined: { color: 'red', text: '已婉拒' },
  reschedule: { color: 'orange', text: '改期中' },
  confirmed: { color: 'purple', text: '已确认' },
};

export default function InterviewListPage() {
  const [interviews, setInterviews] = useState<Interview[]>([]);
  const { isHR } = useAuth();

  const fetch = async () => { setInterviews(await listMyInterviews()); };
  useEffect(() => { fetch(); }, []);

  const handleStatus = async (id: string, status: string) => {
    await updateInterviewStatus(id, status);
    message.success('状态已更新');
    fetch();
  };

  const handleNavigate = async (id: string) => {
    const url = await getNavigateURL(id);
    window.open(url, '_blank');
  };

  const parseAddress = (s: string) => { try { return JSON.parse(s); } catch { return {}; } };

  return (
    <div>
      <Title level={4}>面试邀约</Title>
      {interviews.map(inv => {
        const addr = parseAddress(inv.company_address);
        const st = statusMap[inv.status] || { color: 'default', text: inv.status };
        return (
          <Card key={inv.id} style={{ marginBottom: 12 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <div>
                <Text strong style={{ fontSize: 16 }}>{inv.job_title}</Text>
                <Tag color={st.color} style={{ marginLeft: 8 }}>{st.text}</Tag>
              </div>
            </div>
            <div style={{ marginTop: 12 }}>
              <Tag icon={<CalendarOutlined />}>{inv.scheduled_at?.slice(0, 16)}</Tag>
              <Tag icon={<EnvironmentOutlined />}>{addr.formatted || addr.detail || addr.city}</Tag>
              {inv.contact_phone && <Tag icon={<PhoneOutlined />}>{inv.contact_name}: {inv.contact_phone}</Tag>}
            </div>
            {inv.notes && <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>备注：{inv.notes}</Text>}
            <Space style={{ marginTop: 12 }}>
              {inv.status === 'pending' && (
                <>
                  <Button type="primary" size="small" onClick={() => handleStatus(inv.id, 'accepted')}>接受</Button>
                  <Button size="small" onClick={() => handleStatus(inv.id, 'declined')}>婉拒</Button>
                </>
              )}
              {(inv.status === 'accepted' || inv.status === 'confirmed') && (
                <Button type="primary" size="small" onClick={() => handleNavigate(inv.id)}><EnvironmentOutlined /> 导航到公司</Button>
              )}
            </Space>
          </Card>
        );
      })}
      {interviews.length === 0 && <Text type="secondary">暂无面试邀约</Text>}
    </div>
  );
}
