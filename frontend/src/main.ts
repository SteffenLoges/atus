import { createApp } from "vue";
import { createPinia } from "pinia";
import { createHead } from "@vueuse/head";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./router";
import { getTitleTemplate } from "./router/meta";
import "typeface-roboto";
import "./plugins/moment";

const app = createApp(App);
const pinia = createPinia();

const head = createHead({
  titleTemplate: getTitleTemplate,
});

app.use(router);
app.use(vuetify);
app.use(pinia);
app.use(head);

app.mount("#app");
