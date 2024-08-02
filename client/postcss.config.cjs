const tailwindcss = require("tailwindcss");
const tailwindNesting = require("tailwindcss/nesting");
const postCssNesting = require("postcss-nesting");
const postCssImport = require("postcss-import");
const autoprefixer = require("autoprefixer");

/** @type {{ plugins: import('postcss').Plugin }} */
module.exports = {
  plugins: [
    postCssImport,
    tailwindNesting(postCssNesting),
    tailwindcss(),
    autoprefixer,
  ],
};
