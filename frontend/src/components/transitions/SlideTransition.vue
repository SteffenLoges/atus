<template>
  <Transition name="slide-fade" @before-leave="beforeLeave">
    <slot />
  </Transition>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  setup() {
    return {
      // @see: https://forum.vuejs.org/t/transition-group-leave-transition-w-position-absolute-jumping-to-top-left-flip/12258/4
      beforeLeave: (el: any) => {
        const { marginLeft, marginTop, width, height } = window.getComputedStyle(el);
        el.style.left = `${el.offsetLeft - parseFloat(marginLeft)}px`;
        el.style.top = `${el.offsetTop - parseFloat(marginTop)}px`;
        el.style.width = width;
        el.style.height = height;
      },
    };
  },
});
</script>

<style lang="scss" scoped>
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(100px);
  opacity: 0;
  position: absolute;
}
</style>