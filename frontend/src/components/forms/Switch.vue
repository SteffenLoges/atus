<template>
  <div class="d-flex">
    <v-switch inline color="primary" class="flex-grow-0" hide-details v-bind="$attrs" />

    <div v-if="$attrs.hint || $slots.hint" class="d-flex align-center justify-center">
      <v-tooltip location="bottom">
        <template v-slot:activator="{ props }">
          <v-icon class="hint-icon" v-bind="props" :icon="mdiHelpCircle" />
        </template>
        <span>
          <slot name="hint">{{ $attrs.hint }}</slot>
        </span>
      </v-tooltip>
    </div>
  </div>
</template>


<style lang="scss" scoped>
.hint-icon {
  cursor: help;
  opacity: 0.5;
  margin: 0 0 5px 8px;

  &:hover {
    opacity: 0.8;
  }
}
</style>

<script lang="ts">
import { defineComponent } from "vue";
import { mdiHelpCircle } from "@mdi/js";

export default defineComponent({
  setup() {
    return {
      mdiHelpCircle,
    };
  },
});
</script>
