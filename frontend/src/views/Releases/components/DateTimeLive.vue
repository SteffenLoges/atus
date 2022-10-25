<template>
  {{ humanReadable }}
  <span class="mx-1 d-none d-sm-inline">({{ elapsed }})</span>
</template> 


<script lang="ts">
import { defineComponent, computed, onBeforeUnmount, ref, toRefs, watch } from "vue";
import moment from "moment";

export default defineComponent({
  props: {
    date: {
      type: String,
      required: true,
    },
    format: {
      type: String,
      default: "MM/DD/YYYY HH:mm:ss",
    },
    updateInterval: {
      type: Number,
      default: 5e3,
    },
  },
  setup(props) {
    const { date, format, updateInterval } = toRefs(props);

    const humanReadable = computed(() => moment(date.value).format(format.value));
    const elapsed = ref("");

    const calcElapsed = () => (elapsed.value = moment(date.value).fromNow());

    watch(date, calcElapsed, { immediate: true });

    if (updateInterval.value > 0) {
      const intervalEH = setInterval(() => calcElapsed(), updateInterval.value);
      onBeforeUnmount(() => clearInterval(intervalEH));
    }

    return {
      humanReadable,
      elapsed,
    };
  },
});
</script>
