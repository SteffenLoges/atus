<template>
  <v-overlay :modelValue="true" persistent scrim="rgb(70, 3, 3)" no-click-animation
    class="align-center justify-center error-overlay">
    <Card color="transparent" flat>
      <v-card-title class="text-center text-h4">Oh No! Something went terribly wrong!</v-card-title>
      <v-card-text class="mt-6">
        <small>Error details:</small>
        <v-alert density="compact" max-height="300px" class="overflow-auto text-caption">
          {{ error }}
          <br />
          Route: {{ $route.fullPath }}
        </v-alert>

        <div class="text-caption text-center mt-3">
          Try refreshing the page and restarting the {{ appName }}. If the problem persists, please open an issue on
          GitHub.
        </div>

        <div class="d-flex justify-space-evenly mt-7">
          <v-btn size="small" width="180" color="white" @click.prevent="reloadPage()">Reload Page</v-btn>
          <v-btn size="small" width="180" color="white" target="_blank"
            :href="dereferURL('https://github.com/SteffenLoges/atus/issues')">GitHub issues</v-btn>
        </div>
      </v-card-text>
    </Card>
  </v-overlay>
</template>


<script lang="ts">
import { defineComponent } from "vue";
import { dereferURL } from "@/utils/url";

export default defineComponent({
  props: {
    error: {
      default: "An Unknown Error Occurred",
    },
  },
  setup() {
    const reloadPage = () => window.location.reload();
    const appName = import.meta.env.VITE_APP_NAME

    return {
      appName,
      dereferURL,
      reloadPage,
    };
  },
});
</script>
