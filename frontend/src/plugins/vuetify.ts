import "/src/styles/_overrides.scss";
import { createVuetify } from "vuetify";
import { aliases, mdi } from "vuetify/iconsets/mdi-svg";

export default createVuetify({
  treeshake: true,
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi,
    },
  },
  theme: {
    defaultTheme: "atus",
    themes: {
      atus: {
        dark: true,
        colors: {
          surface: "#1c212e",
          info: "#125a96",
          warning: "#bb7a2a",
        },
      },
    },
  },
});
