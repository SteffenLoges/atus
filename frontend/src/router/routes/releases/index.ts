import { RouteRecordRaw } from "vue-router";

export default <RouteRecordRaw>{
  path: "/releases",
  alias: "/",
  meta: {
    title: "Releases",
  },
  children: [
    {
      name: "releases_browse",
      path: "browse",
      alias: "",
      meta: {
        title: "Browse Releases",
      },
      component: () =>
        import(
          /* webpackChunkName: "releases_browse" */ "@/views/Releases/Browse/Index.vue"
        ),
    },
    {
      name: "releases_details",
      path: "details/:uid/:name",
      meta: {
        title: "Release Details",
      },
      component: () =>
        import(
          /* webpackChunkName: "releases_details" */ "@/views/Releases/Details/Index.vue"
        ),
    },
  ],
};
