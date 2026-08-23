import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../openapi/swagger.json',
  output: 'src/lib/api/generated',
  plugins: [
    '@hey-api/client-fetch',
    '@hey-api/typescript',
    '@hey-api/sdk',
    'valibot',
  ],
})
