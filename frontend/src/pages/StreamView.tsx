import { useParams } from 'react-router-dom';
import { Player } from '@/components/Player';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useState } from 'react';

export default function StreamView() {
  const { uuid } = useParams<{ uuid: string }>();
  const [type, setType] = useState<'webrtc' | 'mse' | 'hls' | 'auto'>('auto');

  if (!uuid) return <div>ID de transmisión inválido</div>;

  return (
    <div className="p-8 h-full flex flex-col min-h-[calc(100vh-4rem)]">
      <h1 className="text-2xl font-bold mb-4">Transmisión: {uuid}</h1>
      <div className="flex-1 min-h-0">
        <Tabs value={type} onValueChange={(v) => setType(v as any)} className="h-full flex flex-col">
          <TabsList>
            <TabsTrigger value="auto">Auto</TabsTrigger>
            <TabsTrigger value="webrtc">WebRTC</TabsTrigger>
            <TabsTrigger value="mse">MSE</TabsTrigger>
            <TabsTrigger value="hls">HLS</TabsTrigger>
          </TabsList>
          <div className="flex-1 mt-4 bg-black rounded-lg overflow-hidden min-h-[500px]">
             <Player uuid={uuid} type={type} />
          </div>
        </Tabs>
      </div>
    </div>
  );
}
