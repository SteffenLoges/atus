interface ISource {
  uid: string;
  name: string;
  favicon: string;
  enabled: boolean;
  timesChecked: number;
  lastChecked: string;
  nextCheck: string;
  sumTorrentsDownloaded: number;
  sumImagesDownloaded: number;
  sumReleasesDownloaded: number;
  metaPath: string;
  metaPathUseAsKey: boolean;
  metaPathAutoDetected: boolean;
  imagePath: string;
  imagePathUseAsKey: boolean;
  imagePathAutoDetected: boolean;
  rssURL: string;
  requiresCookies: boolean;
  cookies: ICookie[];
  rssInterval: number;
  requestWaitTime: number;
}

interface ICookie {
  name: string;
  value: string;
}
