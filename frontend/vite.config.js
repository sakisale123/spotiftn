import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

import fs from 'fs';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    watch: {
      usePolling: true,
    },
    host: true,
    strictPort: true,
    port: 5173,
    https: {
      key: fs.readFileSync('/etc/ssl/certs/server.key'),
      cert: fs.readFileSync('/etc/ssl/certs/server.crt'),
    },
  }
})
