<template>
  <v-progress-linear :modelValue="progress" :color="barColor" :striped="progress < 100" :height="height" reverse
    class="progress" />
</template>


<script lang="ts">
import { defineComponent, toRefs, computed } from "vue";

export default defineComponent({
  props: {
    progress: {
      type: Number,
      default: 0,
    },
    height: {
      type: Number,
      default: 10,
    },
  },
  setup(props) {
    const { progress } = toRefs(props);

    const barColor = computed(() => {
      if (progress.value >= 70) {
        return "green-lighten-1";
      }

      if (progress.value >= 30) {
        return "orange-lighten-2";
      }

      return "red-lighten-2";
    });

    return {
      barColor,
    };
  },
});
</script>



<style lang="scss" scoped>
.progress {
  // flip so the bar moves forward
  transform: scaleX(-1) scaleY(-1);
}
</style>