import React, { useEffect, useState } from 'react';
import { Table, Tag, Button, Space, Modal, Form, Input, DatePicker, message, Typography } from 'antd';
import { CheckOutlined, CloseOutlined, CalendarOutlined } from '@ant-design/icons';
import { listReceivedApplications, updateApplicationStatus, type Application } from '../../api/application';
import { createInterview } from '../../api/interview';

const { Title } = Typography;

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'blue', text: '待查看' },
  viewed: { color: 'orange', text: '已查看' },
  accepted: { color: 'green', text: '已通过' },
  rejected: { color: 'red', text: '未通过' },
};

export default function CandidateListPage() {
  const [apps, setApps] = useState<Application[]>([]);
  const [interviewOpen, setInterviewOpen] = useState(false);
  const [selectedApp, setSelectedApp] = useState<Application | null>(null);
  const [form] = Form.useForm();

  const fetch = async () => { setApps(await listReceivedApplications()); };
  useEffect(() => { fetch(); }, []);

  const handleStatus = async (id: string, status: string) => {
    await updateApplicationStatus(id, status);
    message.success('状态已更新');
    fetch();
  };

  const handleInterview = async (values: any) => {
    if (!selectedApp) return;
    const address = JSON.stringify({
      province: values.province, city: values.city, district: values.district,
      detail: values.detail, lat: 0, lng: 0,
      formatted: `${values.province}${values.city}${values.district}${values.detail}`,
    });
    await createInterview({
      application_id: selectedApp.id,
      scheduled_at: values.scheduled_at?.toISOString(),
      company_address: address,
      contact_name: values.contact_name,
      contact_phone: values.contact_phone,
      notes: values.notes,
    });
    message.success('面试邀约已发送');
    setInterviewOpen(false);
    fetch();
  };

  return (
    <div>
      <Title level={4}>候选人管理</Title>
      <Table dataSource={apps} rowKey="id" pagination={false}
        columns={[
          { title: '候选人', dataIndex: 'user_name', key: 'name' },
          { title: '职位', dataIndex: 'job_title', key: 'job' },
          { title: '状态', dataIndex: 'status', key: 'status', render: (s: string) => {
            const m = statusMap[s] || { color: 'default', text: s };
            return <Tag color={m.color}>{m.text}</Tag>;
          }},
          { title: '投递时间', dataIndex: 'created_at', key: 'time', render: (t: string) => t?.slice(0, 10) },
          { title: '操作', key: 'actions', render: (_: any, record: Application) => (
            <Space>
              {record.status === 'pending' && <Button size="small" onClick={() => handleStatus(record.id, 'viewed')}>标记已读</Button>}
              <Button size="small" type="primary" icon={<CheckOutlined />} onClick={() => handleStatus(record.id, 'accepted')}>通过</Button>
              <Button size="small" danger icon={<CloseOutlined />} onClick={() => handleStatus(record.id, 'rejected')}>拒绝</Button>
              <Button size="small" icon={<CalendarOutlined />} onClick={() => { setSelectedApp(record); setInterviewOpen(true); }}>邀约面试</Button>
            </Space>
          )},
        ]}
      />

      <Modal title="发起面试邀约" open={interviewOpen} onCancel={() => setInterviewOpen(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleInterview}>
          <Form.Item name="scheduled_at" label="面试时间" rules={[{ required: true }]}><DatePicker showTime style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="公司地址">
            <Input.Group compact>
              <Form.Item name="province" noStyle rules={[{ required: true }]}><Input placeholder="省" style={{ width: '25%' }} /></Form.Item>
              <Form.Item name="city" noStyle rules={[{ required: true }]}><Input placeholder="市" style={{ width: '25%' }} /></Form.Item>
              <Form.Item name="district" noStyle><Input placeholder="区" style={{ width: '25%' }} /></Form.Item>
              <Form.Item name="detail" noStyle rules={[{ required: true }]}><Input placeholder="详细地址" style={{ width: '25%' }} /></Form.Item>
            </Input.Group>
          </Form.Item>
          <Form.Item name="contact_name" label="联系人"><Input /></Form.Item>
          <Form.Item name="contact_phone" label="联系电话"><Input /></Form.Item>
          <Form.Item name="notes" label="备注"><Input.TextArea rows={2} /></Form.Item>
          <Form.Item><Button type="primary" htmlType="submit">发送邀约</Button></Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
