import { RouteRecordRaw } from "vue-router";

export default <RouteRecordRaw>{
  name: "settings_filters",
  alias: "",
  path: "filters",
  meta: {
    title: "Filters",
  },
  children: [
    {
      name: "settings_filters_categories",
      path: "categories",
      alias: "",
      meta: {
        title: "Categories",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_filters_categories" */ "@/views/Settings/children/Filters/Categories/Index.vue"
        ),
    },
    {
      name: "settings_filters_misc",
      path: "misc",
      meta: {
        title: "Miscellaneous",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_filters_misc" */ "@/views/Settings/children/Filters/Misc/Index.vue"
        ),
    },
  ],
};
