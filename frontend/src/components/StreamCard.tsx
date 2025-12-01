import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Play, Video } from 'lucide-react';
import { Link } from 'react-router-dom';

export interface Stream {
  uuid: string;
  name: string;
  channels: Record<string, any>;
}

export function StreamCard({ stream, uuid, isAdmin, onDelete }: { stream: Stream; uuid: string; isAdmin?: boolean; onDelete?: () => void }) {
  const channelCount = stream.channels ? Object.keys(stream.channels).length : 0;

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium truncate" title={stream.name}>
          {stream.name}
        </CardTitle>
        <Badge variant="secondary">{channelCount} Channels</Badge>
      </CardHeader>
      <CardContent>
        <div className="aspect-video w-full bg-muted rounded-md flex items-center justify-center mb-4">
           <Video className="h-10 w-10 text-muted-foreground" />
        </div>
        <div className="flex justify-end gap-2">
           {isAdmin && (
             <>
               <Button size="sm" variant="outline" asChild>
                 <Link to={`/stream/${uuid}/edit`}>Edit</Link>
               </Button>
               <Button size="sm" variant="destructive" onClick={onDelete}>
                 Delete
               </Button>
             </>
           )}
           <Button size="sm" asChild>
             <Link to={`/stream/${uuid}`}>
               <Play className="mr-2 h-4 w-4" /> Watch
             </Link>
           </Button>
        </div>
      </CardContent>
    </Card>
  );
}
