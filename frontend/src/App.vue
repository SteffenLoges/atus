<template>
  <v-app fullHeight>
    <ErrorOverlay v-if="error !== null && error.critical" :error="error.err" />
    <SplashScreen v-else-if="!authStatusLoaded" />
    <WebsocketDisconnected v-else-if="isAuthenticated && !websocketConnected" />
    <template v-else>
      <template v-if="isAuthenticated">
        <Drawer />
        <AppBar />
      </template>

      <v-main>
        <ErrorDialog v-if="error !== null" :text="error.err" @dismiss="clearError()" max-width="85vw" />

        <suspense>
          <router-view />

          <template #fallback>
            <v-row align="center" justify="center" class="fill-height">
              <v-progress-circular indeterminate color="primary" size="64" />
            </v-row>
          </template>
        </suspense>
      </v-main>
    </template>
  </v-app>
</template> 

<script lang="ts">
import { defineComponent, ref, watch, onBeforeUnmount, onErrorCaptured } from "vue";
import { useRouter, useRoute, RouteLocationNormalized, isNavigationFailure } from "vue-router";
import { storeToRefs } from "pinia";
import useGlobalStore from "@/store/global";
import useUserStore from "@/store/user";
import AppBar from "@/views/components/layout/AppBar/Index.vue";
import Drawer from "@/views/components/layout/Drawer/Index.vue";
import { isSSL, host, port } from "@/utils/url";
import websocket from "@/utils/websocket";

export default defineComponent({
  components: {
    AppBar,
    Drawer,
  },
  setup() {
    const router = useRouter();
    const route = useRoute();
    const global = useGlobalStore();
    const { error } = storeToRefs(global);
    const userStore = useUserStore();
    const { isAuthenticated } = storeToRefs(userStore);
    const authStatusLoaded = ref(false);

    onErrorCaptured((err: IResponse<string>, vm, info) => {
      global.setError(err, true);
      console.error(err, vm, info);
      return false;
    });

    // handle websocket connection
    const websocketConnected = ref(false);

    const sendCurrentPage = ({ name, params }: RouteLocationNormalized) =>
      websocket
        .send("CURRENT_PAGE", { name, params })
        .catch((e) => console.error(e));

    const wsConnect = (token: string) => websocket.connect({
      url: `ws${isSSL ? "s" : ""}://${host}${port ? `:${port}` : ''}/frontend/api/ws?token=${token}`,
      reconnect: true,
    });

    websocket.onOpen(() => {
      websocketConnected.value = true;
      if (route.name) {
        sendCurrentPage(route);
      }
    })

    websocket.onClose(() => websocketConnected.value = false)

    router.afterEach((to, from, failure) => {
      if (isNavigationFailure(failure)) {
        return;
      }

      if (!authStatusLoaded.value) {
        return;
      }

      const isOnInternalPage = !["login", "register"].includes(to.name as string)

      // user is not authenticated, reload page to get redirected to correct page
      if (isOnInternalPage && !userStore.isAuthenticated) {
        window.location.reload();
        return;
      }

      // redirect to home if the user is on the login page
      if (!isOnInternalPage && userStore.isAuthenticated) {
        router.push("/");
        return;
      }

      if (websocketConnected.value) {
        sendCurrentPage(to);
      }
    });

    // ------------------------

    const authenticate = () => userStore
      .authenticate()
      .then(() => authStatusLoaded.value = true)
      .catch(({ status }) => {
        // we use "StatusPreconditionRequired" to indicate that no user accounts exist (new installation)
        if (status === 428 && route.name !== "register") {
          router.replace({ name: "register" })
            .then(() => authStatusLoaded.value = true);
          return
        }

        // 401 means the user is not authenticated
        if (status === 401 && route.name !== "login") {
          router.replace({ name: "login" })
            .then(() => authStatusLoaded.value = true);
          return
        }

        // other errors (e.g. timeout), retry in 2 seconds
        setTimeout(authenticate, 2e3);
      });


    authenticate()

    watch(isAuthenticated, () => {
      if (isAuthenticated && !websocketConnected.value) {
        wsConnect(userStore.getAuthToken());
      }
    });

    onBeforeUnmount(() => websocket.close())

    return {
      websocketConnected,
      authStatusLoaded,
      isAuthenticated,
      error,
      clearError: global.clearError,
    };
  },
});
</script>

<style lang="scss">
@import "/src/styles/main";
</style>
