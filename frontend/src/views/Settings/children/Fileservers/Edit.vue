<template>
  <FormCard :loading="isLoading" title="Edit Fileserver" @submit="onSubmit">
    <v-card-text>
      <v-alert type="warning" class="mb-4">
        <strong>Warning!</strong>
        <p>
          {{ appName }} will not check URLs for validity
          when editing a fileserver.<br />
          If you are unsure whether the settings are valid
          or not, please delete the fileserver and add a new
          one.
        </p>
      </v-alert>

      <v-card variant="text" v-for="(c, i) in components" :class="{ 'mt-3': i > 0 }" :key="c.title" :title="c.title"
        class="card-accent">
        <v-card-text>
          <component v-if="!isLoading" :is="c.component" v-bind="c.binds.value" v-on="c.handlers" />
        </v-card-text>
      </v-card>
    </v-card-text>

    <v-card-actions class="px-5 justify-end">
      <v-btn color="error" @click.prevent="$router.push({ name: 'settings_fileservers_manage' })">Cancel</v-btn>
      <v-btn color="primary" type="submit">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";
import { useRouter, useRoute } from "vue-router";
import URL from "./components/URL.vue";
import Settings from "./components/Settings.vue";

export default defineComponent({
  async setup() {
    const globalStore = useGlobalStore();
    const router = useRouter();
    const route = useRoute();
    const appName = import.meta.env.VITE_APP_NAME

    const isLoading = ref(false);

    // --------------------------------------------------------------------------

    const uid = route.params.uid as string;
    const name = ref("");
    const url = ref("");
    const listInterval = ref(0);
    const statisticsInterval = ref(0);
    const minFreeDiskSpace = ref(0);
    const diskFreeSpace = ref(0);
    const diskTotalSpace = ref(0);

    // --------------------------------------------------------------------------

    const r: IResponse<IFileserver> = await send("SETTINGS__FILESERVERS_EDIT__GET", { uid })
    name.value = r.payload.name;
    url.value = r.payload.url;
    listInterval.value = r.payload.listInterval;
    statisticsInterval.value = r.payload.statisticsInterval;
    minFreeDiskSpace.value = r.payload.minFreeDiskSpace;
    diskFreeSpace.value = r.payload.diskFreeSpace || 0;
    diskTotalSpace.value = r.payload.diskTotalSpace || 0;

    // --------------------------------------------------------------------------

    const components = [
      {
        component: URL,
        title: "URL",
        binds: computed(() => ({
          url: url.value,
        })),
        handlers: {
          "update:url": (v: string) => (url.value = v),
        },
      },
      {
        component: Settings,
        title: "Settings",
        binds: computed(() => ({
          uid,
          name: name.value,
          listInterval: listInterval.value,
          statisticsInterval: statisticsInterval.value,
          diskFreeSpace: diskFreeSpace.value,
          diskTotalSpace: diskTotalSpace.value,
          minFreeDiskSpace: minFreeDiskSpace.value,
        })),
        handlers: {
          "update:name": (v: string) => (name.value = v),
          "update:listInterval": (v: number) => (listInterval.value = v),
          "update:statisticsInterval": (v: number) => (statisticsInterval.value = v),
          "update:minFreeDiskSpace": (v: number) => (minFreeDiskSpace.value = v),
        },
      },
    ];

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__FILESERVERS_EDIT__SAVE", {
        uid,
        name: name.value,
        url: url.value,
        listInterval: listInterval.value,
        statisticsInterval: statisticsInterval.value,
        minFreeDiskSpace: minFreeDiskSpace.value,
      })
        .then(() => {
          success("Settings saved successfully");
          router.push({ name: "settings_fileservers_manage" });
        })
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => (isLoading.value = false));
    };

    // --------------------------------------------------------------------------

    return {
      appName,
      onSubmit,
      components,
      isLoading,
    };
  },
});
</script>