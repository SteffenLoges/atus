<template>
  <v-tooltip location="bottom" :disabled="!$slots.tooltip && !tooltip">
    <template v-slot:activator="{ props }">
      <v-chip v-bind="{ ...props, ...$attrs }" class="chip justify-center"
        :variant="(($attrs.variant || 'flat') as any)" :label="label" :disabled="disabled || loading"
        :to="(!disabled && to) as RouteLocationRaw">
        <v-icon v-if="iconComputed" :icon="iconComputed" left :class="{ 'icon-spinner': loading }" />
        <slot>{{ text }}</slot>
        <v-badge v-if="badge" color="primary" class="ml-1 mr-n2" :content="badge" inline />
      </v-chip>
    </template>
    <span>
      <slot name="tooltip">{{ tooltip }}</slot>
    </span>
  </v-tooltip>
</template>

<script lang="ts">
import { defineComponent, computed, toRefs, PropType } from "vue";
import { RouteLocationRaw } from "vue-router";
import { mdiLoading } from "@mdi/js";

export default defineComponent({
  props: {
    text: {
      type: String,
      default: "",
    },
    label: {
      type: Boolean,
      default: true,
    },
    icon: {
      type: String,
      default: "",
    },
    loading: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: [String, Boolean],
      default: false,
    },
    to: {
      type: [String, Object] as PropType<RouteLocationRaw>,
    },
    badge: {
      type: [String, Number],
      default: "",
    },
  },
  setup(props) {
    const { icon, loading } = toRefs(props);

    const iconComputed = computed(() => {
      if (loading.value) {
        return mdiLoading;
      }

      return icon.value;
    });

    return {
      iconComputed,
      loading,
    };
  },
});
</script>


<style lang="scss" scoped>
.chip {
  pointer-events: all !important; // required for v-tooltip to work
  user-select: none;

  .v-icon {
    margin: 0 5px 0 -4px;
  }
}
</style>