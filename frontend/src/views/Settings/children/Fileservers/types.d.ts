interface IFileserver {
  uid: string;
  name: string;
  enabled: boolean;
  listInterval: number;
  minFreeDiskSpace: number;
  statisticsInterval: number;
  url: string;
  filesDownloaded: number;
  serverLoad?: number[];
  diskFreeSpace?: number;
  diskTotalSpace?: number;
}

interface IFileserverSettings {
  allocationMethod: "RANDOM" | "FILL" | "MOST_FREE";
  downloadLabel: string;
  uploadLabel: string;
}
