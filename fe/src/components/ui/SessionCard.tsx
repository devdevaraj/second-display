import { Session } from '@/lib/api';
import Link from 'next/link';

export function SessionCard({ session, onStart, onStop, onDelete }: {
 session: Session,
 onStart: (id: string) => void,
 onStop: (id: string) => void,
 onDelete: (id: string) => void
}) {
 return (
  <div className="bg-gray-800 rounded-lg p-6 shadow-lg border border-gray-700 hover:border-blue-500 transition-colors">
   <div className="flex justify-between items-start mb-4">
    <div>
     <h3 className="text-xl font-bold text-white mb-1">{session.outputName}</h3>
     <p className="text-sm text-gray-400">ID: {session.id.substring(0, 8)}</p>
    </div>
    <span className={`px-3 py-1 rounded-full text-xs font-semibold ${session.status === 'running' ? 'bg-green-500/20 text-green-400' :
      session.status === 'error' ? 'bg-red-500/20 text-red-400' :
       'bg-gray-600/50 text-gray-300'
     }`}>
     {session.status.toUpperCase()}
    </span>
   </div>

   <div className="space-y-2 mb-6 text-sm text-gray-300">
    <div className="flex justify-between">
     <span>Resolution:</span>
     <span className="font-mono bg-gray-900 px-2 py-0.5 rounded">{session.resolution}</span>
    </div>
    <div className="flex justify-between">
     <span>FPS / Bitrate:</span>
     <span className="font-mono bg-gray-900 px-2 py-0.5 rounded">{session.fps} / {session.bitrate}k</span>
    </div>
   </div>

   <div className="flex gap-2">
    {session.status !== 'running' ? (
     <button
      onClick={() => onStart(session.id)}
      className="flex-1 bg-blue-600 hover:bg-blue-500 text-white py-2 rounded transition-colors"
     >
      Start
     </button>
    ) : (
     <button
      onClick={() => onStop(session.id)}
      className="flex-1 bg-yellow-600 hover:bg-yellow-500 text-white py-2 rounded transition-colors"
     >
      Stop
     </button>
    )}

    {session.status === 'running' && (
     <Link href={`/session/${session.id}`} className="flex-1 bg-green-600 hover:bg-green-500 text-white py-2 rounded transition-colors text-center inline-block">
      View
     </Link>
    )}

    <button
     onClick={() => onDelete(session.id)}
     className="px-4 bg-gray-700 hover:bg-red-600 text-white rounded transition-colors"
     title="Delete Session"
    >
     🗑️
    </button>
   </div>
  </div>
 );
}
