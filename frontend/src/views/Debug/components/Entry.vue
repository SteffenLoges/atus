<template>
  <div :class="{ 'text-red': logSeverity === 3 }">
    {{ timeComputed }}
    <span class="text-disabled">|</span> {{ severityComputed }}
    <span class="text-disabled">|</span> {{ logType }}
    <span class="text-disabled">|</span> {{ uid }}
    <span class="text-disabled">|</span> {{ message }}
  </div>
</template>



<script lang="ts">
import { defineComponent, toRefs } from "vue";
import moment from "moment";

export default defineComponent({
  inheritAttrs: false,
  props: {
    logSeverity: {
      type: Number,
      required: true,
    },
    logType: {
      type: String,
      required: true,
    },
    time: {
      type: String,
      required: true,
    },
    uid: {
      type: String,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { time, logSeverity } = toRefs(props);

    const severities: { [key: string]: string } = {
      "0": "DEBUG",
      "1": "INFO",
      "2": "WARNING",
      "3": "ERROR",
      "4": "FATAL",
    };

    return {
      timeComputed: moment(time.value).format("MM/DD/YYYY HH:mm:ss"),
      severityComputed: severities[logSeverity.value],
    };
  },
});
</script>

