import { useEffect, useRef, useState } from 'react';
import Hls from 'hls.js';
import client from '@/api/client';

interface PlayerProps {
  uuid: string;
  channel?: string;
  type?: 'webrtc' | 'mse' | 'hls';
}

export function Player({ uuid, channel = "0", type = 'webrtc' }: PlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [error, setError] = useState<string | null>(null);
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const hlsRef = useRef<Hls | null>(null);

  useEffect(() => {
    return () => {
      cleanup();
    };
  }, [uuid, channel, type]);

  useEffect(() => {
    cleanup();
    setError(null);

    if (type === 'webrtc') {
      startWebRTC();
    } else if (type === 'mse') {
      startMSE();
    } else if (type === 'hls') {
      startHLS();
    }
  }, [uuid, channel, type]);

  const cleanup = () => {
    if (pcRef.current) {
      pcRef.current.close();
      pcRef.current = null;
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    if (hlsRef.current) {
      hlsRef.current.destroy();
      hlsRef.current = null;
    }
    if (videoRef.current) {
      videoRef.current.src = '';
      videoRef.current.srcObject = null;
    }
  };

  const startWebRTC = async () => {
    try {
      const pc = new RTCPeerConnection({
        iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
      });
      pcRef.current = pc;

      pc.ontrack = (event) => {
        if (videoRef.current) {
          videoRef.current.srcObject = event.streams[0];
        }
      };

      pc.addTransceiver('video', { direction: 'recvonly' });

      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);

      const payload = btoa(pc.localDescription?.sdp || '');
      const formData = new URLSearchParams();
      formData.append('data', payload);

      const response = await client.post(`/stream/${uuid}/channel/${channel}/webrtc`, formData, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });

      const answer = atob(response.data);
      await pc.setRemoteDescription(new RTCSessionDescription({ type: 'answer', sdp: answer }));

    } catch (err: any) {
      console.error('WebRTC Error:', err);
      setError('WebRTC connection failed: ' + err.message);
    }
  };

  const startMSE = () => {
    const video = videoRef.current;
    if (!video) return;

    const mse = new MediaSource();
    video.src = URL.createObjectURL(mse);

    mse.addEventListener('sourceopen', () => {
      const ws = new WebSocket(`ws://${window.location.hostname}:8083/stream/${uuid}/channel/${channel}/mse`);
      wsRef.current = ws;
      ws.binaryType = 'arraybuffer';

      let sourceBuffer: SourceBuffer | null = null;
      const queue: ArrayBuffer[] = [];
      let isAppending = false;

      ws.onopen = () => {
        console.log('MSE WS Connected');
      };

      ws.onmessage = (event) => {
        const data = new Uint8Array(event.data);
        if (data[0] === 9) { // Mime type packet
             const mimeCodec = new TextDecoder("utf-8").decode(data.slice(1));
             if (!sourceBuffer && mse.readyState === 'open') {
                 try {
                    sourceBuffer = mse.addSourceBuffer(`video/mp4; codecs="${mimeCodec}"`);
                    sourceBuffer.mode = 'segments';
                    sourceBuffer.addEventListener('updateend', () => {
                        isAppending = false;
                        processQueue();
                    });
                 } catch (e) {
                     console.error('MSE addSourceBuffer error', e);
                 }
             }
        } else {
            if (sourceBuffer) {
                queue.push(event.data);
                processQueue();
            }
        }
      };

      const processQueue = () => {
          if (sourceBuffer && !isAppending && queue.length > 0) {
              isAppending = true;
              try {
                  sourceBuffer.appendBuffer(queue.shift()!);
              } catch (e) {
                  console.error('MSE appendBuffer error', e);
                  isAppending = false;
              }
          }
      };
    });
  };

  const startHLS = () => {
    const video = videoRef.current;
    if (!video) return;

    const src = `http://${window.location.hostname}:8083/stream/${uuid}/channel/${channel}/hls/live/index.m3u8`;

    if (Hls.isSupported()) {
      const hls = new Hls();
      hlsRef.current = hls;
      hls.loadSource(src);
      hls.attachMedia(video);
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = src;
    }
  };

  return (
    <div className="relative w-full h-full bg-black rounded-lg overflow-hidden">
      {error && (
        <div className="absolute inset-0 flex items-center justify-center text-red-500 bg-black/80 z-10">
          {error}
        </div>
      )}
      <video
        ref={videoRef}
        autoPlay
        muted
        playsInline
        controls
        className="w-full h-full object-contain"
      />
    </div>
  );
}
