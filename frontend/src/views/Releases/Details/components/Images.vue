<template>
  <section>
    <Lightbox :images="imagesComputed" v-model="lightboxIndex" />
    <v-container fluid class="py-8">
      <v-alert type="info" class="mb-4 mt-n2" :class="{ 'mb-6': imagesComputed.length > 0 }"
        v-if="imagesComputed.length < metaFiles.length">
        We are still processing some images
      </v-alert>

      <v-row>
        <v-col v-for="(image, i) of imagesComputed" :key="image.title" cols="12" lg="3"
          class="d-flex justify-center align-center">
          <img :src="image.src" @click="lightboxIndex = i" class="image cursor-pointer" />
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>


<script lang="ts">
import { defineComponent, PropType, ref, computed } from "vue";
import { getFileURL } from "@/utils/url";
import useMetaFiles from "../../composables/metaFiles";

export default defineComponent({
  props: {
    metaFiles: {
      type: Array as PropType<IMetaFile[]>,
      required: true,
    },
  },
  setup(props) {
    const metaFiles = useMetaFiles();
    const lightboxIndex = ref(-1);

    const imagesComputed = computed(() =>
      props.metaFiles
        .filter(({ state }) => state === "PROCESSED")
        .map(({ releaseUID, fileName, type }) => ({
          src: getFileURL(releaseUID + "/" + fileName),
          title: `${metaFiles.getName(type)} - ${fileName}`,
        }))
    );

    return {
      imagesComputed,
      lightboxIndex,
    };
  },
});
</script>


<style lang="scss" scoped>
.image {
  max-height: 300px;
  max-width: 100%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.16),
    0 3px 6px rgba(0, 0, 0, 0.23);
}
</style>