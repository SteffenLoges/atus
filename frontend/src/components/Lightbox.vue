<template>
  <VueEasyLightbox class="lightbox" :visible="modelValue > -1" :imgs="images" :index="modelValue"
    @hide="$emit('update:modelValue', -1)" />
</template>
  
<script lang="ts">
import { defineComponent, PropType } from "vue";
import VueEasyLightbox from "vue-easy-lightbox";
import { PropsImgs } from "vue-easy-lightbox/types/types/index";

export default defineComponent({
  components: {
    VueEasyLightbox,
  },
  props: {
    images: {
      type: Array as PropType<PropsImgs>,
    },
    modelValue: {
      type: Number,
      default: -1,
    },
  },
  emits: ["update:modelValue"],
});
</script>
  

<style lang="scss">
// hide rotate buttons. showing them without actually rotating the image on the server is pointless imo
.lightbox .toolbar-btn__rotate {
  display: none;
}
</style>