/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./src/*.{css,js,ts}",
    "./src/**/*.{css,js,ts}",
    "../server/*.{go,templ}",
    "../server/**/*.{go,templ}",
  ],
  darkMode: "media",
  theme: {
    extend: {
      colors({ colors }) {
        return {
          background: variables("background", colors.gray),
          "on-background": variables("on-background", {
            50: colors.black,
            100: colors.black,
            200: colors.black,
            300: colors.black,
            400: colors.black,
            500: colors.white,
            600: colors.white,
            700: colors.white,
            800: colors.white,
            900: colors.white,
            950: colors.white,
          }),

          primary: variables("primary", colors.indigo),
          "on-primary": variables("on-primary", {
            50: colors.black,
            100: colors.black,
            200: colors.black,
            300: colors.black,
            400: colors.black,
            500: colors.white,
            600: colors.white,
            700: colors.white,
            800: colors.white,
            900: colors.white,
            950: colors.white,
          }),

          secondary: variables("secondary", colors.emerald),

          "dark-background": {
            50: colors.gray[950],
            100: colors.gray[900],
            200: colors.gray[800],
            300: colors.gray[700],
            400: colors.gray[600],
            500: colors.gray[500],
            600: colors.gray[400],
            700: colors.gray[300],
            800: colors.gray[200],
            900: colors.gray[100],
            950: colors.gray[50],
          },
          "on-dark-background": {
            50: colors.white,
            100: colors.white,
            200: colors.white,
            300: colors.white,
            400: colors.white,
            500: colors.white,
            600: colors.black,
            700: colors.black,
            800: colors.black,
            900: colors.black,
            950: colors.black,
          },
          "dark-primary": {
            50: colors.indigo[950],
            100: colors.indigo[900],
            200: colors.indigo[800],
            300: colors.indigo[700],
            400: colors.indigo[600],
            500: colors.indigo[500],
            600: colors.indigo[400],
            700: colors.indigo[300],
            800: colors.indigo[200],
            900: colors.indigo[100],
            950: colors.indigo[50],
          },
          "on-dark-primary": {
            50: colors.white,
            100: colors.white,
            200: colors.white,
            300: colors.white,
            400: colors.white,
            500: colors.white,
            600: colors.black,
            700: colors.black,
            800: colors.black,
            900: colors.black,
            950: colors.black,
          },
        };
      },
    },
  },
  plugins: [],
};

/**
 *
 * @param {string} name
 * @param {string} defaultValue
 */
function variable(name, defaultValue) {
  return `var(--${name}, ${defaultValue})`;
}

/**
 * @param {string} prefix
 * @param {Record<number, string>} colors
 * @returns {Record<number, string>}
 */
function variables(prefix, colors) {
  return Object.fromEntries(
    Object.entries(colors).map(([key, value]) => [
      key,
      variable(`${prefix}-${key}`, value),
    ])
  );
}
