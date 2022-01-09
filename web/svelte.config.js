import preprocess from "svelte-preprocess";
import adapter from "@sveltejs/adapter-node";
import path from "path";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: [
    preprocess({
      postcss: true,
    }),
  ],
  kit: {
    target: "#root",
    adapter: adapter(),
    prerender: {
      enabled: false,
    },
    vite: {
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
    },
  },
};

export default config;
