<template>
  <div class="d-flex">
    <Chip style="min-width: 85px" label color="blue-lighten-3 justify-center" variant="flat" :text="name"
      v-bind="$attrs" />

    <Chip style="min-width: 85px" class="ml-2 justify-center" color="blue-lighten-4" variant="flat" :text="raw"
      v-bind="$attrs" />
  </div>
</template>


<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  inheritAttrs: false,
  props: {
    name: {
      type: String,
      required: true,
    },
    raw: {
      type: String,
      required: true,
    },
  },
});
</script>
