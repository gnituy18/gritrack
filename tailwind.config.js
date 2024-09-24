/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./template/**/*.gotmpl", "./page/**/*.gotmpl"],
  theme: {
    extend: {
      gridTemplateColumns: {
        track: "minmax(0, 1fr) repeat(31, minmax(0, 20px))",
      },
    },
  },
  plugins: [],
};
