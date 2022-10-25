import { RouteRecordRaw } from "vue-router";

export default <RouteRecordRaw>{
  name: "settings_users",
  alias: "",
  path: "users",
  meta: {
    title: "Users",
  },
  children: [
    {
      name: "settings_users_manage",
      path: "manage",
      alias: "",
      meta: {
        title: "Manage Users",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_users_manage" */ "@/views/Settings/children/Users/Manage/Index.vue"
        ),
    },
    {
      name: "settings_users_add",
      path: "add",
      meta: {
        title: "Add User",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_users_add" */ "@/views/Settings/children/Users/Add.vue"
        ),
    },
    {
      name: "settings_users_edit",
      path: "edit/:uid",
      meta: {
        title: "Add User",
      },
      component: () =>
        import(
          /* webpackChunkName: "settings_users_edit" */ "@/views/Settings/children/Users/Edit.vue"
        ),
    },
  ],
};
