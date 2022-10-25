<template>
  <section class="overflow-hidden py-8">
    <v-container fluid style="max-width: 800px">
      <div class="d-flex justify-center">
        <v-alert type="info" class="my-4" v-if="metaFiles[0].state !== 'PROCESSED'">
          The sample video will appear here once it has been
          processed.
        </v-alert>

        <Video v-else :src="getFileURL(`${metaFiles[0].releaseUID}/${metaFiles[0].fileName}`, false)" />
      </div>
    </v-container>
  </section>
</template>


<script lang="ts">
import { defineComponent, PropType } from "vue";
import { getFileURL } from "@/utils/url";

export default defineComponent({
  props: {
    metaFiles: {
      type: Array as PropType<IMetaFile[]>,
      required: true,
    },
  },
  setup() {
    return {
      getFileURL,
    };
  },
});
</script>
