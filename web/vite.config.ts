import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'
import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

export default defineConfig({
  plugins: [
    TanStackRouterVite({
      target: 'react',
      autoCodeSplitting: true,
    }),
    preact(),
    vanillaExtractPlugin(),
  ],
  resolve: {
    tsconfigPaths: true,
    alias: {
      react: 'preact/compat',
      'react-dom/test-utils': 'preact/test-utils',
      'react-dom': 'preact/compat',
      'react/jsx-runtime': 'preact/jsx-runtime',
    },
  },
  build: {
    outDir: '../pkg/web/dist',
    emptyOutDir: true,
    rollupOptions: {
      onwarn(warning, warn) {
        if (
          warning.code === 'IMPORT_IS_UNDEFINED' &&
          (warning.message?.includes('preact/compat') ||
            warning.id?.includes('@tanstack/react-router'))
        ) {
          return
        }
        warn(warning)
      },
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
