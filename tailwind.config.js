/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./template/**/*.gotmpl", "./page/**/*.gotmpl"],
  theme: {
    extend: {
      colors: {
        eigengrau: {
          50: "#f6f6f9",
          100: "#ececf2",
          200: "#d6d6e1",
          300: "#b2b3c7",
          400: "#888ba8",
          500: "#696d8e",
          600: "#545675",
          700: "#444560",
          800: "#3b3c51",
          900: "#353545",
          950: "#16161d"
        }
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
