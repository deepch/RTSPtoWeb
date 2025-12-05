import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Play, Video, Settings, Trash2 } from 'lucide-react';
import { Link } from 'react-router-dom';

export interface Stream {
  uuid: string;
  name: string;
  channels: Record<string, any>;
}

export function StreamCard({ stream, uuid, isAdmin, onDelete }: { stream: Stream; uuid: string; isAdmin?: boolean; onDelete?: () => void }) {
  const channelCount = stream.channels ? Object.keys(stream.channels).length : 0;

  return (
    <Card className="group overflow-hidden transition-all duration-300 hover:shadow-lg border-border/50 hover:border-primary/50">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 p-4 bg-muted/30">
        <div className="flex flex-col gap-1 min-w-0">
            <CardTitle className="text-base font-semibold truncate leading-none" title={stream.name}>
            {stream.name}
            </CardTitle>
            <div className="text-xs text-muted-foreground font-mono opacity-70 truncate">
                ID: {uuid}
            </div>
        </div>
        <Badge variant={channelCount > 0 ? "default" : "secondary"} className="ml-2 shrikn-0">
            {channelCount > 0 ? 'En Línea' : 'Desconectado'}
        </Badge>
      </CardHeader>

      <CardContent className="p-0">
        <div className="relative aspect-video w-full bg-gradient-to-br from-muted to-muted/50 flex items-center justify-center group-hover:scale-[1.02] transition-transform duration-500">
           <Video className="h-12 w-12 text-muted-foreground/50" />
           {channelCount > 0 && (
               <div className="absolute bottom-2 right-2">
                   <Badge variant="secondary" className="bg-background/80 backdrop-blur-sm text-xs">
                       {channelCount} CH
                   </Badge>
               </div>
           )}
        </div>
      </CardContent>

      <CardFooter className="p-4 flex justify-between items-center gap-2 bg-background">
           <Button size="sm" className="w-full font-medium" asChild>
             <Link to={`/stream/${uuid}`}>
               <Play className="mr-2 h-3.5 w-3.5" /> Ver
             </Link>
           </Button>

           {isAdmin && (
             <div className="flex gap-1 shrink-0">
               <Button size="icon" variant="ghost" className="h-9 w-9 text-muted-foreground hover:text-foreground" title="Editar Transmisión" asChild>
                 <Link to={`/stream/${uuid}/edit`}>
                    <Settings className="h-4 w-4" />
                 </Link>
               </Button>
               <Button size="icon" variant="ghost" className="h-9 w-9 text-muted-foreground hover:text-destructive" title="Eliminar Transmisión" onClick={onDelete}>
                 <Trash2 className="h-4 w-4" />
               </Button>
             </div>
           )}
      </CardFooter>
    </Card>
  );
}
