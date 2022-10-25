<template>
  <h2>{{ name }}</h2>
</template>
  

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
  },
});
</script>


<style lang="scss" scoped>
h2 {
  font-size: clamp(0.9rem, 2vw, 1.1rem);
  line-height: 23px;
}
</style>