let plugin = require("tailwindcss/plugin");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.templ"],
  theme: {
    colors: {
      base: "#faf4ed",
      surface: "#fffaf3",
      overlay: "#f2e9e1",
      muted: "#9893a5",
      subtle: "#797593",
      fg: "#575279",
      love: "#b4637a",
      gold: "#ea9d34",
      rose: "#d7827e",
      pine: "#286983",
      foam: "#56949f",
      iris: "#907aa9",
      highlightLow: "#f4ede8",
      highlightMed: "#dfdad9",
      highlightHigh: "#cecacd",
    },
  },
  plugins: [
    plugin(function ({ addBase, addUtilities, addVariant }) {
      addBase({
        "*": { minWidth: "0" },
      });
    }),
  ],
}

