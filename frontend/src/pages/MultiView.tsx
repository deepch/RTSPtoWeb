import { useState, useEffect } from 'react';
import { Player } from '@/components/Player';
import { Button } from '@/components/ui/button';
import { Plus, X, LayoutGrid } from 'lucide-react';
import client from '@/api/client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

interface Stream {
    name: string;
    channels: Record<string, any>;
}

export default function MultiView() {
  const [layout, setLayout] = useState<number>(4); // 4 for 2x2, 9 for 3x3, 16 for 4x4
  const [slots, setSlots] = useState<Record<number, string>>({});
  const [streams, setStreams] = useState<Record<string, Stream>>({});
  const [selectedSlot, setSelectedSlot] = useState<number | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  useEffect(() => {
    // Load persisted state
    const savedLayout = localStorage.getItem('multiview_layout');
    const savedSlots = localStorage.getItem('multiview_slots');
    if (savedLayout) setLayout(parseInt(savedLayout));
    if (savedSlots) setSlots(JSON.parse(savedSlots));

    fetchStreams();
  }, []);

  const fetchStreams = async () => {
      try {
          const response = await client.get('/streams');
          if (response.data.status === 1) {
              setStreams(response.data.payload);
          }
      } catch (error) {
          console.error("Failed to fetch streams", error);
      }
  }

  const handleLayoutChange = (count: string) => {
    const newLayout = parseInt(count);
    setLayout(newLayout);
    localStorage.setItem('multiview_layout', count);
  };

  const addStreamToSlot = (uuid: string) => {
      if (selectedSlot !== null) {
          const newSlots = { ...slots, [selectedSlot]: uuid };
          setSlots(newSlots);
          localStorage.setItem('multiview_slots', JSON.stringify(newSlots));
          setIsDialogOpen(false);
          setSelectedSlot(null);
      }
  };

  const removeStreamFromSlot = (index: number) => {
      const newSlots = { ...slots };
      delete newSlots[index];
      setSlots(newSlots);
      localStorage.setItem('multiview_slots', JSON.stringify(newSlots));
  };

  const getGridClass = () => {
    switch (layout) {
      case 1: return 'grid-cols-1';
      case 4: return 'grid-cols-2';
      case 9: return 'grid-cols-3';
      case 16: return 'grid-cols-4';
      default: return 'grid-cols-2';
    }
  };

  return (
    <div className="p-4 h-[calc(100vh-4rem)] flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold flex items-center gap-2">
            <LayoutGrid className="h-6 w-6" /> MultiView
        </h1>
        <div className="flex items-center gap-4">
          <Select value={layout.toString()} onValueChange={handleLayoutChange}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select Layout" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">1x1 (Single)</SelectItem>
              <SelectItem value="4">2x2 (4 views)</SelectItem>
              <SelectItem value="9">3x3 (9 views)</SelectItem>
              <SelectItem value="16">4x4 (16 views)</SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" onClick={() => {
              setSlots({});
              localStorage.removeItem('multiview_slots');
          }}>Clear All</Button>
        </div>
      </div>

      <div className={`grid ${getGridClass()} gap-2 flex-1 min-h-0`}>
        {Array.from({ length: layout }).map((_, index) => (
          <div key={index} className="relative bg-muted/50 rounded-lg overflow-hidden border border-border flex items-center justify-center">
            {slots[index] ? (
              <>
                <Player uuid={slots[index]} />
                <Button
                  variant="destructive"
                  size="icon"
                  className="absolute top-2 right-2 h-6 w-6 opacity-0 hover:opacity-100 transition-opacity"
                  onClick={() => removeStreamFromSlot(index)}
                >
                  <X className="h-4 w-4" />
                </Button>
              </>
            ) : (
                <Dialog open={isDialogOpen && selectedSlot === index} onOpenChange={(open) => {
                    if (open) setSelectedSlot(index);
                    setIsDialogOpen(open);
                }}>
                    <DialogTrigger asChild>
                        <Button variant="ghost" className="h-full w-full flex flex-col gap-2 hover:bg-muted/80">
                            <Plus className="h-8 w-8 text-muted-foreground" />
                            <span className="text-muted-foreground">Add Camera</span>
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Select Stream</DialogTitle>
                        </DialogHeader>
                        <div className="grid gap-2">
                            {Object.entries(streams).map(([uuid, stream]) => (
                                <Button key={uuid} variant="outline" className="justify-start" onClick={() => addStreamToSlot(uuid)}>
                                    {stream.name || uuid}
                                </Button>
                            ))}
                            {Object.keys(streams).length === 0 && <div className="text-center text-muted-foreground">No streams available</div>}
                        </div>
                    </DialogContent>
                </Dialog>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
