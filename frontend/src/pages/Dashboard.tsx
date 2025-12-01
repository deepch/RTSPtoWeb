import { useEffect, useState } from 'react';
import client from '@/api/client';
import { StreamCard, type Stream } from '@/components/StreamCard';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';
import { Plus } from 'lucide-react';

export default function Dashboard() {
  const [streams, setStreams] = useState<Record<string, Stream>>({});
  const [loading, setLoading] = useState(true);
  const { user } = useAuth();

  const fetchStreams = async () => {
    try {
      const response = await client.get('/streams');
      if (response.data.status === 1) {
        setStreams(response.data.payload);
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
    if (!confirm('Are you sure you want to delete this stream?')) return;
    try {
      await client.get(`/stream/${uuid}/delete`);
      fetchStreams();
    } catch (error) {
      console.error('Failed to delete stream', error);
      alert('Failed to delete stream');
    }
  };

  if (loading) {
    return <div className="p-8">Loading streams...</div>;
  }

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <div className="flex items-center gap-4">
          <div className="text-sm text-muted-foreground">
            Logged in as {user?.username}
          </div>
          {user?.role === 'admin' && (
            <Button asChild>
              <Link to="/stream/add">
                <Plus className="mr-2 h-4 w-4" /> Add Stream
              </Link>
            </Button>
          )}
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {Object.entries(streams).map(([uuid, stream]) => (
          <StreamCard
            key={uuid}
            uuid={uuid}
            stream={stream}
            isAdmin={user?.role === 'admin'}
            onDelete={() => handleDelete(uuid)}
          />
        ))}
      </div>

      {Object.keys(streams).length === 0 && (
        <div className="text-center text-muted-foreground mt-12">
          No streams found.
        </div>
      )}
    </div>
  );
}
