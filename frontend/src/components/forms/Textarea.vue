<template>
  <!-- Info: auto-grow is broken in vuetify 3.0.0 -->
  <v-textarea v-bind="$attrs" :rules="computedRules">
    <slot />
  </v-textarea>
</template>



<script lang="ts">
import { defineComponent, computed, toRefs } from "vue";

export default defineComponent({
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { rules, required } = toRefs(props);

    const computedRules = computed<any[]>(() => {
      let r = rules.value;

      if (required.value) {
        r = [(v: any) => !!v || "Value is required", ...r];
      }

      return r;
    });

    return {
      required,
      computedRules,
    };
  },
});
</script>