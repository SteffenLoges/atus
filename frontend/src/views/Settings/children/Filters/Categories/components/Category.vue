<template>
  <Switch v-model="enabledComputed" label="Enabled" />

  <VSlideYTransition>
    <div v-if="enabledComputed">
      <v-row no-gutters>
        <v-col>
          <TextField v-model="maxSizeComputed" type="number" :min="0" required label="Maximum size of a release in GiB"
            hint="Use 0 to disable this filter" persistent-hint class="mb-2" />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" lg="6">
          <Textarea hide-details v-model="includesComputed" placeholder="e.g.&#10;1080p&#10;720p&#10;bluray" :rows="4"
            label="Includes" />
        </v-col>
        <v-col cols="12" lg="6">
          <Textarea hide-details v-model="excludesComputed" placeholder="e.g.&#10;subfrench&#10;subbed&#10;anime"
            :rows="4" label="Excludes" />
        </v-col>
      </v-row>
    </div>
  </VSlideYTransition>
</template>


<script lang="ts">
import { defineComponent, PropType, toRefs, computed } from "vue";

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
    enabled: {
      type: Boolean,
      required: true,
    },
    includes: {
      type: Array as PropType<string[]>,
      required: true,
    },
    excludes: {
      type: Array as PropType<string[]>,
      required: true,
    },
    maxSize: {
      type: Number,
      required: true,
    },
  },
  emits: [
    "update:enabled",
    "update:includes",
    "update:excludes",
    "update:maxSize",
  ],
  setup(props, { emit }) {
    const { enabled, includes, excludes } = toRefs(props);

    const enabledComputed = computed({
      get: () => enabled.value,
      set: (v: boolean) => emit("update:enabled", v),
    });

    const includesComputed = computed({
      get: () => includes.value.join("\n"),
      set: (v: string) => emit("update:includes", v
        .split("\n")
        .filter((s) => s.length > 0)
        .map((s) => s.trim().toLowerCase())
      ),
    });

    const excludesComputed = computed({
      get: () => excludes.value.join("\n"),
      set: (v: string) => emit("update:excludes", v
        .split("\n")
        .filter((s) => s.length > 0)
        .map((s) => s.trim().toLowerCase())
      ),
    });

    const maxSizeComputed = computed({
      get: () => props.maxSize,
      set: (v: number) => emit("update:maxSize", parseInt("" + v)),
    });

    return {
      enabledComputed,
      includesComputed,
      excludesComputed,
      maxSizeComputed,
    };
  },
});
</script>
