export const env = {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? '/api',
  apiToken: import.meta.env.VITE_API_TOKEN ?? '',
} as const
