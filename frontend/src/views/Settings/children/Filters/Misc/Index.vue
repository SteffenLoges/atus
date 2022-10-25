<template>
  <FormCard :loading="isLoading" title="Miscellaneous Filters" @submit="onSubmit">
    <v-card-text>
      <v-alert type="info" class="mb-4">
        {{ appName }} relies on <a :href="dereferURL('https://predb.ovh/')" target="_blank">predb.ovh</a> to fetch
        informations (e.g. pre-time) about releases.<br>
        <p class="mt-2">We are not affiliated with predb.ovh in any way, therefore we cannot guarantee the accuracy of
          the information provided.</p>
      </v-alert>

      <v-card variant="text" class="card-accent">
        <v-card-text>
          <TextField v-model="maxAge" type="number" :min="0" :max="999999999" :maxlength="9" required
            label="Maximum age in minutes since pre" hint="Use 0 to disable this filter" persistent-hint>
            <template #append-inner>
              <span class="text-no-wrap text-high-emphasis">{{ maxAgeHumanized }}</span>
            </template>
          </TextField>
        </v-card-text>
      </v-card>
    </v-card-text>

    <v-card-actions class="px-5 justify-end">
      <v-btn color="primary" type="submit">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import moment from "moment";
import useGlobalStore from "@/store/global";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import { dereferURL } from "@/utils/url";


export default defineComponent({
  async setup() {
    const globalStore = useGlobalStore();
    const appName = import.meta.env.VITE_APP_NAME

    let isLoading = ref(false);
    let maxAge = ref(0);

    const maxAgeHumanized = computed(() => {
      return moment
        .duration(maxAge.value, "minutes")
        .format("d[d] h[h] mm[m] ss[s]");
    });

    // --------------------------------------------------------------------------

    const resp: IResponse<IFiltersMisc> = await send("SETTINGS__FILTERS_MISC__GET_ALL")
    maxAge.value = resp.payload.maxAge

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__FILTERS_MISC__SAVE", { maxAge: maxAge.value })
        .then(() => success("Settings saved successfully"))
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => isLoading.value = false);
    };

    // --------------------------------------------------------------------------

    return {
      appName,
      maxAge,
      maxAgeHumanized,
      onSubmit,
      isLoading,
      dereferURL,
    };
  },
});
</script>
