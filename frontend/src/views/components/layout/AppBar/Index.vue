<template>
  <v-app-bar :elevation="0" :height="mdAndUp ? 35 : 50" class="main-app-bar pr-1">
    <v-btn class="mr-2 d-md-none" @click="toggleForceDrawer" :icon="mdiMenu" />

    <Breadcrumbs class="d-none d-md-block ml-4" />
    <v-spacer />
    <Donate class="mr-3" />
    <UserMenu class="mr-2" />
  </v-app-bar>
</template>


<script lang="ts">
import { defineComponent } from "vue";
import { useDisplay } from "vuetify";
import useGlobalStore from "@/store/global";
import Breadcrumbs from "./components/Breadcrumbs.vue";
import UserMenu from "./components/UserMenu.vue";
import Donate from "./components/Donate.vue";
import { mdiMenu } from "@mdi/js";

export default defineComponent({
  components: {
    Breadcrumbs,
    UserMenu,
    Donate,
  },
  setup() {
    const globalStore = useGlobalStore();
    const { mdAndUp } = useDisplay();

    return {
      mdAndUp,
      toggleForceDrawer: globalStore.toggleForceDrawer,
      mdiMenu
    };
  },
});
</script>


<style lang="scss">
@use "vuetify/styles/settings/variables" as v;

.main-app-bar {
  .v-toolbar__content {
    padding: 0 !important;
  }

  overflow: visible !important;

  // add rounded corner
  &:before {
    $border-width: 20px;
    pointer-events: none;
    position: absolute;
    bottom: -$border-width;
    left: -$border-width;
    width: 100%;
    height: $border-width * 2;
    content: "";
    background-color: transparent;
    border-color: rgb(var(--v-theme-surface));
    border-style: solid;
    border-width: $border-width 0 0 $border-width;
    border-top-left-radius: $border-width * 2;
    visibility: hidden;

    @media #{(map-get(v.$display-breakpoints, "md-and-up"))} {
      visibility: visible;
    }
  }
}

.divider {
  margin-top: -10px;
}
</style>