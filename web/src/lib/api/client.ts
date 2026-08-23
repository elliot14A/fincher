import { client } from './generated/client.gen'

client.setConfig({
  baseUrl: (import.meta.env.VITE_API_URL as string) || '/api',
})

export * from './generated'
export { client }
