import { useEffect, useState, useMemo } from 'react';
import client from '@/api/client';
import { StreamCard, type Stream } from '@/components/StreamCard';
import { useAuth } from '@/context/AuthContext';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Link, useNavigate } from 'react-router-dom';
import { Plus, Search, Activity, Video } from 'lucide-react';

export default function Dashboard() {
  const [streams, setStreams] = useState<Record<string, Stream>>({});
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect non-admin users to MultiView
    if (user && user.role !== 'admin') {
      navigate('/multiview');
    }
  }, [user, navigate]);

  const fetchStreams = async () => {
    try {
      const response = await client.get('/streams');
      if (response.data.status === 1) {
        // Ensure payload is an object, defaulting to {} if null/undefined
        setStreams(response.data.payload || {});
      }
    } catch (error) {
      console.error('Failed to fetch streams', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStreams();
  }, []);

  const handleDelete = async (uuid: string) => {
    if (!confirm('¿Estás seguro de que quieres eliminar esta transmisión?')) return;
    try {
      await client.get(`/stream/${uuid}/delete`);
      fetchStreams();
    } catch (error) {
      console.error('Failed to delete stream', error);
      alert('Error al eliminar la transmisión');
    }
  };

  const filteredStreams = useMemo(() => {
    return Object.entries(streams).filter(([_, stream]) =>
      stream.name.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [streams, searchQuery]);

  const stats = useMemo(() => {
    const totalStreams = Object.keys(streams).length;
    const totalChannels = Object.values(streams).reduce((acc, stream) =>
      acc + (stream.channels ? Object.keys(stream.channels).length : 0), 0
    );
    return { totalStreams, totalChannels };
  }, [streams]);

  if (loading) {
     // Elegant loading state
    return (
        <div className="p-8 space-y-8 animate-pulse">
            <div className="h-8 w-48 bg-muted rounded"></div>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <div className="h-32 bg-muted rounded-xl"></div>
                <div className="h-32 bg-muted rounded-xl"></div>
            </div>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {[1, 2, 3, 4].map(i => (
                    <div key={i} className="h-64 bg-muted rounded-xl"></div>
                ))}
            </div>
        </div>
    )
  }

  // If redirected, might briefly show nothing or loading, but logic above handles it.
  // We check again for safety to render null if not admin (though useEffect handles redirect)
  if (user?.role !== 'admin') return null;

  return (
    <div className="p-8 space-y-8 max-w-[1600px] mx-auto">
      {/* Header & Stats */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div>
          <h1 className="text-4xl font-bold tracking-tight bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            Panel de Administración
          </h1>
          <p className="text-muted-foreground mt-2">
            Gestiona tus transmisiones y configuraciones RTSP.
          </p>
        </div>

        <div className="flex gap-4">
            <Card className="min-w-[140px]">
                <CardContent className="p-4 flex flex-col justify-between h-full">
                    <div className="flex justify-between items-start">
                        <span className="text-sm text-muted-foreground font-medium">Transmisiones</span>
                        <Video className="h-4 w-4 text-primary" />
                    </div>
                    <div className="text-2xl font-bold">{stats.totalStreams}</div>
                </CardContent>
            </Card>
            <Card className="min-w-[140px]">
                <CardContent className="p-4 flex flex-col justify-between h-full">
                    <div className="flex justify-between items-start">
                        <span className="text-sm text-muted-foreground font-medium">Canales Activos</span>
                        <Activity className="h-4 w-4 text-primary" />
                    </div>
                    <div className="text-2xl font-bold">{stats.totalChannels}</div>
                </CardContent>
            </Card>
        </div>
      </div>

      {/* Action Bar */}
      <div className="flex items-center space-x-2">
         <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
                placeholder="Buscar transmisiones..."
                className="pl-9 bg-background/50 backdrop-blur-sm"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
            />
         </div>
      </div>

      {/* Content Grid */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {/* Add New Stream Card */}
        <Link to="/stream/add" className="group">
            <div className="h-full min-h-[280px] border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 hover:bg-muted/30 rounded-xl flex flex-col items-center justify-center gap-4 transition-all duration-300 cursor-pointer">
                <div className="h-16 w-16 rounded-full bg-muted group-hover:bg-primary/10 flex items-center justify-center transition-colors">
                    <Plus className="h-8 w-8 text-muted-foreground group-hover:text-primary transition-colors" />
                </div>
                <div className="text-center">
                    <h3 className="font-semibold text-lg">Agregar Transmisión</h3>
                    <p className="text-sm text-muted-foreground">Configurar nueva fuente RTSP</p>
                </div>
            </div>
        </Link>

        {/* Stream Cards */}
        {filteredStreams.map(([uuid, stream]) => (
          <StreamCard
            key={uuid}
            uuid={uuid}
            stream={stream}
            isAdmin={true}
            onDelete={() => handleDelete(uuid)}
          />
        ))}
      </div>

      {filteredStreams.length === 0 && searchQuery && (
        <div className="text-center py-20 text-muted-foreground">
            No se encontraron transmisiones que coincidan con "{searchQuery}"
        </div>
      )}
    </div>
  );
}
