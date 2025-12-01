import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '@/api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function AddStream() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [onDemand, setOnDemand] = useState(true);
  const [audio, setAudio] = useState(false);
  const [debug, setDebug] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const uuid = crypto.randomUUID();
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
      await client.post(`/stream/${uuid}/add`, payload);
      navigate('/');
    } catch (error) {
      console.error('Failed to add stream', error);
      alert('Failed to add stream');
    }
  };

  return (
    <div className="p-8 flex justify-center">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Add New Stream</CardTitle>
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
            <Button type="submit" className="w-full">Add Stream</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
