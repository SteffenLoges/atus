<template>
  <v-text-field v-bind="$attrs" :rules="computedRules" :type="type" ref="input">
    <template v-if="$slots.prepend" #prepend>
      <slot name="prepend" />
    </template>
    <template v-if="$slots['prepend-inner']" #prepend-inner>
      <slot name="prepend-inner" />
    </template>
    <template v-if="$slots['append-inner'] || copyable" #append-inner>
      <slot name="append-inner" />

      <v-tooltip location="bottom">
        <template v-slot:activator="{ props }">
          <v-btn v-if="copyable" @click="copyToClipboard" :icon="mdiContentCopy" color="primary" size="x-small"
            class="mt-n1" v-bind="props" />
        </template>
        <span>Copy to clipboard</span>
      </v-tooltip>
    </template>
    <template v-if="$slots.append" #append>
      <slot name="append" />
    </template>
  </v-text-field>
</template>


<script lang="ts">
import { defineComponent, computed, toRefs, ref } from "vue";
import { success, error } from "@/plugins/toast";
import { mdiContentCopy } from "@mdi/js";

export default defineComponent({
  inheritAttrs: false,
  props: {
    rules: {
      type: Array,
      default: () => [],
    },
    min: {
      type: Number,
      default: -Infinity,
    },
    max: {
      type: Number,
      default: Infinity,
    },
    type: {
      type: String,
      default: "text",
    },
    required: {
      type: Boolean,
      default: false,
    },
    copyable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { rules, required, type, min, max } = toRefs(props);
    const input = ref<any>(null);

    const computedRules = computed<any[]>(() => {
      let r = rules.value;

      if (required.value) {
        r.push((v: any) => !!v || v === 0 || "Value is required");
      }

      if (type.value === "number") {
        if (min.value !== -Infinity) {
          r.push((v: any) => v >= min.value || `Value must be at least ${min.value}`);
        }

        if (max.value !== Infinity) {
          r.push((v: any) => v <= max.value || `Value must be at most ${max.value}`);
        }
      }

      return r;
    });



    // @see https://stackoverflow.com/a/71876238
    const unsecuredCopyToClipboard = (text: string) => {
      const textArea = document.createElement("textarea");
      textArea.value = text; document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      try {
        document.execCommand('copy')
        success("Copied to clipboard")
      } catch (err) {
        console.error('Failed to copy!', err);
        error("Failed to copy to clipboard")
      }
      document.body.removeChild(textArea)
    };

    /**
     * Copies the text passed as param to the system clipboard
     * Check if using HTTPS and navigator.clipboard is available
     * Then uses standard clipboard API, otherwise uses fallback
    */
    const copyToClipboard = () => {
      if ((input?.value?.value ?? "") === "") {
        return;
      }

      if (window.isSecureContext && navigator.clipboard) {
        navigator.clipboard
          .writeText(input.value.value)
          .then(() => success("Copied to clipboard"))
          .catch((err) => error("Failed to copy to clipboard", err));

        return
      }

      unsecuredCopyToClipboard(input.value.value);
    };

    return {
      input,
      computedRules,
      copyToClipboard,
      mdiContentCopy,
    };
  },
});
</script>
