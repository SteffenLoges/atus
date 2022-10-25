<template>
  <ConfirmDialog :modelValue="deleteUID != ''" @cancel="deleteUID = ''"
    @confirm="onDeleteConfirm(deleteUID); deleteUID = '';">
    Are you sure you want to delete this fileserver?
  </ConfirmDialog>

  <Card :loading="isLoading" title="Fileservers">
    <v-card-text class="px-0">
      <v-alert v-if="!fileservers.length" type="info" class="mx-3">You didn't add any fileservers yet.</v-alert>

      <Fileserver v-for="(fs, i) in fileservers" :class="{ 'mt-4': i > 0 }" :key="fs.uid" :uid="fs.uid" :name="fs.name"
        :enabled="fs.enabled" :filesDownloaded="fs.filesDownloaded" :serverLoad="fs.serverLoad"
        :diskFreeSpace="fs.diskFreeSpace" :diskTotalSpace="fs.diskTotalSpace" @delete="deleteUID = fs.uid"
        @toggle="toggle(fs.uid, $event)" />
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="primary" :to="{ name: 'settings_fileservers_add' }">Add New Fileserver</v-btn>
    </v-card-actions>
  </Card>
</template>


<script lang="ts">
import { ref, defineComponent, onBeforeUnmount } from "vue";
import { send } from "@/utils/websocket";
import { success, error } from "@/plugins/toast";
import Fileserver from "./components/Fileserver.vue";

export default defineComponent({
  components: {
    Fileserver,
  },
  async setup() {
    const isLoading = ref(false);

    // --------------------------------------------------------------------------

    const fileservers = ref<IFileserver[]>([]);

    const loadFileservers = () =>
      send("SETTINGS__FILESERVERS_MANAGE__GET_ALL")
        .then(({ payload }: IResponse<IFileserver[]>) => fileservers.value = payload || []);

    let loadFileserversIntervalEH = setInterval(loadFileservers, 3e3);
    onBeforeUnmount(() => clearInterval(loadFileserversIntervalEH));

    await loadFileservers()

    // --------------------------------------------------------------------------

    const deleteUID = ref("");
    const onDeleteConfirm = (uid: string) => {
      isLoading.value = true;

      send("SETTINGS__FILESERVERS_MANAGE__DELETE", { uid })
        .then(async () => {
          await loadFileservers();
          success("Fileserver deleted successfully.")
        })
        .catch(({ payload }: IResponse<string>) => error("Fileserver couldn't be deleted", payload))
        .finally(() => isLoading.value = false)
    };

    // --------------------------------------------------------------------------

    const toggle = (uid: string, start: boolean) => {
      isLoading.value = true;

      send("SETTINGS__FILESERVERS_MANAGE__TOGGLE", { uid, start })
        .then(async () => {
          await loadFileservers();
          success(`Fileserver ${start ? "enabled" : "disabled"}`);
        })
        .catch(({ payload }: IResponse<string>) => error("Fileserver couldn't be enabled", payload))
        .finally(() => isLoading.value = false)
    };

    // --------------------------------------------------------------------------

    return {
      fileservers,
      isLoading,
      deleteUID,
      toggle,
      onDeleteConfirm,
    };
  },
});
</script>
