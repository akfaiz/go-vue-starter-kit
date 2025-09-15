// ui/src/types/global.d.ts
export {};

declare global {
  interface Window {
    __APP_CONFIG__?: {
      env: string;
      apiBaseUrl: string;
    };
  }
}
