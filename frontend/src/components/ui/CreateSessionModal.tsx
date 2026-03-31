"use client";

import { useState } from 'react';

export function CreateSessionModal({
 isOpen,
 onClose,
 onCreate
}: {
 isOpen: boolean;
 onClose: () => void;
 onCreate: (config: { resolution: string, fps: number, bitrate: number }) => Promise<void>;
}) {
 const [resolution, setResolution] = useState("1920x1080");
 const [fps, setFps] = useState(30);
 const [bitrate, setBitrate] = useState(8000);
 const [loading, setLoading] = useState(false);

 if (!isOpen) return null;

 const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  setLoading(true);
  try {
   await onCreate({ resolution, fps, bitrate });
   onClose();
  } catch (err) {
   alert("Failed to create session. " + err);
  } finally {
   setLoading(false);
  }
 };

 return (
  <div className="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50 backdrop-blur-sm">
   <div className="bg-gray-800 rounded-xl p-6 w-full max-w-md border border-gray-700 shadow-2xl">
    <h2 className="text-2xl font-bold text-white mb-6">Create Virtual Display</h2>

    <form onSubmit={handleSubmit} className="space-y-4">
     <div>
      <label className="block text-sm font-medium text-gray-300 mb-1">Resolution</label>
      <select
       value={resolution}
       onChange={e => setResolution(e.target.value)}
       className="w-full bg-gray-900 border border-gray-600 rounded p-2 text-white focus:outline-none focus:border-blue-500 transition-colors"
      >
       <option value="1920x1080">1920x1080 (1080p)</option>
       <option value="2560x1440">2560x1440 (1440p)</option>
       <option value="3840x2160">3840x2160 (4K)</option>
       <option value="1280x720">1280x720 (720p)</option>
      </select>
     </div>

     <div>
      <label className="block text-sm font-medium text-gray-300 mb-1">FPS (Frames Per Second)</label>
      <input
       type="number"
       value={fps}
       onChange={e => setFps(Number(e.target.value))}
       min={15} max={144}
       className="w-full bg-gray-900 border border-gray-600 rounded p-2 text-white focus:outline-none focus:border-blue-500 transition-colors"
      />
     </div>

     <div>
      <label className="block text-sm font-medium text-gray-300 mb-1">Bitrate (kbps)</label>
      <input
       type="number"
       value={bitrate}
       onChange={e => setBitrate(Number(e.target.value))}
       min={1000} max={50000} step={1000}
       className="w-full bg-gray-900 border border-gray-600 rounded p-2 text-white focus:outline-none focus:border-blue-500 transition-colors"
      />
     </div>

     <div className="pt-4 flex justify-end gap-3">
      <button
       type="button"
       onClick={onClose}
       className="px-4 py-2 text-gray-300 hover:text-white transition-colors"
      >
       Cancel
      </button>
      <button
       type="submit"
       disabled={loading}
       className="px-6 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded font-medium disabled:opacity-50 transition-colors"
      >
       {loading ? 'Creating...' : 'Create Display'}
      </button>
     </div>
    </form>
   </div>
  </div>
 );
}
