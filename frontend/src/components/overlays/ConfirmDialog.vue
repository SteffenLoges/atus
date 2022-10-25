<template>
  <v-dialog v-bind="$attrs" persistent>
    <v-card class="dialog-card align-self-center">
      <v-toolbar color="blue-grey-darken-4" height="40">
        <span class="text-disabled text-white px-5">
          <slot name="title">{{ $attrs.title || "Confirm your action" }}</slot>
        </span>
      </v-toolbar>
      <v-card-text>
        <slot>Are you sure?</slot>
      </v-card-text>
      <v-card-actions class="justify-center">
        <v-btn text @click="onConfirm()">Yes</v-btn>
        <v-btn text @click="onCancel()">No</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  emits: ["always", "confirm", "cancel"],
  setup(props, { emit }) {

    const onConfirm = () => {
      emit("always");
      emit("confirm");
    };

    const onCancel = () => {
      emit("always");
      emit("cancel");
    };

    return {
      onConfirm,
      onCancel,
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