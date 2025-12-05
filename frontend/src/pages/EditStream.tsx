import { useEffect, useState } from 'react';

import { useNavigate, useParams } from 'react-router-dom';
import client from '@/api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

export default function EditStream() {
  const navigate = useNavigate();
  const { uuid } = useParams<{ uuid: string }>();
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [onDemand, setOnDemand] = useState(true);
  const [audio, setAudio] = useState(false);
  const [debug, setDebug] = useState(false);
  const [loading, setLoading] = useState(true);
  const [defaultProtocol, setDefaultProtocol] = useState<string>('auto');

  useEffect(() => {
    if (uuid) {
        const savedProtocol = localStorage.getItem(`stream_protocol_${uuid}`);
        if (savedProtocol) setDefaultProtocol(savedProtocol);
    }
  }, [uuid]);

  useEffect(() => {
    const fetchStream = async () => {
      try {
        const response = await client.get(`/stream/${uuid}/info`);
        if (response.data.status === 1) {
          const stream = response.data.payload;
          setName(stream.name);
          // Assuming single channel "0" for now as per AddStream simplification
          const channel = stream.channels?.["0"];
          if (channel) {
            setUrl(channel.url);
            setOnDemand(channel.on_demand);
            setAudio(channel.audio);
            setDebug(channel.debug);
          }
        }
      } catch (error) {
        console.error('Failed to fetch stream info', error);
      } finally {
        setLoading(false);
      }
    };

    if (uuid) {
      fetchStream();
    }
  }, [uuid]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const payload = {
      name,
      channels: {
        "0": {
          name: "ch1",
          url,
          on_demand: onDemand,
          debug,
          audio,
        }
      }
    };

    try {
      await client.post(`/stream/${uuid}/edit`, payload);

      if (defaultProtocol) {
          localStorage.setItem(`stream_protocol_${uuid}`, defaultProtocol);
      } else {
          localStorage.removeItem(`stream_protocol_${uuid}`);
      }

      navigate('/');
    } catch (error) {
      console.error('Failed to edit stream', error);
      alert('Failed to edit stream');
    }
  };

  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="p-8 flex justify-center">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Edit Stream</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Stream Name</Label>
              <Input id="name" value={name} onChange={e => setName(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="url">RTSP URL</Label>
              <Input id="url" value={url} onChange={e => setUrl(e.target.value)} required />
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="onDemand" checked={onDemand} onCheckedChange={(c) => setOnDemand(!!c)} />
              <Label htmlFor="onDemand">On Demand</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="audio" checked={audio} onCheckedChange={(c) => setAudio(!!c)} />
              <Label htmlFor="audio">Audio</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="debug" checked={debug} onCheckedChange={(c) => setDebug(!!c)} />
              <Label htmlFor="debug">Debug</Label>
            </div>

            <div className="space-y-2">
                <Label htmlFor="protocol">Default Protocol</Label>
                <Select value={defaultProtocol} onValueChange={setDefaultProtocol}>
                    <SelectTrigger>
                        <SelectValue placeholder="Select Protocol" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="auto">Auto (WebRTC to MSE to HLS)</SelectItem>
                        <SelectItem value="webrtc">WebRTC</SelectItem>
                        <SelectItem value="mse">MSE</SelectItem>
                        <SelectItem value="hls">HLS</SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <Button type="submit" className="w-full">Save Changes</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
