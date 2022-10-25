<template>
  <v-card v-bind="$attrs" class="custom-card">
    <div class="d-flex" v-if="title || $slots['title'] || $slots['title-actions']">
      <v-card-title class="flex-grow-1">
        <slot name="title">{{ title }}</slot>
      </v-card-title>
      <v-card-actions class="justify-end text-disabled align-center">
        <slot name="title-actions" />
      </v-card-actions>
    </div>

    <v-overlay v-if="hasLoadingAttr" :model-value="($attrs.loading as boolean)" persistent contained
      class="justify-center align-center overlay-loading">
      <v-progress-circular indeterminate class="d-flex mx-auto" color="primary" size="64" />
      <div v-if="loadingText != '' || $slots['loading-text']" class="mt-8 text-subtitle-2 text-high-emphasis">
        <slot name="loading-text">{{ loadingText }}</slot>
      </div>
    </v-overlay>
    <slot />
  </v-card>
</template>



<script lang="ts">
import { defineComponent, toRefs } from "vue";

export default defineComponent({
  props: {
    loadingText: {
      type: String,
      default: "",
    },
    title: {
      type: String,
      default: "",
    },
  },
  setup(props, { attrs }) {
    const { loadingText } = toRefs(props);
    const hasLoadingAttr = "loading" in attrs;

    return {
      loadingText,
      hasLoadingAttr,
    };
  },
});
</script>



<style lang="scss">
.custom-card {
  .v-card__loader {
    display: none;
  }

  .overlay-loading .v-overlay__scrim {
    backdrop-filter: blur(2px);
    opacity: 0.65 !important;
    background-color: #222 !important;
  }
}
</style>