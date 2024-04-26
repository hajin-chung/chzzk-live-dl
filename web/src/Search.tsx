import { useEffect } from "react";
import { useState } from "react";
import { Input } from "./components/ui/input";
import { Channel, Streamers, searchChannels } from "@/lib/api";
import { Button } from "./components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "./components/ui/dialog";
import { SearchIcon } from "./components/icons";

type SearchProps = {
  onAdd: (id: string) => void;
  streamers: Streamers;
};

export function Search({ onAdd, streamers }: SearchProps) {
  const [query, setQuery] = useState("");
  const [channels, setChannels] = useState<Channel[]>([]);

  useEffect(() => {
    if (query.trim().length == 0) return;
    const timeout = setTimeout(async () => {
      const resultChannels = await searchChannels(query.trim());
      setChannels(resultChannels);
    }, 500);
    return () => clearTimeout(timeout);
  }, [query]);

  function handleAdd(id: string) {
    setQuery("");
    onAdd(id);
  }

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="outline">
          <SearchIcon />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <div className="flex flex-col gap-4 h-[296px]">
          <Input value={query} onChange={(e) => setQuery(e.target.value)} />
          {channels.length > 0 && query.length > 0 && (
            <div className="w-full max-h-[240px] overflow-y-scroll p-4">
              <div className="flex flex-col gap-2 w-full">
                {channels.map(
                  ({
                    channelId,
                    channelName,
                    channelImageUrl,
                    openLive,
                    followerCount,
                  }) => (
                    <div
                      className="flex items-center w-full justify-between"
                      key={channelId}
                    >
                      <div className="flex gap-3 items-center">
                        <img
                          className={`w-7 h-7 rounded-full object-cover p-0.5 border-2 ${openLive ? "border-red-400" : "border-neutral-400"}`}
                          src={channelImageUrl}
                        />
                        <div className="flex flex-col">
                          <p>{channelName}</p>
                          <p className="text-sm">{followerCount}</p>
                        </div>
                      </div>
                      <Button
                        variant="outline"
                        disabled={!!streamers[channelId]}
                        onClick={() => handleAdd(channelId)}
                        className="justify-self-end"
                      >
                        Add
                      </Button>
                    </div>
                  ),
                )}
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
