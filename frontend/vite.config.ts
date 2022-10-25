import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/vite";

// https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
import vuetify from "vite-plugin-vuetify";
import { fileURLToPath, URL } from "url";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vuetify({
      autoImport: true,
      // styles: "expose",
    }),
    Components(),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  base: "/",
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: ["./src/styles/variables"]
          .map((s) => `@use "${s}";`)
          .join("\n"),
      },
    },
  },
  build: {
    chunkSizeWarningLimit: 800,
  },
  server: {
    strictPort: true,
    port: 3000,
  },
});
