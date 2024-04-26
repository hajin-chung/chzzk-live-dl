import { useEffect, useState } from "react";
import { Layout } from "./Layout";
import { Search } from "./Search";
import {
  Streamers,
  addStreamer,
  deleteStreamer,
  disableAutoDownload,
  enableAutoDownload,
  getStreamers,
  startDownload,
  stopDownload,
} from "./lib/api.ts";
import { Button } from "./components/ui/button.tsx";
import {
  DeleteIcon,
  DownloadIcon,
  RefreshIcon,
  StopIcon,
} from "./components/icons.tsx";
import { Credentials } from "./Credentials.tsx";

function App() {
  const { streamers, updateStreamers, setStreamers } = useStreamers();

  function onAdd(id: string) {
    addStreamer(id).then(setStreamers);
  }

  function onDelete(id: string) {
    deleteStreamer(id).then(setStreamers);
  }

  function onSetAutoDownload(id: string, autoDownload: boolean) {
    if (autoDownload) enableAutoDownload(id).then(setStreamers);
    else disableAutoDownload(id).then(setStreamers);
  }

  function onSetDownload(id: string, download: boolean) {
    if (download) startDownload(id).then(setStreamers);
    else stopDownload(id).then(setStreamers);
  }

  return (
    <Layout>
      <div className="w-full flex justify-between">
        <Button variant="ghost" onClick={updateStreamers}>
          <RefreshIcon />
        </Button>
        <div className="flex gap-3">
          <Credentials />
          <Search onAdd={onAdd} streamers={streamers} />
        </div>
      </div>
      <div className="h-4" />
      <div className="flex flex-col gap-2">
        {Object.values(streamers).map(
          ({
            id,
            name,
            image,
            isLive,
            isDownloading,
            followerCount,
            autoDownload,
          }) => (
            <div className="flex items-center w-full justify-between" key={id}>
              <div className="flex gap-3 items-center">
                <img
                  className={`w-7 h-7 rounded-full object-cover p-0.5 border-2 ${isLive ? "border-red-400" : "border-neutral-400"}`}
                  src={image}
                />
                <div className="flex flex-col">
                  <p>{name}</p>
                  <p className="text-sm">{followerCount}</p>
                </div>
              </div>
              <div className="flex gap-3 items-center">
                <Button
                  variant="outline"
                  className={`${autoDownload && "bg-neutral-700"}`}
                  onClick={() => onSetAutoDownload(id, !autoDownload)}
                >
                  <RefreshIcon />
                </Button>
                <Button
                  variant="outline"
                  onClick={() => onSetDownload(id, !isDownloading)}
                >
                  {isDownloading ? <StopIcon /> : <DownloadIcon />}
                </Button>
                <Button variant="outline" onClick={() => onDelete(id)}>
                  <DeleteIcon />
                </Button>
              </div>
            </div>
          ),
        )}
      </div>
    </Layout>
  );
}

function useStreamers() {
  const [streamers, setStreamers] = useState<Streamers>({});

  const updateStreamers = () => {
    getStreamers().then(setStreamers);
  };

  useEffect(() => {
    updateStreamers();
    const interval = setInterval(updateStreamers, 30 * 1000);
    return () => clearInterval(interval);
  }, []);

  return { streamers, setStreamers, updateStreamers };
}

export default App;
