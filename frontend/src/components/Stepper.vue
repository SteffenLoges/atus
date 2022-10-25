<template>
  <div class="wrapper option-1 option-1-1">
    <ol class="c-stepper">
      <li class="c-stepper__item" :class="{ active: activeStep >= i }" v-for="(step, i) of steps" :key="step">
        <h3 class="c-stepper__title">{{ step }}</h3>
      </li>
    </ol>
  </div>
</template>
 

<script lang="ts">
import { defineComponent, PropType } from "vue";

export default defineComponent({
  props: {
    steps: {
      type: Array as PropType<string[]>,
      required: true,
    },
    activeStep: {
      type: Number,
      required: true,
    },
  },
});
</script>



<style lang="scss" scoped>
// credit: https://ishadeed.com/article/stepper-component-html-css/
$active-color: #551fd3;

.wrapper {
  --circle-size: clamp(1.5rem, 5vw, 2.5rem);
  --spacing: clamp(0.25rem, 2vw, 0.5rem);
}

.c-stepper {
  display: flex;

  .c-stepper__item {
    display: flex;
    flex-direction: column;
    flex: 1;
    text-align: center;

    .c-stepper__title {
      font-weight: 300;
      opacity: 0.4;
    }

    &:before {
      content: "";
      display: block;
      width: var(--circle-size);
      height: var(--circle-size);
      border-radius: 50%;
      background-color: desaturate($active-color, 50%);
      opacity: 0.3;
      margin: 0 auto 1rem;
      border: 5px solid #9a9a9a;
    }

    &.active:before {
      background-color: $active-color;
    }

    &:not(:first-child) {
      &:after {
        content: "";
        opacity: 0.3;
        position: relative;
        top: calc(var(--circle-size) / 2);
        width: calc(100% - var(--circle-size) - calc(var(--spacing) * 2));
        right: calc(50% - calc(var(--circle-size) / 2 + var(--spacing)));
        height: 2px;
        background-color: grey;
        order: -1;
      }
    }

    &.active:after {
      background-color: $active-color;
    }

    &.active {

      &:before,
      &:after,
      .c-stepper__title {
        opacity: 1;
      }

      .c-stepper__title {
        font-weight: 500;
      }
    }
  }
}
</style>