<template>
  <FormCard :loading="isLoading" title="Sample Settings" @submit="onSubmit">
    <v-card-text>
      <v-card variant="text" title="Main Settings" class="card-accent mb-4">
        <v-card-text>
          <v-select v-model="allocationMethod as any" :items="allocationMethods" label="Allocation Method" />
        </v-card-text>
      </v-card>

      <v-card variant="text" title="Main Settings" class="card-accent mb-4">
        <v-card-text>
          <v-alert type="warning" class="mb-4">
            Changing labels will break already queued
            releases.
            <br />
            Make sure that there are no pending releases
            that haven't been uploaded before changing
            labels.
            <br>
            <p class="mt-2">
              <i>Should not be left blank.<br> Labels are used to filter unrelated files from the fileserver wich
                significantly reduces the amount of data that needs to be transferred and processed.</i>
            </p>
          </v-alert>

          <TextField v-model="downloadLabel" label="Download Label" hint="Will be visible in ruTorrent" persistent-hint
            class="mb-2" />

          <TextField v-model="uploadLabel" label="Upload Label" hint="Will be visible in ruTorrent" persistent-hint />
        </v-card-text>
      </v-card>
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
import useGlobalStore from "@/store/global";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import { bytesHumanReadable } from "@/utils/conversion";

export default defineComponent({
  async setup() {
    const globalStore = useGlobalStore();

    const isLoading = ref(false);

    // --------------------------------------------------------------------------

    const allocationMethod = ref("");
    const allocationMethods = [
      { title: "Use the same fileserver until it is full", value: "FILL" },
      { title: "Use fileserver with most free space", value: "MOST_FREE" },
      { title: "Pick a random fileserver", value: "RANDOM" },
    ];
    const downloadLabel = ref("");
    const uploadLabel = ref("");

    // --------------------------------------------------------------------------

    const r: IResponse<IFileserverSettings> = await send("SETTINGS__FILESERVERS_SETTINGS__GET")
    allocationMethod.value = r.payload.allocationMethod;
    downloadLabel.value = r.payload.downloadLabel;
    uploadLabel.value = r.payload.uploadLabel;

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__FILESERVERS_SETTINGS__SAVE", {
        allocationMethod: allocationMethod.value,
        downloadLabel: downloadLabel.value,
        uploadLabel: uploadLabel.value,
      })
        .then(() => success("Settings saved successfully"))
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => (isLoading.value = false));
    };

    // --------------------------------------------------------------------------

    return {
      allocationMethod,
      allocationMethods,
      downloadLabel,
      uploadLabel,
      onSubmit,
      isLoading,
      bytesHumanReadable,
    };
  },
});
</script>
