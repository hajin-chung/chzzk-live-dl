import { ft } from "./utils";

export type Channel = {
  channelId: string;
  channelName: string;
  channelImageUrl: string;
  openLive: string;
  followerCount: number;
};

export type Streamer = {
  id: string;
  name: string;
  image: string;
  followerCount: number;
  isLive: boolean;
  isDownloading: boolean;
  autoDownload: boolean;
}

export type Streamers = Record<string, Streamer>

export async function searchChannels(query: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/search?query=${query}`,
  );
  const data = (await res.json()) as Channel[];
  data.sort((a, b) => b.followerCount - a.followerCount);
  return data;
}

export async function addStreamer(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}`,
    { method: "POST" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function deleteStreamer(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}`,
    { method: "DELETE" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function startDownload(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}/download`,
    { method: "POST" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function stopDownload(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}/download`,
    { method: "DELETE" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function enableAutoDownload(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}/autoDownload`,
    { method: "POST" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function disableAutoDownload(id: string) {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer/${id}/autoDownload`,
    { method: "DELETE" },
  );
  const data = (await res.json()) as Streamers;
  return data;
}

export async function getStreamers() {
  const res = await ft(
    `${import.meta.env.VITE_API_ENDPOINT}/api/streamer`,
  );
  const data = (await res.json()) as Streamers;
  return data;
}
