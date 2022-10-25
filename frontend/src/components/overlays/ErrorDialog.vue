<template>
  <v-dialog v-bind="$attrs" persistent :modelValue="true">
    <v-card min-width="300px" class="dialog-card align-self-center">
      <v-toolbar color="red-darken-4" height="40">
        <span class="text-grey-lighten-4 px-5">
          <slot name="title">{{ $attrs.title || "Error" }}</slot>
        </span>
      </v-toolbar>
      <v-card-text>
        <slot>{{ $attrs.text || "An error occurred." }}</slot>
        <div class="mt-4" v-if="$slots.details">
          <small>Error details:</small>
          <v-alert density="compact" max-height="300px" class="overflow-auto text-caption">
            <slot name="details" />
          </v-alert>
        </div>
      </v-card-text>
      <v-card-actions class="justify-center">
        <v-btn text block @click.prevent="onDismiss">OK</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  emits: ["dismiss", "update:modelValue"],
  setup(props, { emit }) {
    const onDismiss = () => {
      emit("dismiss");
      emit("update:modelValue", false);
    };

    return {
      onDismiss,
    };
  },
});
</script>


<style lang="scss" scoped>
@use "vuetify/styles/settings/variables" as v;

.dialog-card {
  min-width: 300px;
  max-width: 100vw;

  @media #{(map-get(v.$display-breakpoints, "lg-and-up"))} {
    max-width: 60vw;
  }
}
</style>