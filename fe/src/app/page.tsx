"use client";

import { useEffect, useState } from "react";
import { api, Session } from "@/lib/api";
import { SessionCard } from "@/components/ui/SessionCard";
import { CreateSessionModal } from "@/components/ui/CreateSessionModal";

export default function Dashboard() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [isModalOpen, setModalOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchSessions = async () => {
    try {
      const data = await api.getSessions();
      setSessions(data || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || "Failed to fetch sessions. Ensure Go backend is running.");
    }
  };

  useEffect(() => {
    fetchSessions();
    const interval = setInterval(fetchSessions, 3000);
    return () => clearInterval(interval);
  }, []);

  const handleCreate = async (config: { resolution: string, fps: number, bitrate: number }) => {
    await api.createSession(config);
    await fetchSessions();
  };

  const handleStart = async (id: string) => {
    try {
      await api.startSession(id);
      fetchSessions();
    } catch (err: any) { alert(err); }
  };

  const handleStop = async (id: string) => {
    try {
      await api.stopSession(id);
      fetchSessions();
    } catch (err: any) { alert(err); }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this session forever?")) return;
    try {
      await api.deleteSession(id);
      fetchSessions();
    } catch (err: any) { alert(err); }
  };

  return (
    <main className="min-h-screen bg-gray-900 text-white p-8">
      <div className="max-w-6xl mx-auto">
        <header className="flex justify-between items-center mb-12 border-b border-gray-800 pb-6">
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent">
              Virtual Display Manager
            </h1>
            <p className="text-gray-400 mt-2">Manage headless Wayland outputs remotely</p>
          </div>
          <button
            onClick={() => setModalOpen(true)}
            className="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-lg font-semibold shadow-lg shadow-blue-500/30 transition-all flex items-center gap-2"
          >
            <span>+</span> New Display
          </button>
        </header>

        {error && (
          <div className="bg-red-500/10 border border-red-500 text-red-400 p-4 rounded-lg mb-8">
            {error}
          </div>
        )}

        {sessions.length === 0 && !error ? (
          <div className="text-center py-20 border-2 border-dashed border-gray-800 rounded-xl">
            <div className="text-gray-500 mb-4 text-6xl">🖥️</div>
            <h3 className="text-xl font-medium text-gray-300">No active displays</h3>
            <p className="text-gray-500 mt-2">Create your first virtual display to get started.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {sessions.map(session => (
              <SessionCard
                key={session.id}
                session={session}
                onStart={handleStart}
                onStop={handleStop}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        <CreateSessionModal
          isOpen={isModalOpen}
          onClose={() => setModalOpen(false)}
          onCreate={handleCreate}
        />
      </div>
    </main>
  );
}
