const defaultTheme = require("tailwindcss/defaultTheme");

const config = {
  mode: "jit",
  purge: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    fontFamily: {
      sans: ["Roboto", "Noto Sans TC", ...defaultTheme.fontFamily.sans],
    },
    extend: {},
  },

  plugins: [],
};

module.exports = config;
