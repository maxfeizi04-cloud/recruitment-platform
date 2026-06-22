import React, { useEffect, useState } from 'react';
import { Card, Typography, Button, Popconfirm, Modal, Input, message, Tag, Upload, Form, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, StarOutlined, StarFilled, UploadOutlined, EnvironmentOutlined, FileOutlined } from '@ant-design/icons';
import { listResumes, createResume, updateResume, deleteResume, setDefaultResume, uploadAttachment, type Resume } from '../../api/resume';

const { Title, Text } = Typography;

function parseContent(content: string) {
  try { return JSON.parse(content); } catch { return {}; }
}

export default function ResumeListPage() {
  const [resumes, setResumes] = useState<Resume[]>([]);
  const [loading, setLoading] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [editing, setEditing] = useState<Resume | null>(null);
  const [saving, setSaving] = useState(false);
  const [form] = Form.useForm();

  const fetch = async () => { setLoading(true); try { setResumes(await listResumes()); } finally { setLoading(false); } };
  useEffect(() => { fetch(); }, []);

  const openEdit = (r?: Resume) => {
    setEditing(r || null);
    if (r) {
      const c = parseContent(r.content);
      form.setFieldsValue({
        title: r.title,
        name: c.name || '',
        phone: c.phone || '',
        email: c.email || '',
        city: c.city || '',
        skills: (c.skills || []).join(', '),
        experience: c.experience || '',
        education: c.education || '',
        summary: c.summary || '',
      });
    } else {
      form.resetFields();
    }
    setEditOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields().catch(() => null);
    if (!values) return;
    setSaving(true);
    const content = JSON.stringify({
      name: values.name || '',
      phone: values.phone || '',
      email: values.email || '',
      city: values.city || '',
      skills: (values.skills || '').split(',').map((s: string) => s.trim()).filter(Boolean),
      experience: values.experience || '',
      education: values.education || '',
      summary: values.summary || '',
    });
    try {
      if (editing) {
        await updateResume(editing.id, values.title, content);
        message.success('简历已更新');
      } else {
        await createResume(values.title, content);
        message.success('简历已创建');
      }
      setEditOpen(false);
      fetch();
    } catch { /* handled by interceptor */ }
    finally { setSaving(false); }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <Title level={4} className="!mb-0">我的简历</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => openEdit()}>新建简历</Button>
      </div>

      {loading && <Text type="secondary">加载中...</Text>}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {resumes.map(r => {
          const c = parseContent(r.content);
          return (
            <Card key={r.id} hoverable className="rounded-xl shadow-sm border-slate-100"
              actions={[
                <Button type="link" icon={<EditOutlined />} onClick={() => openEdit(r)}>编辑</Button>,
                <Button type="link" icon={r.is_default ? <StarFilled className="text-amber-400" /> : <StarOutlined />}
                  onClick={() => setDefaultResume(r.id).then(fetch)}>
                  {r.is_default ? '默认简历' : '设为默认'}
                </Button>,
                <Popconfirm title="确定删除此简历？" onConfirm={() => deleteResume(r.id).then(fetch)}>
                  <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
                </Popconfirm>,
              ]}
            >
              <div className="flex items-start justify-between">
                <div>
                  <Text strong className="text-base">{r.title}</Text>
                  {r.is_default && <Tag color="gold" className="ml-2">默认</Tag>}
                </div>
              </div>

              {/* Preview info */}
              <div className="mt-3 space-y-1 text-sm text-slate-500">
                {c.name && <div><span className="text-slate-400">姓名：</span>{c.name}</div>}
                {c.city && <div><EnvironmentOutlined className="mr-1 text-xs" />{c.city}</div>}
                {c.skills?.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-2">
                    {c.skills.map((s: string) => <Tag key={s} className="text-xs bg-blue-50 text-blue-600 border-0">{s}</Tag>)}
                  </div>
                )}
              </div>

              {/* Attachments */}
              <div className="mt-3 pt-3 border-t border-slate-50">
                {r.attachment_urls?.length > 0 ? (
                  r.attachment_urls.map((url: string, idx: number) => (
                    <Tag key={idx} className="mb-1"><a href={url} target="_blank" rel="noopener noreferrer">附件{idx + 1}</a></Tag>
                  ))
                ) : null}
                <Upload accept=".pdf,.doc,.docx,.jpg,.jpeg,.png" showUploadList={false}
                  beforeUpload={async (file) => {
                    if (file.size > 20 * 1024 * 1024) { message.error('文件不能超过20MB'); return false; }
                    try { await uploadAttachment(r.id, file); message.success('上传成功'); fetch(); } catch {}
                    return false;
                  }}
                >
                  <Button icon={<UploadOutlined />} size="small" type="link">上传附件</Button>
                </Upload>
              </div>

              <div className="mt-2 text-xs text-slate-400">更新于 {r.updated_at?.slice(0, 10)}</div>
            </Card>
          );
        })}
      </div>

      {resumes.length === 0 && !loading && (
        <div className="text-center py-16 text-slate-400">
          <FileOutlined style={{ fontSize: 48 }} />
          <p className="mt-3">暂无简历，点击"新建简历"开始</p>
        </div>
      )}

      {/* Edit Modal */}
      <Modal title={editing ? '编辑简历' : '新建简历'} open={editOpen} onOk={handleSave}
        onCancel={() => setEditOpen(false)} confirmLoading={saving} width={640} okText="保存" cancelText="取消"
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item name="title" label="简历标题" rules={[{ required: true, message: '请输入简历标题' }]}>
            <Input placeholder="如：3年前端工程师简历" />
          </Form.Item>
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="name" label="姓名"><Input placeholder="你的姓名" /></Form.Item>
            <Form.Item name="phone" label="手机号"><Input placeholder="手机号" /></Form.Item>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="email" label="邮箱"><Input placeholder="邮箱地址" /></Form.Item>
            <Form.Item name="city" label="期望城市"><Input placeholder="如：上海" /></Form.Item>
          </div>
          <Form.Item name="skills" label="技能标签（逗号分隔）">
            <Input placeholder="React, TypeScript, Go, Docker" />
          </Form.Item>
          <Form.Item name="experience" label="工作经历">
            <Input.TextArea rows={3} placeholder="简述你的工作经历..." />
          </Form.Item>
          <Form.Item name="education" label="教育背景">
            <Input placeholder="如：浙江大学 · 计算机科学 · 本科 · 2019" />
          </Form.Item>
          <Form.Item name="summary" label="个人简介">
            <Input.TextArea rows={2} placeholder="简单介绍一下自己..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
