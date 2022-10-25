<template>
  <v-overlay :model-value="true" persistent no-click-animation scrim="#1c212e" class="align-center justify-center">
    <v-card color="transparent" flat>
      <v-card-text class="text-center">
        <div class="py-4 px-1">
          <v-progress-linear indeterminate color="white" size="100" />
        </div>
        <div class="text-overline">Connecting to server, please wait...</div>
        <div class="d-flex justify-space-evenly mt-5">
          <v-btn size="small" width="150" color="white" @click.prevent="reloadPage()">Reload Page</v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-overlay>
</template>


<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  setup() {
    return {
      reloadPage: () => window.location.reload(),
    };
  },
});
</script>