const defaultTheme = require("tailwindcss/defaultTheme");
const forms = require("@tailwindcss/forms");

const config = {
  mode: "jit",
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    fontFamily: {
      sans: ["Roboto", "Noto Sans TC", ...defaultTheme.fontFamily.sans],
    },
    extend: {},
  },

  plugins: [forms],
};

module.exports = config;
