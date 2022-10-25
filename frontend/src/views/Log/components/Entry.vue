<template>
  <Card :color="severities[entry.severity]?.color || ''" class="mt-2">
    <v-card-text class="py-2">
      <div class="text-body-2" v-text="entry.message"></div>
      <div class="text-caption text-high-emphasis d-flex flex-column flex-lg-row">
        <div class="font-weight-bold" v-text="severities[entry.severity]?.text || '---'"></div>
        <div class="ml-lg-2 text-capitalize" v-text="entry.type.toLowerCase().replaceAll('_', ' ')"></div>
        <div class="ml-lg-2">{{ entry.added }}</div>
        <div v-if="showRelease && entry.releaseUID && entry.releaseName">
          <span class="font-weight-bold mx-1 d-none d-lg-inline-block">|</span>
          <router-link class="text-black font-weight-medium text-decoration-none text-high-emphasis"
            :to="{ name: 'releases_details', params: { uid: entry.releaseUID, name: entry.releaseName } }"
            v-text="entry.releaseName" />
        </div>
      </div>
    </v-card-text>
  </Card>
</template>


<script lang="ts">
import { defineComponent, PropType } from "vue";
import { severities } from "../composables/log";

export default defineComponent({
  props: {
    entry: {
      type: Object as PropType<ILogEntry>,
      required: true,
    },
    showRelease: {
      type: Boolean,
      default: true,
    },
  },
  setup() {
    return { severities };
  },
});
</script>
