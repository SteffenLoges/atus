<template>
  <v-navigation-drawer ref="drawer" :modelValue="forceDrawer || mdAndUp" :permanent="mdAndUp"
    :rail="!lgAndUp && !forceDrawer" :temporary="lgAndUp || forceDrawer" expand-on-hover floating class="main-drawer">
    <template #prepend>
      <v-btn block stacked flat to="/" class="rounded-0 py-6 px-0 justify-start" variant="plain" :ripple="false">
        <div class="logo-wrapper d-flex align-center">
          <img src="/src/assets/images/logo.png" class="logo" />
          <img src="/src/assets/images/logo-text.png" class="logo-text" />
        </div>
      </v-btn>

      <v-divider class="mb-2" />
    </template>

    <Menu />

    <template #append>
      <v-divider class="mb-1" />
      <FileserverStatistics />
    </template>
  </v-navigation-drawer>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from "vue";
import { useDisplay } from "vuetify";
import { storeToRefs } from "pinia";
import useGlobalStore from "@/store/global";
import FileserverStatistics from "./components/Fileservers/Index.vue";
import Menu from "./components/Menu.vue";

export default defineComponent({
  components: {
    FileserverStatistics,
    Menu,
  },
  setup() {
    const globalStore = useGlobalStore()
    const { forceDrawer } = storeToRefs(globalStore)
    const { lgAndUp, mdAndUp } = useDisplay()

    const drawer = ref(null);

    // close drawer when user clicks outside
    const onClick: EventListener = (e: any) => {
      if (e.target?.closest(".v-navigation-drawer__scrim")) {
        globalStore.setForceDrawer(false);
      }
    };

    onMounted(() => document.addEventListener("click", onClick))
    onBeforeUnmount(() => document.removeEventListener("click", onClick))

    return {
      drawer,
      mdAndUp,
      lgAndUp,
      forceDrawer,
    };
  },
});
</script>


<style lang="scss">
.main-drawer {
  .logo-wrapper {
    overflow: hidden;
    white-space: nowrap;
    height: 75px;

    .logo {
      margin-left: 8px;
      height: 37px;
    }

    .logo-text {
      width: 0;
      margin-left: 15px;
      margin-top: 10px;
      height: 44px;
      opacity: 0;
    }

    .logo,
    .logo-text {
      transition: width 0.2s ease;
      transition-property: width, height, opacity,
        margin-left;
      transform-origin: right;
    }
  }

  &.v-navigation-drawer--is-hovering,
  &:not(.v-navigation-drawer--rail) {
    .logo-wrapper {
      .logo {
        margin-left: 35px;
      }

      .logo-text {
        width: 120px;
        opacity: 1;
      }
    }
  }

  &.v-navigation-drawer--rail:not(.v-navigation-drawer--is-hovering) .v-navigation-drawer__content {
    overflow: hidden;
  }
}
</style>