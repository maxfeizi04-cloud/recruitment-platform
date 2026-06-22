import React, { useEffect, useState } from 'react';
import { Card, Form, Input, Button, Typography, message } from 'antd';
import { BankOutlined } from '@ant-design/icons';

const { Title } = Typography;

export default function CompanyPage() {
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem('hr_company');
    if (saved) {
      try {
        form.setFieldsValue(JSON.parse(saved));
      } catch { /* ignore */ }
    }
  }, [form]);

  const handleSave = async (values: any) => {
    setSaving(true);
    try {
      localStorage.setItem('hr_company', JSON.stringify(values));
      message.success('企业信息已保存');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <Title level={4}><BankOutlined /> 企业信息</Title>
      <Card>
        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item name="company_name" label="企业名称" rules={[{ required: true, message: '请输入企业名称' }]}>
            <Input placeholder="如：放心科技有限公司" />
          </Form.Item>
          <Form.Item label="企业地址（发布职位时自动复用）" required>
            <Input.Group compact>
              <Form.Item name="province" noStyle rules={[{ required: true, message: '请输入省份' }]}>
                <Input placeholder="省" style={{ width: '25%' }} />
              </Form.Item>
              <Form.Item name="city" noStyle rules={[{ required: true, message: '请输入城市' }]}>
                <Input placeholder="市" style={{ width: '25%' }} />
              </Form.Item>
              <Form.Item name="district" noStyle>
                <Input placeholder="区" style={{ width: '25%' }} />
              </Form.Item>
              <Form.Item name="detail" noStyle>
                <Input placeholder="详细地址" style={{ width: '25%' }} />
              </Form.Item>
            </Input.Group>
          </Form.Item>
          <Form.Item name="description" label="企业简介">
            <Input.TextArea rows={3} placeholder="介绍一下公司..." />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>保存</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
