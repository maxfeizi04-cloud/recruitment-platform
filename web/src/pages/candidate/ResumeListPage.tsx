import React, { useEffect, useState } from 'react';
import { Card, Typography, Button, Popconfirm, Modal, Input, message, Tag } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, StarOutlined, StarFilled } from '@ant-design/icons';
import { listResumes, createResume, updateResume, deleteResume, setDefaultResume, type Resume } from '../../api/resume';

const { Title, Text } = Typography;

export default function ResumeListPage() {
  const [resumes, setResumes] = useState<Resume[]>([]);
  const [loading, setLoading] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [editing, setEditing] = useState<Resume | null>(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [saving, setSaving] = useState(false);

  const fetch = async () => { setLoading(true); try { setResumes(await listResumes()); } finally { setLoading(false); } };
  useEffect(() => { fetch(); }, []);

  const handleSave = async () => {
    if (!title.trim()) { message.error('请输入简历标题'); return; }
    setSaving(true);
    try {
      if (editing) {
        await updateResume(editing.id, title, content);
      } else {
        await createResume(title, content || '{}');
      }
      message.success(editing ? '已更新' : '已创建');
      setEditOpen(false);
      fetch();
    } finally { setSaving(false); }
  };

  const openEdit = (r?: Resume) => {
    setEditing(r || null);
    setTitle(r?.title || '');
    setContent(r?.content || '');
    setEditOpen(true);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>我的简历</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => openEdit()}>新建简历</Button>
      </div>
      {resumes.map(r => (
        <Card key={r.id} style={{ marginBottom: 12 }} actions={[
          <Button type="link" icon={<EditOutlined />} onClick={() => openEdit(r)}>编辑</Button>,
          <Button type="link" icon={r.is_default ? <StarFilled /> : <StarOutlined />} onClick={() => setDefaultResume(r.id).then(fetch)}>{r.is_default ? '默认' : '设为默认'}</Button>,
          <Popconfirm title="确定删除？" onConfirm={() => deleteResume(r.id).then(fetch)}><Button type="link" danger icon={<DeleteOutlined />}>删除</Button></Popconfirm>,
        ]}>
          <Text strong>{r.title}</Text>
          {r.is_default && <Tag color="gold" style={{ marginLeft: 8 }}>默认</Tag>}
          <div style={{ marginTop: 8 }}><Text type="secondary">更新于 {r.updated_at?.slice(0, 10)}</Text></div>
        </Card>
      ))}
      {resumes.length === 0 && !loading && <Text type="secondary">暂无简历，点击"新建简历"创建</Text>}

      <Modal title={editing ? '编辑简历' : '新建简历'} open={editOpen} onOk={handleSave} onCancel={() => setEditOpen(false)} confirmLoading={saving}>
        <Input placeholder="简历标题（如：3年前端工程师）" value={title} onChange={e => setTitle(e.target.value)} style={{ marginBottom: 12 }} />
        <Input.TextArea rows={6} placeholder='简历内容（JSON格式，如：{"skills":["React"]}' value={content} onChange={e => setContent(e.target.value)} />
      </Modal>
    </div>
  );
}
