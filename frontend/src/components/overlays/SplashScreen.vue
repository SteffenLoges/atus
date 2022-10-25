<template>
  <v-overlay :model-value="true" persistent no-click-animation scrim="#1c212e" class="align-center justify-center">
    <v-card color="transparent" flat>
      <v-card-text class="text-center">
        <div class="d-flex justify-center mb-4">
          <v-img src="/src/assets/images/logo-big.png" maxWidth="150" />
        </div>

        <div class="py-4 px-1">
          <v-progress-linear indeterminate color="primary" size="100" />
        </div>
        <div class="text-overline">
          <slot />
        </div>
        <div class="mt-5 btn-row" :class="{ visible: showRealoadBtn }">
          <v-btn size="small" width="150" color="primary" variant="tonal" @click.prevent="reloadPage()">Retry</v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-overlay>
</template>


<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from "vue";

export default defineComponent({
  setup() {
    const showRealoadBtn = ref(false);

    let showRealoadBtnTimeout = -1;
    onMounted(() => showRealoadBtnTimeout = setTimeout(() => showRealoadBtn.value = true, 5e3))
    onBeforeUnmount(() => clearTimeout(showRealoadBtnTimeout))

    return {
      showRealoadBtn,
      reloadPage: () => window.location.reload(),
    };
  },
});
</script>


<style lang="scss" scoped>
.btn-row {
  transition: opacity 0.7s ease-in-out;
  opacity: 0;

  &.visible {
    opacity: 1;
  }
}
</style>
