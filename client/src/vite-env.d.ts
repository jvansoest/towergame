/// <reference types="vite/client" />

declare module "*.glb" {
  const src: string;
  export default src;
}

interface ImportMetaEnv {
  readonly VITE_SERVER_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
