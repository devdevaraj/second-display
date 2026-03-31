export interface Session {
  id: string;
  outputName: string;
  resolution: string;
  fps: number;
  bitrate: number;
  status: 'created' | 'running' | 'stopped' | 'error';
  streamUrl?: string;
}

const API_BASE = 'http://localhost:8080/api';

export const api = {
  getSessions: async (): Promise<Session[]> => {
    const res = await fetch(`${API_BASE}/sessions`, { cache: 'no-store' });
    if (!res.ok) throw new Error('Failed to fetch sessions');
    return res.json();
  },

  createSession: async (options: { resolution: string; fps: number; bitrate: number }): Promise<Session> => {
    const res = await fetch(`${API_BASE}/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(options),
    });
    if (!res.ok) throw new Error('Failed to create session');
    return res.json();
  },

  startSession: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/sessions/${id}/start`, { method: 'POST' });
    if (!res.ok) throw new Error('Failed to start session');
  },

  stopSession: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/sessions/${id}/stop`, { method: 'POST' });
    if (!res.ok) throw new Error('Failed to stop session');
  },

  deleteSession: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/sessions/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error('Failed to delete session');
  },
  
  sendWebRTCOffer: async (id: string, offer: RTCSessionDescriptionInit): Promise<RTCSessionDescriptionInit> => {
    const res = await fetch(`${API_BASE}/sessions/${id}/webrtc`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(offer),
    });
    if (!res.ok) throw new Error('Failed to negotiate WebRTC');
    return res.json();
  }
};
