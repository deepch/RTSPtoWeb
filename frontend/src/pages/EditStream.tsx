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
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
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
          const channel = stream.channels?.["0"];
          if (channel) {
            // Parse existing URL to extract credentials
            try {
                // If the URL is just an IP or invalid, this might fail or produce empty parts
                let fullUrl = channel.url;
                if (!fullUrl.includes('://')) fullUrl = 'rtsp://' + fullUrl;

                const u = new URL(fullUrl);

                if (u.username) setUsername(decodeURIComponent(u.username));
                if (u.password) setPassword(decodeURIComponent(u.password));

                // Set URL field to the "clean" URL without credentials
                // Reconstruct: protocol + host + path...
                // Using u.host includes port.
                u.username = '';
                u.password = '';
                setUrl(u.toString());
            } catch (e) {
                // If parsing fails, just set the raw URL
                setUrl(channel.url);
            }

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

    let finalUrl = url;
    try {
        if (username || password) {
            let parseUrl = url;
            if (!parseUrl.includes('://')) {
                parseUrl = 'rtsp://' + parseUrl;
            }
            const u = new URL(parseUrl);

            const encodedUser = encodeURIComponent(username);
            const encodedPass = encodeURIComponent(password);
            const auth = `${encodedUser}:${encodedPass}`;

            finalUrl = `${u.protocol}//${auth}@${u.host}${u.pathname}${u.search}${u.hash}`;
        }
    } catch (err) {
        console.error("Invalid URL format", err);
        alert("URL inválida");
        return;
    }

    const payload = {
      name,
      channels: {
        "0": {
          name: "ch1",
          url: finalUrl,
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
      alert('Error al editar la transmisión');
    }
  };

  if (loading) return <div className="p-8">Cargando...</div>;

  return (
    <div className="p-8 flex justify-center">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Editar Transmisión</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nombre de la Transmisión</Label>
              <Input id="name" value={name} onChange={e => setName(e.target.value)} required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="url">URL RTSP (Sin credenciales)</Label>
              <Input id="url" value={url} onChange={e => setUrl(e.target.value)} required />
              <p className="text-xs text-muted-foreground">La URL se limpia automáticamente al cargar si tiene credenciales.</p>
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="username">Usuario (Opcional)</Label>
                  <Input id="username" value={username} onChange={e => setUsername(e.target.value)} placeholder="admin" />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Contraseña (Opcional)</Label>
                  <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="••••••" />
                </div>
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox id="onDemand" checked={onDemand} onCheckedChange={(c) => setOnDemand(!!c)} />
              <Label htmlFor="onDemand">Bajo Demanda</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="audio" checked={audio} onCheckedChange={(c) => setAudio(!!c)} />
              <Label htmlFor="audio">Audio</Label>
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox id="debug" checked={debug} onCheckedChange={(c) => setDebug(!!c)} />
              <Label htmlFor="debug">Depuración</Label>
            </div>

            <div className="space-y-2">
                <Label htmlFor="protocol">Protocolo Predeterminado</Label>
                <Select value={defaultProtocol} onValueChange={setDefaultProtocol}>
                    <SelectTrigger>
                        <SelectValue placeholder="Seleccionar Protocolo" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="auto">Auto (WebRTC a MSE a HLS)</SelectItem>
                        <SelectItem value="webrtc">WebRTC</SelectItem>
                        <SelectItem value="mse">MSE</SelectItem>
                        <SelectItem value="hls">HLS</SelectItem>
                    </SelectContent>
                </Select>
            </div>
            <Button type="submit" className="w-full">Guardar Cambios</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
