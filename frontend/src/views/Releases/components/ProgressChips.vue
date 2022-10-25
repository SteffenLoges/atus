<template>
  <div class="d-flex flex-column flex-md-row">
    <Chip v-for="chip in chips" :key="chip.text" class="mb-1 mb-md-0 mr-md-2" text-color="white" :color="chip.color"
      :icon="chip.icon" :size="size" :text="chip.text" :loading="!chip.completed"
      :badge="chip.cnt > 1 ? chip.cnt : undefined" :tooltip="chip.tooltip" v-bind="$attrs" />
  </div>
</template>


<script lang="ts">
import { defineComponent, PropType, toRefs, computed } from "vue";
import useMetaFiles from "../composables/metaFiles";
import { mdiAlert, mdiCheck } from "@mdi/js";

export default defineComponent({
  inheritAttrs: false,
  props: {
    metaFiles: {
      type: Array as PropType<IMetaFile[]>,
      required: true,
    },
    size: {
      type: String,
      default: "default",
    },
  },
  setup(props) {
    const { getName: getMetaFileName } = useMetaFiles();
    const { metaFiles } = toRefs(props);

    const chips = computed(() =>
      metaFiles.value
        .filter((status, index, self) => index === self.findIndex((s) => s.type === status.type))
        .map((s) => {
          const allSameType = metaFiles.value.filter((m) => m.type === s.type);
          const anyError = allSameType.some((m) => m.state === "ERROR");
          const completed = allSameType.find((status) => !['PROCESSED', 'ERROR'].includes(status.state)) === undefined

          return {
            text: getMetaFileName(s.type),
            icon: anyError ? mdiAlert : mdiCheck,
            color: s.state === 'ERROR' ? 'red-darken-2 text-white' : 'yellow',
            tooltip: anyError ? 'Error while processing file. Check log!' : !completed ? 'File is still being processed' : undefined,
            cnt: metaFiles.value.filter((status) => status.type === s.type)?.length || 0,
            completed,
          }
        })
        .sort((a, b) => b.text.localeCompare(a.text))
    );

    return {
      chips,
    };
  },
});
</script>
