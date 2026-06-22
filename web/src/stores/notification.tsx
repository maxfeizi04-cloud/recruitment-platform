import { createContext, useContext, useEffect, useRef, useState, useCallback } from 'react';
import { message } from 'antd';
import { useAuth } from './auth';

interface WSEvent {
  type: string;
  payload: any;
  time: string;
}

interface NotificationCtx {
  connected: boolean;
}

const Ctx = createContext<NotificationCtx>({ connected: false });
export const useNotification = () => useContext(Ctx);

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const auth = useAuth();
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!auth.token) return;

    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${proto}//${location.host}/api/ws?token=${auth.token}`);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);

    ws.onmessage = (e) => {
      try {
        const event: WSEvent = JSON.parse(e.data);
        switch (event.type) {
          case 'application.updated':
            message.info(`投递状态已更新为：${event.payload?.status || '已更新'}`);
            break;
          case 'interview.created':
            message.success(event.payload?.message || '收到新的面试邀约');
            break;
          case 'interview.updated':
            message.info(`面试状态已更新：${event.payload?.status || '已更新'}`);
            break;
        }
      } catch { /* ignore */ }
    };

    return () => { ws.close(); };
  }, [auth.token]);

  return <Ctx.Provider value={{ connected }}>{children}</Ctx.Provider>;
}
