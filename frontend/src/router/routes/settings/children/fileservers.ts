import { RouteRecordRaw } from "vue-router";

export default <RouteRecordRaw>{
  name: "settings_fileservers",
  path: "fileservers",
  meta: {
    title: "Fileservers",
  },
  children: [
    {
      name: "settings_fileservers_manage",
      path: "manage",
      alias: "",
      meta: {
        title: "Manage Fileservers",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_fileservers_manage" */ "@/views/Settings/children/Fileservers/Manage/Index.vue"
        ),
    },
    {
      name: "settings_fileservers_add",
      path: "add",
      meta: {
        title: "Add Fileserver",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_fileservers_add" */ "@/views/Settings/children/Fileservers/Add.vue"
        ),
    },
    {
      name: "settings_fileservers_edit",
      path: "edit/:uid",
      meta: {
        title: "Edit Fileserver",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_fileservers_edit" */ "@/views/Settings/children/Fileservers/Edit.vue"
        ),
    },
    {
      name: "settings_fileservers_settings",
      path: "settings",
      meta: {
        title: "Fileserver Settings",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_fileservers_settings" */ "@/views/Settings/children/Fileservers/Settings.vue"
        ),
    },
  ],
};
