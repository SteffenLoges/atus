/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_NAME: string;
  readonly VITE_APP_BACKEND_HOST: string;
  readonly VITE_APP_BACKEND_PORT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
