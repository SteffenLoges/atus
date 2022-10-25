<template>
  <v-container fluid class="d-flex align-center justify-center mt-4">
    <Card class="px-4 py-3" title="Let's get started!">
      <v-card-text>
        <p class="pr-6 text-subtitle-1">
          Welcome to <strong>{{appName}}</strong>.<br />
          You need to configure a couple of things before you can start using the app.
        </p>

        <ul style="list-style-type: none;" class="mt-5">
          <li v-for="t of todos" :key="t.title">
            <v-btn :to="t.to" variant="text" class="text-none justify-start mx-n4" block>
              <v-icon :color="t.done ? 'success' : 'error'" :icon="t.done ? mdiCheck : mdiClose" />
              <span class="ml-2">{{t.title}}</span>
            </v-btn>
          </li>
        </ul>
      </v-card-text>
    </Card>
  </v-container>
</template>


<script lang="ts">
import { defineComponent } from "vue";
import { storeToRefs } from "pinia";
import useGlobalStore from "@/store/global";
import { mdiCheck, mdiClose } from "@mdi/js";

export default defineComponent({
  setup() {
    const globalStore = useGlobalStore();
    const { setupStatus } = storeToRefs(globalStore);

    const todos = [
      {
        title: 'Configure the upload settings',
        done: setupStatus.value.UPLOAD_CONFIGURED,
        to: { name: 'settings_destinations' },
      },
      {
        title: 'Add a fileservers',
        done: setupStatus.value.FILESERVER_ADDED,
        to: { name: 'settings_fileservers_add' },
      },
      {
        title: 'Add a source',
        done: setupStatus.value.SOURCE_ADDED,
        to: { name: 'settings_sources_add' },
      }
    ]

    return {
      appName: import.meta.env.VITE_APP_NAME,
      todos,
      mdiCheck,
      mdiClose,
    }
  },
});
</script>