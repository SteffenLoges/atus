import {
  createRouter,
  createWebHistory,
  createWebHashHistory,
} from "vue-router";
import { getTitleTemplate } from "./meta";
import settings from "./routes/settings";
import releases from "./routes/releases";
import useGlobalStore from "@/store/global";

const router = createRouter({
  history:
    import.meta.env.MODE === "production"
      ? createWebHistory()
      : createWebHashHistory(),
  routes: [
    releases,
    settings,
    {
      path: "/login",
      name: "login",
      meta: {
        title: "Login",
      },
      component: () =>
        import(
          /* webpackChunkName: "login" */ "../views/Login/Index.vue"
        ),
    },
    {
      path: "/setup",
      name: "setup",
      meta: {
        title: "Getting Started",
      },
      component: () =>
        import(
          /* webpackChunkName: "setup" */ "../views/Setup/Index.vue"
        ),
    },
    {
      path: "/register",
      name: "register",
      meta: {
        title: "Welcome",
      },
      component: () =>
        import(
          /* webpackChunkName: "register" */ "../views/Register/Index.vue"
        ),
    },
    {
      path: "/api",
      name: "api",
      meta: {
        title: "API",
      },
      component: () =>
        import(
          /* webpackChunkName: "api" */ "../views/API/Index.vue"
        ),
    },
    {
      path: "/log",
      name: "log",
      meta: {
        title: "Log",
      },
      component: () =>
        import(
          /* webpackChunkName: "log" */ "../views/Log/Index.vue"
        ),
    },
    {
      path: "/debug",
      name: "debug",
      meta: {
        title: "Debug",
      },
      component: () =>
        import(
          /* webpackChunkName: "debug" */ "../views/Debug/Index.vue"
        ),
    },
    {
      path: "/donate",
      name: "donate",
      meta: {
        title: "Donate",
      },
      component: () =>
        import(
          /* webpackChunkName: "donate" */ "../views/Donate/Index.vue"
        ),
    },
    {
      path: "/404",
      meta: {
        title: "Page not found",
      },
      component: () =>
        import(
          /* webpackChunkName: "error_404" */ "../views/404/Index.vue"
        ),
    },
    { path: "/:catchAll(.*)", redirect: "/404" },
  ],
});

router.beforeEach((to, from, next) => {
  // close drawer
  const globalStore = useGlobalStore();
  globalStore.setForceDrawer(false);

  // scroll to top
  window.scrollTo(0, 0);

  next();
});

router.afterEach((to) => {
  // set title
  document.title = getTitleTemplate(
    <string | undefined>to.meta?.title
  );
});

export default router;
