<template>
  <ConfirmLeaveDialog v-model="formDirty" v-if="!noConfirm" />
  <Card v-bind="$attrs" :loadingText="loadingText">
    <template v-if="$slots['title']" #title>
      <slot name="title" />
    </template>
    <v-form @submit.prevent="submit" @change="formDirty = true" ref="form">
      <slot />
    </v-form>
  </Card>
</template>



<script lang="ts">
import { defineComponent, ref } from "vue";
import { error } from "@/plugins/toast";

export default defineComponent({
  props: {
    loadingText: {
      type: String,
      default: "",
    },
    noConfirm: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["submit"],
  setup(props, { emit }) {
    const formDirty = ref(false);
    const form = ref<any>(null);

    const submit = async (e: any) => {
      if (form.value) {
        let v = await form.value.validate();
        if (!v.valid) {
          error("Form is not valid", "Please check your inputs and try again.");
          return;
        }
      }

      formDirty.value = false;
      emit("submit", e);
    };

    return {
      form,
      formDirty,
      submit,
    };
  },
});
</script>