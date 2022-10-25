type IMetaFileType =
  | "TORRENT"
  | "NFO"
  | "SOURCE_IMAGE"
  | "IMAGE"
  | "PROOF_IMAGE"
  | "SCREEN_IMAGE"
  | "SCREEN_IMAGE__FROM_SAMPLE"
  | "SAMPLE_VIDEO";

type IMetaFileState =
  | "UNKNOWN"
  | "DOWNLOADING"
  | "PROCESSED"
  | "ERROR";

interface IMetaFile {
  releaseUID: string;
  type: IMetaFileType;
  state: IMetaFileState;
  fileName: string;
  info: any[] | null;
}

interface IRelease {
  uid: string;
  name: string;
  nameRaw: string;
  pre: string;
  category: string;
  categoryRaw: string;
  size: number;
  added: string;
  fileserverName: string;
  sourceName: string;
  metaFiles: IMetaFile[];
  state: IReleaseState;
  downloadState?: IDownloadState;
}

interface IDownloadState {
  hash: string;
  downloadRate: number;
  eta: number;
  done: number;
  state:
    | "STARTED"
    | "PAUSED"
    | "STOPPED"
    | "HASHING"
    | "CHECKING"
    | "ERROR";
}

interface IReleaseState {
  state:
    | "NEW"
    | "DOWNLOAD_INIT"
    | "DOWNLOADING"
    | "DOWNLOADED"
    | "UPLOADED"
    | "GENERAL_ERROR"
    | "UPLOAD_ERROR";
  uploadDate: string;
}
