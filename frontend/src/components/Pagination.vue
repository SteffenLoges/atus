<template>
  <v-pagination :modelValue="modelValue" :length="pages" :totalVisible="lgAndUp ? 10 : 5"
    @update:modelValue="$emit('update:modelValue', $event)" />
</template>
 
<script lang="ts">
import { defineComponent } from "vue";
import { useDisplay } from "vuetify";

export default defineComponent({
  props: {
    pages: {
      type: Number,
      required: true,
    },
    modelValue: {
      type: Number,
      required: true,
    },
  },
  emits: ["update:modelValue"],
  setup() {
    const { lgAndUp } = useDisplay();

    return {
      lgAndUp,
    };
  },
});
</script>




