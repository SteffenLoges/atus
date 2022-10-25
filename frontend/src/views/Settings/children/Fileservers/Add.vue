<template>
  <ErrorDialog v-model="errorState.error">
    {{ errorState.text }}
    <template #details v-if="errorState.details != ''">{{ errorState.details }}</template>
  </ErrorDialog>

  <FormCard :loading="loadingState.loading" :loadingText="loadingState.text" title="Add New Fileserver"
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
      <v-btn color="error" @click.prevent="$router.push({ name: 'settings_fileservers_manage' })">Cancel</v-btn>
      <v-btn color="primary" type="submit">{{ selectedStepIndex == steps.length - 1 ? "Finish" : "Next" }}</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import { useRouter } from "vue-router";
import URL from "./components/URL.vue";
import Settings from "./components/Settings.vue";

export default defineComponent({
  setup() {
    const router = useRouter();

    // --------------------------------------------------------------------------

    const url = ref("");
    const uid = ref("");
    const name = ref("");
    const listInterval = ref(0);
    const statisticsInterval = ref(0);
    const diskFreeSpace = ref(0);
    const diskTotalSpace = ref(0);
    const minFreeDiskSpace = ref(0);

    const loadingState = ref({ loading: false, text: "" });
    const errorState = ref({ error: false, text: "", details: "" });

    // --------------------------------------------------------------------------

    const steps = [
      {
        component: URL,
        title: "URL",
        loadingText: "Checking URL, this may take a while...",
        errorText: "Error checking URL",
        onSubmit: () =>
          send("SETTINGS__FILESERVERS_ADD__SET_URL", { url: url.value })
            .then(({ payload }: IResponse<IFileserver>) => {
              uid.value = payload.uid;
              name.value = payload.name;
              listInterval.value = payload.listInterval;
              statisticsInterval.value = payload.statisticsInterval;
              diskFreeSpace.value = payload.diskFreeSpace || 0
              diskTotalSpace.value = payload.diskTotalSpace || 0
              minFreeDiskSpace.value = payload.minFreeDiskSpace;
            }),
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
        loadingText: "Saving Settings...",
        errorText: "Error while saving Settings",
        onSubmit: () =>
          send("SETTINGS__FILESERVERS_ADD__SET_SETTINGS", {
            uid: uid.value,
            name: name.value,
            listInterval: parseInt("" + listInterval.value),
            statisticsInterval: parseInt("" + statisticsInterval.value),
            minFreeDiskSpace: parseInt("" + minFreeDiskSpace.value),
          }).then(() => {
            success("Settings saved", "Fileserver added successfully");
            router.push({ name: "settings_fileservers_manage", });
          }),
        binds: computed(() => ({
          uid: uid.value,
          name: name.value,
          listInterval: listInterval.value,
          statisticsInterval: statisticsInterval.value,
          diskFreeSpace: diskFreeSpace.value,
          diskTotalSpace: diskTotalSpace.value,
          minFreeDiskSpace: minFreeDiskSpace.value,
        })),
        handlers: {
          "update:name": (v: string) => name.value = v,
          "update:listInterval": (v: number) => listInterval.value = v,
          "update:statisticsInterval": (v: number) => statisticsInterval.value = v,
          "update:minFreeDiskSpace": (v: number) => minFreeDiskSpace.value = v,
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
              return;
            }
          })
          .catch(({ payload }: IResponse<any>) => {
            errorState.value = {
              error: true,
              text: selectedStep.value.errorText || "",
              details: payload,
            };
          })
          .finally(() => {
            loadingState.value = { loading: false, text: "" };
          });
      }
    };

    // --------------------------------------------------------------------------

    return {
      steps,
      selectedStepIndex,
      selectedStep,
      loadingState,
      errorState,
      name,
      onSubmit,
    };
  },
});
</script>

