interface IFileserverStatistics {
  diskFreeSpace: number;
  diskTotalSpace: number;
  serverLoad: [number, number, number];
}

interface IFileserver {
  uid: string;
  name: string;
  enabled: boolean;
  statistics?: IFileserverStatistics;
}
