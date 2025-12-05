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
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [onDemand, setOnDemand] = useState(true);
  const [audio, setAudio] = useState(false);
  const [debug, setDebug] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const uuid = crypto.randomUUID();

    let finalUrl = url;
    try {
        if (username || password) {
            // Flexible parsing: if user entered just "192.168.1.1", add rtsp://
            let parseUrl = url;
            if (!parseUrl.includes('://')) {
                parseUrl = 'rtsp://' + parseUrl;
            }
            const u = new URL(parseUrl);

            // Encode credentials to handle special characters like $$ -> %24%24
            if (username) u.username = encodeURIComponent(username);
            if (password) u.password = encodeURIComponent(password);

            // Note: new URL() might double encode if we pass encoded chars to properties?
            // Actually, u.username = 'val' encodes dangerous chars, but not all.
            // But if we explicitly used encodeURIComponent above, we might get double encoding if we set it to .username directly?
            // Let's verify: u.username = '%24' -> toString() has '%2524'.
            // So we should NOT use encodeURIComponent if setting to .username, OR we should manually construct string.
            // Given the specific issue with $$, which URL object doesn't encode by default as it's allowed:
            // We must manually construct the string to guarantee encoding.

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
          on_demand: onDemand, // Note: payload uses snake_case in backend? Previous code had on_demand: onDemand.
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
      alert('Error al agregar la transmisión');
    }
  };

  return (
    <div className="p-8 flex justify-center">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Agregar Nueva Transmisión</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Nombre de la Transmisión</Label>
              <Input id="name" value={name} onChange={e => setName(e.target.value)} required placeholder="Ej. Cámara Entrada" />
            </div>

            <div className="space-y-2">
              <Label htmlFor="url">URL RTSP (Sin credenciales)</Label>
              <Input id="url" value={url} onChange={e => setUrl(e.target.value)} required placeholder="rtsp://192.168.1.50:554/live" />
              <p className="text-xs text-muted-foreground">Ingresa la URL sin usuario ni contraseña.</p>
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
            <Button type="submit" className="w-full">Agregar Transmisión</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
