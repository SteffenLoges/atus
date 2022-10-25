<template>
  <TextField :modelValue="name" @update:modelValue="$emit('update:name', $event)" persistent-hint label="Name"
    persistent-placeholder placeholder="e.g. Awesome Tracker" :counter="100" :maxlength="100" class="mb-2"
    hint="Enter a name for this source." required>
    <template #prepend-inner v-if="favicon">
      <v-img :src="getFileURL(`favicons/${favicon}`)" :width="30" :height="30" class="mr-3" />
    </template>
  </TextField>

  <TextField type="number" :modelValue="rssInterval" @update:modelValue="$emit('update:rssInterval', $event)" :min="30"
    :max="108000" :maxlength="6" required label="Seconds between checks for new torrents"
    hint="Default: 120 seconds (2 minutes)." persistent-hint class="mb-2">
    <template #append-inner>
      <span class="text-no-wrap text-high-emphasis">{{ rssIntervalHumanized }}
      </span>
    </template>
  </TextField>

  <TextField type="number" :modelValue="requestWaitTime" @update:modelValue="$emit('update:requestWaitTime', $event)"
    :min="30" :max="10000" :maxlength="6" required
    hint="Wait time when downloading torrent / image files. Default: 300ms. Low values may result in IP bans."
    class="mb-2" persistent-hint label="Wait time between requests in milliseconds">
    <template #append-inner>
      <span class="text-no-wrap text-high-emphasis">{{ requestWaitTimeHumanized }}</span>
    </template>
  </TextField>

  <small class="font-italic bg-grey-darken-3 px-2 py-1 text-medium-emphasis">
    Internal ID: {{ uid }}
  </small>
</template>


<script lang="ts">
import { defineComponent, toRefs, computed } from "vue";
import { getFileURL } from "@/utils/url";
import moment from "moment";

export default defineComponent({
  props: {
    isSetup: {
      type: Boolean,
      default: false,
    },
    uid: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    favicon: {
      type: String,
      required: true,
    },
    rssInterval: {
      type: Number,
      required: true,
    },
    requestWaitTime: {
      type: Number,
      required: true,
    },
  },
  emits: [
    "update:name",
    "update:rssInterval",
    "update:requestWaitTime",
  ],
  setup(props) {
    const { rssInterval, requestWaitTime } = toRefs(props);

    const rssIntervalHumanized = computed(() => {
      try {
        return moment
          .duration(rssInterval.value, "seconds")
          .format("d[d]hh[h]mm[m]ss[s]");
      } catch (e) {
        return "";
      }
    });

    const requestWaitTimeHumanized = computed(() => {
      try {
        return moment
          .duration(requestWaitTime.value, "milliseconds")
          .format("ss[s]SSS[ms]");
      } catch (e) {
        return "";
      }
    });

    return {
      appName: import.meta.env.VITE_APP_NAME,
      rssIntervalHumanized,
      requestWaitTimeHumanized,
      getFileURL,
    };
  },
});
</script>
