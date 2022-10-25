<template>
  <div class="h-100">
    <v-menu>
      <template v-slot:activator="{ props: menu }">
        <v-btn class="rounded-0 text-none h-100" block v-bind="menu">
          {{ name }}
          <v-icon right class="ml-2" :icon="mdiMenuDown" />
        </v-btn>
      </template>
      <v-list>
        <template v-for="item in items" :key="item.title">
          <v-divider v-if="item.divider" />

          <v-list-item v-else v-bind="item.binds" v-on="item.listeners || {}" class="pr-8">
            <template v-slot:prepend>
              <v-icon class="mr-3" :icon="item.icon" />
            </template>
            <v-list-item-title v-text="item.title" />
          </v-list-item>
        </template>
      </v-list>
    </v-menu>
  </div>
</template>


<script lang="ts">
import { defineComponent } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import useUserStore from "@/store/user";
import websocket from "@/utils/websocket";
import { dereferURL } from "@/utils/url";
import { mdiLock, mdiLogout, mdiBug, mdiMenuDown } from "@mdi/js";

export default defineComponent({
  setup() {
    const userStore = useUserStore();
    const { deleteAuthToken } = userStore;
    const { uid, name } = storeToRefs(userStore);

    const router = useRouter();

    const logout = () => {
      deleteAuthToken();
      websocket.close();
      router.push({ name: "login" });
    };

    const items = [
      {
        icon: mdiLock,
        title: "Edit Account",
        binds: {
          to: {
            name: "settings_users_edit",
            params: { uid: uid.value },
          },
        },
      },
      {
        icon: mdiLogout,
        title: "Logout",
        listeners: {
          click: logout,
        },
      },
      { divider: true },
      {
        icon: mdiBug,
        title: "Report a bug",
        binds: {
          href: dereferURL("https://github.com/SteffenLoges/atus/issues"),
          target: "_blank",
        },
      },
    ];

    return {
      name,
      items,
      logout,
      mdiMenuDown,
    };
  },
});
</script>
