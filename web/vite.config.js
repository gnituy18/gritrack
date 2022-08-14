import { sveltekit } from "@sveltejs/kit/vite";
import path from "path";

/** @type {import('vite').UserConfig} */
const config = {
  plugins: [sveltekit()],
  envPrefix: "ENV_",
  resolve: {
    alias: {
      $: path.resolve("./src"),
      $types: path.resolve("./src/types"),
      $components: path.resolve("./src/components"),
      $stores: path.resolve("./src/stores"),
      $apis: path.resolve("./src/apis"),
    },
  },
};

export default config;
