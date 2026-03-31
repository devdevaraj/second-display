import { api } from './api';

export class WebRTCClient {
 private pc: RTCPeerConnection | null = null;
 private videoElement: HTMLVideoElement;
 private sessionId: string;

 constructor(sessionId: string, videoElement: HTMLVideoElement) {
  this.sessionId = sessionId;
  this.videoElement = videoElement;
 }

 async connect() {
  this.pc = new RTCPeerConnection({
   iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
  });

  // When the backend sends a track
  this.pc.ontrack = (event) => {
   console.log('Received track:', event.track.kind);
   if (this.videoElement.srcObject !== event.streams[0]) {
    this.videoElement.srcObject = event.streams[0];
    console.log('Attached stream to video element');
   }
  };

  // Pion backend requires transceiver for recvonly
  this.pc.addTransceiver('video', { direction: 'recvonly' });

  const offer = await this.pc.createOffer();
  await this.pc.setLocalDescription(offer);

  // Wait for ICE candidates to be gathered before sending the offer
  // This allows Pion to receive all candidates in the initial SDP
  await new Promise<void>((resolve) => {
   if (this.pc?.iceGatheringState === 'complete') {
    resolve();
   } else {
    const checkState = () => {
     if (this.pc?.iceGatheringState === 'complete') {
      this.pc.removeEventListener('icegatheringstatechange', checkState);
      resolve();
     }
    };
    this.pc?.addEventListener('icegatheringstatechange', checkState);
    // Timeout after 2 seconds to not block forever
    setTimeout(resolve, 2000);
   }
  });

  // Send complete offer to backend
  const answer = await api.sendWebRTCOffer(this.sessionId, this.pc.localDescription!);
  await this.pc.setRemoteDescription(new RTCSessionDescription(answer));
 }

 disconnect() {
  this.pc?.close();
  this.pc = null;
  if (this.videoElement) {
   this.videoElement.srcObject = null;
  }
 }
}
