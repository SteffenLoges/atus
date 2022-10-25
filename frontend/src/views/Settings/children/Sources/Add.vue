<template>
  <div style="position: relative">
    <ErrorDialog v-model="errorState.error" max-width="85vw">
      {{ errorState.text }}
      <template #details v-if="errorState.details != ''">{{ errorState.details }}</template>
    </ErrorDialog>

    <FormCard :loading="loadingState.loading" :loadingText="loadingState.text" title="Add New Source"
      @submit="onSubmit">
      <v-card-text>
        <Stepper :steps="steps.map((s) => s.title)" :activeStep="selectedStepIndex" />

        <SlideTransition>
          <div class="mt-8 mb-4 mx-3" :key="selectedStep.title">
            <component :is="selectedStep.component" isSetup v-bind="selectedStep.binds.value"
              v-on="selectedStep.handlers" />
          </div>
        </SlideTransition>
      </v-card-text>

      <v-card-actions class="px-5">
        <v-btn v-if="selectedStepIndex > 0" @click="selectedStepIndex--">Back</v-btn>

        <v-spacer />
        <v-btn color="error" @click.prevent="$router.push({ name: 'settings_sources_manage' })">Cancel</v-btn>
        <v-btn color="primary" type="submit">{{ selectedStepIndex == steps.length - 1 ? "Finish" : "Next" }}</v-btn>
      </v-card-actions>
    </FormCard>
  </div>
</template>


<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import RSSURL from "./components/RSSURL.vue";
import Settings from "./components/Settings.vue";
import ImageURL from "./components/ImageURL.vue";
import MetaFile from "./components/MetaFile.vue";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import { useRouter } from "vue-router";


export default defineComponent({
  setup() {
    const router = useRouter();

    const loadingState = ref({ loading: false, text: "" });
    const errorState = ref({ error: false, text: "", details: "" });

    // --------------------------------------------------------------------------
    // setting default values for the new source
    const uid = ref("");
    const name = ref("");
    const favicon = ref("");
    const rssURL = ref("");
    const requiresCookies = ref(false);
    const cookies = ref<ICookie[]>([
      { name: "uid", value: "" },
      { name: "pass", value: "" },
    ]);
    const rssInterval = ref(60 * 2);
    const requestWaitTime = ref(60 * 5);
    const metaPath = ref("");
    const metaPathUseAsKey = ref(false);
    const metaPathAutoDetected = ref(false);
    const imagePath = ref("");
    const imagePathUseAsKey = ref(false);
    const imagePathAutoDetected = ref(false);

    // --------------------------------------------------------------------------

    const steps = [
      {
        component: RSSURL,
        title: "RSS URL",
        loadingText: "Checking RSS Feed, this may take a while...",
        errorText: "Error checking RSS Feed",
        onSubmit: () =>
          send("SETTINGS__SOURCES_ADD__SET_RSS_URL", {
            url: rssURL.value,
            cookies: requiresCookies.value ? cookies.value : [],
          }).then(({ payload }: IResponse<ISource>) => {
            uid.value = payload.uid;
            name.value = payload.name;
            favicon.value = payload.favicon;
            metaPath.value = payload.metaPath;
            metaPathUseAsKey.value = payload.metaPathUseAsKey;
            metaPathAutoDetected.value = payload.metaPathAutoDetected;
            imagePath.value = payload.imagePath;
            imagePathUseAsKey.value = payload.imagePathUseAsKey;
            imagePathAutoDetected.value = payload.imagePathAutoDetected;
          }),
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
        loadingText: "Checking Meta File...",
        errorText: "Error checking Meta File",
        onSubmit: () =>
          send("SETTINGS__SOURCES_ADD__SET_META_PATH", {
            uid: uid.value,
            metaPath: metaPath.value,
            metaPathUseAsKey: metaPathUseAsKey.value,
          }),
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
        loadingText: "Checking Image File",
        errorText: "Error checking Image File",
        onSubmit: () =>
          send("SETTINGS__SOURCES_ADD__SET_IMAGE_PATH", {
            uid: uid.value,
            imagePath: imagePath.value,
            imagePathUseAsKey: imagePathUseAsKey.value,
          }),
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
      {
        component: Settings,
        title: "Settings",
        loadingText: "Saving Settings...",
        errorText: "Error while saving Settings",
        onSubmit: () =>
          send("SETTINGS__SOURCES_ADD__SET_SETTINGS", {
            uid: uid.value,
            name: name.value,
            rssInterval: parseInt("" + rssInterval.value),
            requestWaitTime: parseInt("" + requestWaitTime.value),
          }).then(() => {
            success("Source added successfully");
            router.push({ name: "settings_sources_manage" });
          }),

        binds: computed(() => ({
          uid: uid.value,
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
    ];

    let selectedStepIndex = ref(0);
    let selectedStep = computed(() => steps[selectedStepIndex.value]);

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      loadingState.value = {
        loading: true,
        text: selectedStep.value.loadingText || "",
      };

      if (selectedStep.value.onSubmit) {
        selectedStep.value
          .onSubmit()
          .then(() => {
            if (selectedStepIndex.value < steps.length - 1) {
              selectedStepIndex.value++;
            }
          })
          .catch(({ payload }: IResponse<any>) =>
            errorState.value = {
              error: true,
              text: selectedStep.value.errorText || "",
              details: payload,
            }
          )
          .finally(() => loadingState.value = { loading: false, text: "" });
      }
    };

    // --------------------------------------------------------------------------

    return {
      onSubmit,
      steps,
      selectedStepIndex,
      selectedStep,
      loadingState,
      errorState,
    };
  },
});
</script>

