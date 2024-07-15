import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { crx, defineManifest } from "@crxjs/vite-plugin";

const manifest = defineManifest({
  manifest_version: 3,
  name: "OpenAltBrowser",
  description:
    "OpenAltBrowser allows you to effortlessly open the current tab in another browser. Enhance your browsing experience by seamlessly switching between browsers with just a click. Perfect for developers, testers, and anyone who uses multiple web browsers",
  version: "1.0.0",
  action: {
    //default_icon:  "",
    default_title: "Open in Another Browser",
  },
  options_page: "src/options/index.html",
  background: {
    service_worker: "src/background/index.ts",
  },
  permissions: [
    "nativeMessaging",
    "tabs",
    "storage",
    "activeTab",
    "contextMenus",
  ],
});

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte(), crx({ manifest })],
  server: {
    port: 5173,
    strictPort: true,
    hmr: {
      port: 5173,
    },
  },
});
