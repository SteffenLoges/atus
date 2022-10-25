<template>
  <div>
    <Chip style="min-width: 85px" class="justify-center" color="grey" variant="outlined"
      :text="bytesHumanReadable(value, false, 2)" v-bind="$attrs" />
  </div>
</template>


<script lang="ts">
import { defineComponent } from "vue";
import { bytesHumanReadable } from "@/utils/conversion";

export default defineComponent({
  props: {
    value: {
      type: Number,
      required: true,
    },
  },
  setup() {
    return {
      bytesHumanReadable,
    };
  },
});
</script>
