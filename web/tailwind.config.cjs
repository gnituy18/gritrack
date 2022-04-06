const defaultTheme = require("tailwindcss/defaultTheme");
const forms = require("@tailwindcss/forms");

const config = {
  mode: "jit",
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    fontFamily: {
      sans: ["Roboto", "Noto Sans TC", ...defaultTheme.fontFamily.sans],
    },
    extend: {
      colors: {
        blue: {
          100: "#adcdff",
          200: "#82b1ff",
          300: "#5c9aff",
          400: "#3381ff",
          500: "#0a68ff",
          600: "#0056e0",
          700: "#0046b8",
          800: "#00378f",
          900: "#002766",
        },
      },
    },
  },

  plugins: [forms],
};

module.exports = config;
