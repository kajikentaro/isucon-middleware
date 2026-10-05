import react from "@vitejs/plugin-react";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  // relative paths, so the Go side can serve the built UI under any prefix
  base: "./",
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  build: {
    outDir: "out",
  },
  // listen on all interfaces so that the dev_front container (docker-compose.yaml) can publish the port
  server: {
    host: true,
  },
});
