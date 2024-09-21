/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./template/**/*.html"],
  theme: {
    extend: {
      gridTemplateColumns: {
        'track': 'minmax(0, auto) repeat(31, minmax(0, 20px))',
      }
    },
  },
  plugins: [],
}

