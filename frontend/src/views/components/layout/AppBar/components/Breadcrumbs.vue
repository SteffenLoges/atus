<template>
  <v-breadcrumbs class="breadcrumbs">
    <transition-group name="list" @before-leave="beforeLeave">
      <template v-for="(item, i) of items" :key="item.meta.title">
        <v-breadcrumbs-item :to="item.path" :disabled="i === items.length - 1" :title="(item.meta.title as string)"
          class="breadcrumbs-item" />
        <v-breadcrumbs-divider v-if="i !== items.length - 1" :key="`${item.meta.title}-divider`">
          <v-icon :icon="mdiChevronRight" />
        </v-breadcrumbs-divider>
      </template>
    </transition-group>
  </v-breadcrumbs>
</template>


<script lang="ts">
import { defineComponent, computed } from "vue";
import { useRoute } from "vue-router";
import { mdiChevronRight } from "@mdi/js";

export default defineComponent({
  setup() {
    const route = useRoute();
    const items = computed(() => route.matched.filter((r) => r.meta.title));

    return {
      items,
      beforeLeave: (el: any) => {
        const { marginLeft, marginTop, width, height } = window.getComputedStyle(el);
        el.style.left = `${el.offsetLeft - parseFloat(marginLeft)}px`;
        el.style.top = `${el.offsetTop - parseFloat(marginTop)}px`;
        el.style.width = width;
        el.style.height = height;
      },
      mdiChevronRight,
    };
  },
});
</script>

<style lang="scss" scoped>
.breadcrumbs {
  user-select: none;
  height: 35px !important;
  padding: 0;

  .breadcrumbs-item {
    height: 100%;
  }

  .v-breadcrumbs-item--disabled {
    opacity: 0.6;
  }
}

.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  position: absolute;
  opacity: 0;
  transform: translateX(30px);
}
</style>