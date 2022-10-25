<template>
  <v-list density="compact" class="pa-0 mb-4">
    <v-list-subheader>Fileservers</v-list-subheader>

    <v-list-item v-if="!fileserverStatistics?.length">
      <template #prepend>
        <v-icon size="small" class="mr-3" :icon="mdiAlertCircle" />
      </template>

      <v-list-item-title class="text-caption">No active fileservers found</v-list-item-title>
      <v-list-item-action class="mt-1">
        <v-btn color="primary" size="x-small" variant="tonal" :to="{
          name: 'settings_fileservers_manage',
        }">Manage fileservers</v-btn>
      </v-list-item-action>
    </v-list-item>

    <Fileserver v-for="fs of fileserverStatistics" :key="fs.uid" :enabled="fs.enabled" :name="fs.name"
      :statistics="fs.statistics || null" rounded="shaped" />
  </v-list>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { storeToRefs } from "pinia";
import useGlobalStore from "@/store/global";
import Fileserver from "./components/Fileserver.vue";
import { mdiAlertCircle } from "@mdi/js";

export default defineComponent({
  components: {
    Fileserver,
  },
  setup() {
    const gloablStore = useGlobalStore();

    const { fileserverStatistics } = storeToRefs(gloablStore);

    return {
      fileserverStatistics,
      mdiAlertCircle,
    };
  },
});
</script>