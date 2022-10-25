<template>
  <div class="h-100">
    <v-btn block class="rounded-0 h-100" color="red-lighten-2" :to="{name:'donate'}">
      <v-icon class="mr-2 mt-n1" left :icon="mdiHeart" />
      Donate
    </v-btn>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { mdiHeart } from "@mdi/js";

export default defineComponent({
  setup() {
    return {
      mdiHeart,
    };
  },
});
</script>