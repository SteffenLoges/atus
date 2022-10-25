<template>
  <FormCard :loading="isLoading" title="Edit Source" @submit="onSubmit">
    <v-card-text>
      <v-alert type="warning" class="mb-4">
        <strong>Warning!</strong>
        <p>
          {{ appName }} will not check source URLs for
          validity when editing a source.<br />
          If you are unsure whether the settings are valid
          or not, please delete the source and add a new
          one.
        </p>
      </v-alert>

      <v-card v-for="(c, i) in components" :key="c.title" variant="text" class="card-accent" :class="{ 'mt-3': i > 0 }"
        :title="c.title">
        <v-card-text>
          <component v-if="!isLoading" :is="c.component" v-bind="c.binds.value" v-on="c.handlers" />
        </v-card-text>
      </v-card>
    </v-card-text>

    <v-card-actions class="px-5 justify-end">
      <v-btn color="error" @click.prevent="$router.push({ name: 'settings_sources_manage' })">Cancel</v-btn>
      <v-btn color="primary" type="submit">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import RSSURL from "./components/RSSURL.vue";
import Settings from "./components/Settings.vue";
import ImageURL from "./components/ImageURL.vue";
import MetaFile from "./components/MetaFile.vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";
import { useRouter, useRoute } from "vue-router";


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
    const favicon = ref("");
    const rssURL = ref("");
    const requiresCookies = ref(false);
    const cookies = ref<ICookie[]>();
    const rssInterval = ref(0);
    const requestWaitTime = ref(0);
    const metaPath = ref("");
    const metaPathUseAsKey = ref(false);
    const metaPathAutoDetected = ref(false);
    const imagePath = ref("");
    const imagePathUseAsKey = ref(false);
    const imagePathAutoDetected = ref(false);

    // --------------------------------------------------------------------------

    const r: IResponse<ISource> = await send("SETTINGS__SOURCES_EDIT__GET", { uid })
    name.value = r.payload.name;
    favicon.value = r.payload.favicon;
    rssURL.value = r.payload.rssURL;
    cookies.value = r.payload.cookies || [];
    requiresCookies.value = cookies.value.length > 0;
    rssInterval.value = r.payload.rssInterval;
    requestWaitTime.value = r.payload.requestWaitTime;
    metaPath.value = r.payload.metaPath;
    metaPathUseAsKey.value = r.payload.metaPathUseAsKey;
    metaPathAutoDetected.value = false;
    imagePath.value = r.payload.imagePath;
    imagePathUseAsKey.value = r.payload.imagePathUseAsKey;
    imagePathAutoDetected.value = false;

    // --------------------------------------------------------------------------

    const components = [
      {
        component: Settings,
        title: "Settings",
        binds: computed(() => ({
          uid,
          name: name.value,
          favicon: favicon.value,
          rssInterval: rssInterval.value,
          requestWaitTime: requestWaitTime.value,
        })),
        handlers: {
          "update:name": (v: string) => name.value = v,
          "update:rssInterval": (v: number) => rssInterval.value = v,
          "update:requestWaitTime": (v: number) => requestWaitTime.value = v,
        },
      },
      {
        component: RSSURL,
        title: "RSS URL",
        binds: computed(() => ({
          rssURL: rssURL.value,
          requiresCookies: requiresCookies.value,
          cookies: cookies.value,
        })),
        handlers: {
          "update:rssURL": (v: string) => rssURL.value = v,
          "update:requiresCookies": (v: boolean) => requiresCookies.value = v,
          "update:cookies": (v: ICookie[]) => cookies.value = v,
        },
      },
      {
        component: MetaFile,
        title: "Meta File",
        binds: computed(() => ({
          metaPath: metaPath.value,
          metaPathUseAsKey: metaPathUseAsKey.value,
          metaPathAutoDetected: metaPathAutoDetected.value,
        })),
        handlers: {
          "update:metaPath": (v: string) => metaPath.value = v,
          "update:metaPathUseAsKey": (v: boolean) => metaPathUseAsKey.value = v,
        },
      },
      {
        component: ImageURL,
        title: "Image",
        binds: computed(() => ({
          imagePath: imagePath.value,
          imagePathUseAsKey: imagePathUseAsKey.value,
          imagePathAutoDetected: imagePathAutoDetected.value,
        })),
        handlers: {
          "update:imagePath": (v: string) => imagePath.value = v,
          "update:imagePathUseAsKey": (v: boolean) => imagePathUseAsKey.value = v,
        },
      },
    ];

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__SOURCES_EDIT__SAVE", {
        uid,
        name: name.value,
        rssURL: rssURL.value,
        cookies: requiresCookies.value ? cookies.value : [],
        rssInterval: parseInt("" + rssInterval.value),
        requestWaitTime: parseInt("" + requestWaitTime.value),
        metaPath: metaPath.value,
        metaPathUseAsKey: metaPathUseAsKey.value,
        imagePath: imagePath.value,
        imagePathUseAsKey: imagePathUseAsKey.value,
      })
        .then(() => {
          success("Settings saved successfully");
          router.push({ name: "settings_sources_manage" });
        })
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => isLoading.value = false);
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