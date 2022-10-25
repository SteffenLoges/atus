<template>
  <ConfirmDialog :modelValue="deleteUUID != ''" @cancel="deleteUUID = ''" @confirm="onDeleteConfirm(deleteUUID)">
    Are you sure you want to delete this source?
  </ConfirmDialog>

  <Card :loading="isLoading" title="Sources">
    <v-card-text class="px-0">
      <v-alert v-if="!sources.length" type="info" class="mx-3">You didn't add any sources yet.</v-alert>

      <Source v-for="(source, i) in sources" :class="{ 'mt-4': i > 0 }" :key="source.uid" v-bind="source"
        @toggle="toggle(source.uid, $event)" @delete="deleteUUID = source.uid" />
    </v-card-text>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="primary" :to="{ name: 'settings_sources_add' }">Add New Soruce</v-btn>
    </v-card-actions>
  </Card>
</template>


<script lang="ts">
import { ref, defineComponent, onBeforeUnmount } from "vue";
import { send } from "@/utils/websocket";
import { success, error } from "@/plugins/toast";
import Source from "./components/Source.vue";


export default defineComponent({
  components: {
    Source,
  },
  async setup() {
    const isLoading = ref(false);
    const sources = ref<ISource[]>([]);

    let loadSourcesIntervalEH = -1
    onBeforeUnmount(() => clearInterval(loadSourcesIntervalEH));

    // --------------------------------------------------------------------------

    const loadSources = () =>
      send("SETTINGS__SOURCES_MANAGE__GET_ALL")
        .then(({ payload }: IResponse<ISource[]>) => sources.value = payload || []);

    await loadSources()
    loadSourcesIntervalEH = setInterval(loadSources, 5e3);

    // --------------------------------------------------------------------------

    const deleteUUID = ref("");
    const onDeleteConfirm = (uid: string) => {
      deleteUUID.value = "";
      isLoading.value = true;

      send("SETTINGS__SOURCES_MANAGE__DELETE", { uid })
        .then(() => {
          loadSources().then(() => isLoading.value = false);
          success("Source deleted");
        })
        .catch(({ payload }: IResponse<string>) => {
          isLoading.value = false;
          error("Source couldn't be deleted", payload);
        });
    };

    // --------------------------------------------------------------------------

    const toggle = (uid: string, start: boolean) => {
      isLoading.value = true;

      send("SETTINGS__SOURCES_MANAGE__TOGGLE", { uid, start })
        .then(() => {
          loadSources().then(() => isLoading.value = false);
          success(start ? "Source started" : "Stopping source...");
        })
        .catch(({ payload }: IResponse<string>) => {
          isLoading.value = false;
          error(`Source couldn't be ${start ? "started" : "stopped"}`, payload);
        });
    };

    // --------------------------------------------------------------------------

    return {
      sources,
      isLoading,
      deleteUUID,
      toggle,
      onDeleteConfirm,
    };
  },
});
</script>
