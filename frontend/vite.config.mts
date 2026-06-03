import react from '@vitejs/plugin-react-swc'
import ViteYaml from '@modyfi/vite-plugin-yaml'
import { defineConfig } from 'vite'

export default defineConfig(() => {
  return {
    plugins: [
      react(),
      ViteYaml(),
    ],
    base: '',
    build: {
      emptyOutDir: true,
    }
  }
})
