interface ILogEntry {
  severity: number;
  type: string;
  added: string;
  message: string;
  releaseUID?: string;
  releaseName?: string;
}

interface ILogResponse {
  entries: ILogEntry[];
  types: string[];
  count: number;
}
