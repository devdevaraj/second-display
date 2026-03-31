"use client";

import { useEffect, useRef, useState } from "react";
import { usePathname, useParams } from "next/navigation";
import { api, Session } from "@/lib/api";
import { WebRTCClient } from "@/lib/webrtc";
import Link from "next/link";

export default function SessionStream() {
 const params = useParams();
 const id = params.id as string;
 const videoRef = useRef<HTMLVideoElement>(null);

 const [session, setSession] = useState<Session | null>(null);
 const [isConnected, setIsConnected] = useState(false);
 const [error, setError] = useState<string | null>(null);

 useEffect(() => {
  // Fetch initial details
  api.getSessions().then(sessions => {
   const sess = sessions.find(s => s.id === id);
   if (sess) setSession(sess);
   else setError("Session not found");
  }).catch(err => setError(err.message));
 }, [id]);

 useEffect(() => {
  if (!videoRef.current || !id) return;

  let client: WebRTCClient | null = null;

  const initStream = async () => {
   try {
    client = new WebRTCClient(id, videoRef.current!);
    await client.connect();
    setIsConnected(true);
   } catch (err: any) {
    setError("Failed to connect WebRTC: " + err.message);
   }
  };

  initStream();

  return () => {
   if (client) {
    client.disconnect();
   }
  };
 }, [id]);

 return (
  <div className="min-h-screen bg-black flex flex-col">
   <header className="bg-gray-900 border-b border-gray-800 p-4 flex justify-between items-center text-white z-10 w-full shadow-lg h-16 shrink-0">
    <div className="flex items-center gap-4">
     <Link href="/" className="text-gray-400 hover:text-white transition-colors">
      ← Back
     </Link>
     <div className="h-6 w-px bg-gray-700"></div>
     <h1 className="font-semibold">{session ? session.outputName : "Loading..."}</h1>
     <span className={`px-2 py-0.5 rounded text-xs font-mono ml-2 ${isConnected ? 'bg-green-500/20 text-green-400' : 'bg-yellow-500/20 text-yellow-500'}`}>
      {isConnected ? 'LIVE' : 'CONNECTING...'}
     </span>
    </div>

    {session && (
     <div className="text-sm font-mono text-gray-400">
      {session.resolution} @ {session.fps}FPS
     </div>
    )}
   </header>

   <main className="flex-1 relative w-full h-full flex items-center justify-center bg-zinc-950 overflow-hidden">
    {error && (
     <div className="absolute inset-x-8 top-8 z-20 bg-red-500/10 border border-red-500 text-red-400 p-4 rounded-lg backdrop-blur-sm">
      {error}
     </div>
    )}

    {/* Fill available space while keeping aspect ratio */}
    <div className="absolute inset-0 p-4 md:p-8 flex items-center justify-center pointer-events-none">
     <video
      ref={videoRef}
      autoPlay
      playsInline
      muted // Necessary for autoplay in many browsers
      className="max-w-full max-h-full aspect-video object-contain shadow-2xl shadow-black rounded-lg bg-black pointer-events-auto ring-1 ring-gray-800"
     />
    </div>
   </main>
  </div>
 );
}
