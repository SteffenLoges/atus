import { RouteRecordRaw } from "vue-router";

export default <RouteRecordRaw>{
  name: "settings_sources",
  alias: "",
  path: "sources",
  meta: {
    title: "Sources",
  },
  children: [
    {
      name: "settings_sources_manage",
      path: "manage",
      alias: "",
      meta: {
        title: "Manage Sources",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_sources_manage" */ "@/views/Settings/children/Sources/Manage/Index.vue"
        ),
    },
    {
      name: "settings_sources_add",
      path: "add",
      meta: {
        title: "Add Source",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_sources_add" */ "@/views/Settings/children/Sources/Add.vue"
        ),
    },
    {
      name: "settings_sources_edit",
      path: "edit/:uid",
      meta: {
        title: "Edit Source",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_sources_edit" */ "@/views/Settings/children/Sources/Edit.vue"
        ),
    },
  ],
};
