<template>
  <v-list-item rounded="shaped">
    <template v-slot:prepend>
      <v-icon :color="enabled ? 'success' : 'error'" size="small" class="mr-3" :icon="mdiCircle" />
    </template>

    <v-list-item-title v-text="name" />
    <v-list-item-subtitle>
      <template v-if="!enabled">
        <div class="text-caption">Fileserver disabled</div>
      </template>
      <template v-else-if="!statistics?.serverLoad">
        <div class="text-caption">Status unknown</div>
      </template>
      <v-progress-linear v-else :modelValue="percentUsage" :color="percentUsage > 80 ? 'red' : 'blue-grey'" height="18"
        class="text-caption" rounded>
        {{ bytesHumanReadable(statistics.diskFreeSpace) }}
        /
        {{ bytesHumanReadable(statistics.diskTotalSpace) }}
      </v-progress-linear>
    </v-list-item-subtitle>
  </v-list-item>
</template>


<script lang="ts">
import { defineComponent, computed, PropType, toRefs } from "vue";
import { bytesHumanReadable } from "@/utils/conversion";
import { mdiCircle } from "@mdi/js";

export default defineComponent({
  props: {
    enabled: {
      type: Boolean,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    statistics: {
      type: Object as PropType<IFileserverStatistics | null>,
      required: true,
    },
  },
  setup(props) {
    const { statistics } = toRefs(props);

    const percentUsage = computed(() => {
      if (!statistics.value?.diskTotalSpace) {
        return 0;
      }

      return (((statistics.value.diskTotalSpace - statistics.value.diskFreeSpace) / statistics.value.diskTotalSpace) * 100);
    });

    return {
      percentUsage,
      bytesHumanReadable,
      mdiCircle
    };
  },
});
</script>