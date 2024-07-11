import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { crx, defineManifest } from "@crxjs/vite-plugin";
//import { viteStaticCopy } from "vite-plugin-static-copy"; // ↓ この行を追加

const manifest = defineManifest({
  manifest_version: 3,
  name: "OpenAltBrowser",
  description:
    "OpenAltBrowser allows you to effortlessly open the current tab in another browser. Enhance your browsing experience by seamlessly switching between browsers with just a click. Perfect for developers, testers, and anyone who uses multiple web browsers",
  version: "1.0.0",
  action: {
    //default_popup: "src/popup/index.html",
    //default_icon:  "",
    default_title: "Open in Another Browser",
  },
  options_ui: {
    page: "src/options/index.html",
  },
  //content_scripts: [
  //  {
  //    matches: ["http://*/*", "https://*/*"],
  //    js: ["src/contentscript/index.ts"],
  //    run_at: "document_start",
  //  },
  //],
  //side_panel: {
  //  default_path: "src/sidepanel/index.html",
  //},
  background: {
    service_worker: "src/background/index.ts",
  },
  // host_permissions: ["<all_urls>"],
  permissions: [
    "nativeMessaging",
    "tabs",
    //  "storage",
    //  "sidePanel",
    //  "contextMenus",
    //  "activeTab",
    //  "scripting",
    //  "unlimitedStorage",
    //  "alarms",
  ],
});

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    /*
    viteStaticCopy({
      targets: [
        {
          src: "node_modules/bootstrap/dist/js/bootstrap.bundle.min.js",
          dest: ".",
        },
        {
          src: "node_modules/bootstrap/dist/css/bootstrap.min.css",
          dest: ".",
        },
      ],
    }),
    */
    svelte(),
    crx({ manifest }),
  ],
});
