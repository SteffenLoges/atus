<template>
  <FormCard :loading="isLoading" title="Sample Settings" @submit="onSubmit">
    <v-card-text>
      <v-alert type="info" class="mb-4">
        Changes may not effect already queued releases.
      </v-alert>

      <Switch label="Enable sample generation" v-model="enabled" />

      <TextField v-model="sumScreenshots" type="number" :min="0" :max="10" :maxlength="2" required
        label="Number of screenshots to generate" hint="Default: 3. Set to 0 to disable." persistent-hint
        class="mb-2" />

      <TextField v-model="minSize" type="number" :min="1" :maxlength="20" required
        label="Minimum size of a sample file in MiB" hint="Used to filter out false positives. Default: 2"
        persistent-hint class="mb-2" />

      <TextField v-model="maxSize" type="number" :min="1" :maxlength="20" required
        label="Maximum size of a sample file in MiB"
        hint="Used to filter out false positives. Larger files require more time to process. Default: 200"
        persistent-hint />
    </v-card-text>

    <v-card-actions class="justify-end">
      <v-btn color="primary" type="submit" :disabled="isLoading">
        Save
      </v-btn>
    </v-card-actions>
  </FormCard>
</template>



<script lang="ts">
import { defineComponent, ref } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";

export default defineComponent({
  async setup() {
    const globalStore = useGlobalStore();

    const isLoading = ref(false);
    const enabled = ref(false);
    const sumScreenshots = ref(0);
    const minSize = ref(0);
    const maxSize = ref(0);

    // --------------------------------------------------------------------------

    const r: IResponse<ISampleSettings> = await send("SETTINGS__SAMPLES_MANAGE__GET_ALL")
    enabled.value = r.payload.enabled;
    sumScreenshots.value = r.payload.sumScreenshots;
    minSize.value = r.payload.minSize;
    maxSize.value = r.payload.maxSize;

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__SAMPLES_MANAGE__SAVE", {
        enabled: enabled.value,
        sumScreenshots: parseInt("" + sumScreenshots.value),
        minSize: parseInt("" + minSize.value),
        maxSize: parseInt("" + maxSize.value),
      })
        .then(() => success("Settings saved successfully"))
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => isLoading.value = false)
    };

    // --------------------------------------------------------------------------

    return {
      enabled,
      sumScreenshots,
      minSize,
      maxSize,
      onSubmit,
      isLoading,
    };
  },
});
</script>
