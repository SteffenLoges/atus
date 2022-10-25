import { RouteRecordRaw } from "vue-router";
import sources from "./children/sources";
import fileservers from "./children/fileservers";
import users from "./children/users";
import filters from "./children/filters";

export default <RouteRecordRaw>{
  path: "/settings",
  name: "settings",
  meta: {
    title: "Settings",
  },
  component: () =>
    import(
      /* webpackChunkName: "settings" */ "@/views/Settings/Index.vue"
    ),
  children: [
    sources,
    fileservers,
    users,
    filters,
    {
      name: "settings_destinations",
      path: "destinations",
      meta: {
        title: "Upload Settings",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_destinations" */ "@/views/Settings/children/Destinations/Index.vue"
        ),
    },
    {
      name: "settings_samples",
      path: "samples",
      meta: {
        title: "Sample Settings",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_samples" */ "@/views/Settings/children/Samples/Index.vue"
        ),
    },
  ],
};
