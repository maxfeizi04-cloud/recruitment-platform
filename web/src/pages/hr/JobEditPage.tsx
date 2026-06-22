import React, { useEffect, useState } from 'react';
import { Card, Form, Input, Button, Typography, message, InputNumber } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { createJob, getJob, updateJob } from '../../api/job';

const { Title } = Typography;

export default function JobEditPage() {
  const { id } = useParams<{ id: string }>();
  const isEdit = id && id !== 'new';
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  useEffect(() => {
    if (isEdit) {
      getJob(id).then(job => {
        try {
          const loc = JSON.parse(job.location);
          const sal = JSON.parse(job.salary_range);
          const req = JSON.parse(job.requirements);
          form.setFieldsValue({
            title: job.title, description: job.description,
            province: loc.province, city: loc.city, district: loc.district,
            salary_min: sal.min / 1000, salary_max: sal.max / 1000,
            skills: (req.skills || []).join(', '),
          });
        } catch { form.setFieldsValue(job); }
      });
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setLoading(true);
    const jobData = {
      title: values.title,
      description: values.description || '',
      requirements: JSON.stringify({ skills: (values.skills || '').split(',').map((s: string) => s.trim()).filter(Boolean) }),
      salary_range: JSON.stringify({ min: (values.salary_min || 0) * 1000, max: (values.salary_max || 0) * 1000, period: 'monthly' }),
      location: JSON.stringify({ province: values.province || '', city: values.city || '', district: values.district || '' }),
    };
    try {
      if (isEdit) {
        await updateJob(id, jobData);
        message.success('职位已更新');
      } else {
        await createJob(jobData);
        message.success('职位已发布');
      }
      navigate('/jobs/manage');
    } finally { setLoading(false); }
  };

  return (
    <div>
      <Title level={4}>{isEdit ? '编辑职位' : '发布新职位'}</Title>
      <Card>
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item name="title" label="职位名称" rules={[{ required: true }]}><Input placeholder="如：高级前端工程师" /></Form.Item>
          <Form.Item name="description" label="职位描述"><Input.TextArea rows={4} /></Form.Item>
          <Form.Item label="工作地点">
            <Input.Group compact>
              <Form.Item name="province" noStyle><Input placeholder="省" style={{ width: '30%' }} /></Form.Item>
              <Form.Item name="city" noStyle><Input placeholder="市" style={{ width: '35%' }} /></Form.Item>
              <Form.Item name="district" noStyle><Input placeholder="区" style={{ width: '35%' }} /></Form.Item>
            </Input.Group>
          </Form.Item>
          <Form.Item label="薪资范围（K/月）">
            <Input.Group compact>
              <Form.Item name="salary_min" noStyle><InputNumber placeholder="最低" min={0} style={{ width: '45%' }} /></Form.Item>
              <span style={{ lineHeight: '32px', margin: '0 8px' }}>-</span>
              <Form.Item name="salary_max" noStyle><InputNumber placeholder="最高" min={0} style={{ width: '45%' }} /></Form.Item>
            </Input.Group>
          </Form.Item>
          <Form.Item name="skills" label="技能要求（逗号分隔）"><Input placeholder="React, TypeScript, Node.js" /></Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>{isEdit ? '保存' : '发布'}</Button>
            <Button style={{ marginLeft: 8 }} onClick={() => navigate(-1)}>取消</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
