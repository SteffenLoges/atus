<template>
  <div class="state text-high-emphasis text-subtitle-2">
    Status:
    <span class="text-uppercase font-weight-medium" :class="stateComputed.class" v-text="stateComputed.text"></span>
    <span v-if="state.state === 'UPLOADED' && uploadDate" class="font-weight-light text-high-emphasis">
      on
      <DateTimeLive :date="uploadDate" :updateInterval="-1" />
    </span>

    <template v-if="downloadState && downloadState?.done < 100">
      <span class="ml-1 font-weight-bold">
        @
        {{ (downloadState?.done ?? 0).toFixed(2) }}%
      </span>

      <span class="ml-1 text-medium-emphasis text-caption" v-if="downloadState?.state === 'STARTED'">
        -
        <span class="ml-1" v-text="`${bytesHumanReadable(downloadState.downloadRate)}/s`"></span>
        <span class="ml-2">eta: {{ etaHumanized }}</span>
      </span>
    </template>
  </div>
</template>


<script lang="ts">
import { defineComponent, PropType, toRefs, computed } from "vue";
import moment from "moment";
import { bytesHumanReadable } from "@/utils/conversion";
import DateTimeLive from "./DateTimeLive.vue";

export default defineComponent({
  components: {
    DateTimeLive,
  },
  props: {
    state: {
      type: Object as PropType<IReleaseState>,
      required: true,
    },
    uploadDate: {
      type: String,
      required: true,
    },
    downloadState: {
      type: Object as PropType<IDownloadState>,
    },
  },
  setup(props) {
    const { state, downloadState } = toRefs(props);

    const stateMap = {
      DOWNLOADED: { text: "Finished Downloading", class: "text-blue-lighten-4" },
      AWAITING_FS_RESP: { text: "Waiting for fileserver", class: "text-yellow" },
      UPLOADED: { text: "Uploaded", class: "text-green" },
      NEW: { text: "New", class: "text-blue" },
      STARTED: { text: "Downloading", class: "text-blue-lighten-1" },
      PAUSED: { text: "Paused", class: "text-yellow" },
      ERROR: { text: "Error - check ruTorrent", class: "text-red" },
      GENERAL_ERROR: { text: "Error - check log", class: "text-red" },
      UPLOAD_ERROR: { text: "Upload Error - check log", class: "text-red" },
      STOPPED: { text: "Stopped - check rtorrent", class: "text-red" },
      HASHING: { text: "Hashing", class: "text-orange" },
      CHECKING: { text: "Checking", class: "text-orange" },
      DOWNLOAD_INIT: { text: "Initializing Download", class: "text-blue-lighten-1" },
      DOWNLOADING: { text: "Downloading", class: "text-blue-lighten-1" },
    };

    const stateComputed = computed(() => {
      if (["UPLOADED", "NEW", "DOWNLOAD_INIT", "UPLOAD_ERROR"].includes(state.value.state)) {
        return stateMap[state.value.state];
      }

      if (!downloadState.value) {
        return stateMap.AWAITING_FS_RESP;
      }

      if (downloadState.value.done === 100) {
        return stateMap.DOWNLOADED;
      }

      return stateMap[downloadState.value.state] || { text: downloadState.value.state }
    });

    const etaHumanized = computed(() => {
      if (!downloadState.value?.eta) {
        return "∞";
      }

      return moment
        .duration(downloadState.value.eta, "seconds")
        .format("d[d] h[h] mm[m] ss[s]");
    });

    return {
      bytesHumanReadable,
      etaHumanized,
      stateComputed,
    };
  },
});
</script>