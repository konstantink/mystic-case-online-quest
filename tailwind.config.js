/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/templates/**/*.html",
    "./node_modules/flowbite/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        "mc-yellow": "rgb(255, 214, 68)",
        "mc-purple": "#3A3185",
        "mc-light-purple": "#938CD1",
        "mc-light-green": "#B3D138",
        "mc-light-black": "rgba(0, 0, 0, 0.7)",
        "mc-dark-purple": "#231e52",
      },
      container: {
        center: true,
      },
    },
  },
  plugins: [],
}

