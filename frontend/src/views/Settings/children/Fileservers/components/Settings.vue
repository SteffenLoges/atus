<template>
  <v-alert type="success" class="mb-6" v-if="isSetup">
    Good News!
    <p class="my-2">The server is reachable and ready to use.</p>
  </v-alert>

  <TextField :modelValue="name" @update:modelValue="$emit('update:name', $event)" persistent-hint label="Name"
    persistent-placeholder placeholder="e.g. The big boy" :counter="100" :maxlength="100" class="mb-2"
    hint="Enter a name for this fileserver." required />

  <TextField :modelValue="listInterval" @update:modelValue="$emit('update:listInterval', parseInt($event))"
    type="number" :min="1" :max="300" :maxlength="6" required label="Filelist update interval in seconds" class="mb-2"
    hint="Lower values will result in faster uploads, but may cause higher CPU usage on the fileserver. Default: 5"
    persistent-hint>
    <template #append-inner>
      <span class="text-no-wrap text-high-emphasis">{{ listIntervalHumanized }}</span>
    </template>
  </TextField>

  <TextField :modelValue="statisticsInterval" @update:modelValue="$emit('update:statisticsInterval', parseInt($event))"
    class="mb-2" type="number" :min="1" :max="300" :maxlength="6" required label="Statistics update interval in seconds"
    hint="Default: 10" persistent-hint>
    <template #append-inner>
      <span class="text-no-wrap text-high-emphasis">{{ statisticsIntervalHumanized }}</span>
    </template>
  </TextField>

  <TextField :modelValue="minFreeDiskSpace" @update:modelValue="$emit('update:minFreeDiskSpace', parseInt($event))"
    type="number" :min="1" :max="5000" :maxlength="6" required label="Minimum free disk space in GiB" hint="Default: 25"
    class="mb-2" persistent-hint>
    <template #append-inner v-if="diskTotalSpace > 0">
      <div class="text-no-wrap mt-n1">
        <v-chip color="blue-grey-lighten-2">
          {{ bytesHumanReadable(diskFreeSpace) }} free of
          {{ bytesHumanReadable(diskTotalSpace) }}
        </v-chip>
      </div>
    </template>
  </TextField>

  <small class="font-italic bg-grey-darken-3 px-2 py-1 text-medium-emphasis">
    Internal ID: {{ uid }}
  </small>
</template>


<script lang="ts">
import { defineComponent, toRefs, computed } from "vue";
import moment from "moment";
import { bytesHumanReadable } from "@/utils/conversion";

export default defineComponent({
  props: {
    isSetup: {
      type: Boolean,
      default: false,
    },
    uid: {
      type: String,
      default: "",
    },
    name: {
      type: String,
      default: "",
    },
    listInterval: {
      type: Number,
      default: 0,
    },
    statisticsInterval: {
      type: Number,
      default: 0,
    },
    diskFreeSpace: {
      type: Number,
      default: 0,
    },
    diskTotalSpace: {
      type: Number,
      default: 0,
    },
    minFreeDiskSpace: {
      type: Number,
      default: 0,
    },
  },
  emits: [
    "update:name",
    "update:listInterval",
    "update:statisticsInterval",
    "update:minFreeDiskSpace",
  ],
  setup(props) {
    const { listInterval, statisticsInterval } = toRefs(props);

    const listIntervalHumanized = computed(() => {
      try {
        return moment
          .duration(listInterval.value, "seconds")
          .format("hh[h]mm[m]ss[s]");
      } catch (e) {
        return "";
      }
    });

    const statisticsIntervalHumanized = computed(() => {
      try {
        return moment
          .duration(statisticsInterval.value, "seconds")
          .format("hh[h]mm[m]ss[s]");
      } catch (e) {
        return "";
      }
    });

    return {
      listIntervalHumanized,
      statisticsIntervalHumanized,
      bytesHumanReadable,
    };
  },
});
</script>
