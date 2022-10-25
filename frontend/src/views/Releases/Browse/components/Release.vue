<template>
  <Card class="release d-flex jusifty-center" elevation="2" flat :to="{
    name: 'releases_details',
    params: { uid: release.uid, name: release.name },
  }">
    <div class="image flex-grow-0 bg-grey-darken-4 d-none d-sm-flex">
      <v-img :src="coverURL" cover />
    </div>
    <div class="content-wrapper d-flex flex-grow-1 flex-column justify-space-between">
      <v-card-text class="content flex-grow-1">
        <ReleaseName class="text-truncate" :name="release.name" />

        <div class="mt-3 text-caption text-medium-emphasis">
          <span class="d-inline-block" style="min-width: 45px">Found:</span>
          <DateTimeLive :date="release.added" />
          <span class="d-none d-sm-inline">
            on {{ release.sourceName }}
          </span>
        </div>
        <div class="mt-1 text-caption text-medium-emphasis">
          <span class="d-inline-block" style="min-width: 45px">Pre:</span>
          <DateTimeLive :date="release.pre" />
        </div>

        <div class="mt-3 d-flex flex-column flex-md-row">
          <Category size="small" :name="release.category" :raw="release.categoryRaw" />

          <Size size="small" :value="release.size" class="mt-1 mt-md-0 ml-md-1" />
        </div>

        <div class="mt-4 d-none d-md-flex">
          <ProgressChips :metaFiles="metaFiles" size="small" />
        </div>

        <div class="mt-3">
          <State :state="state" :downloadState="downloadState" :uploadDate="state.uploadDate" />
        </div>
      </v-card-text>

      <ProgressBar :progress="progress" :height="5" />
    </div>
  </Card>
</template>


<script lang="ts">
import { defineComponent, PropType, onMounted, onBeforeUnmount } from "vue";
import useRelease from "../../composables/release";
import ProgressBar from "../../components/ProgressBar.vue";
import ProgressChips from "../../components/ProgressChips.vue";
import State from "../../components/State.vue";
import Category from "../../components/Category.vue";
import DateTimeLive from "../../components/DateTimeLive.vue";
import ReleaseName from "../../components/ReleaseName.vue";
import Size from "../../components/Size.vue";

export default defineComponent({
  components: {
    ProgressBar,
    ProgressChips,
    State,
    Category,
    Size,
    DateTimeLive,
    ReleaseName,
  },
  props: {
    release: {
      type: Object as PropType<IRelease>,
      required: true,
    },
  },
  setup(props) {
    const { progress, coverURL, state, backgroundImage, metaFiles, downloadState, addEventHandlers, removeEventHandlers } =
      useRelease(props.release.uid, props.release.state, props.release.metaFiles, props.release.downloadState);

    onMounted(() => addEventHandlers());
    onBeforeUnmount(() => removeEventHandlers());

    return {
      progress,
      state,
      coverURL,
      backgroundImage,
      metaFiles,
      downloadState,
    };
  },
});
</script>


<style lang="scss" scoped>
.release {
  height: 230px;

  .image {
    $container-width: 160px;

    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-basis: $container-width;
    min-width: $container-width;
    max-width: $container-width;

    .v-img {
      min-height: 100%;
      max-height: 100%;
      min-width: 100%;
      max-width: 100%;
      transition: all 0.3s ease-in-out;
      // filter: grayscale(25%);
    }
  }

  &:hover .image .v-img {
    transform: scale(1.05);
    // filter: grayscale(0%);
  }

  .content-wrapper {
    min-width: 0;
    color: #fff;

    .content {
      position: relative;
      overflow: hidden;

      &:before {
        position: absolute;
        content: "";
        background-color: #e0e0e0;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;

        background: v-bind(backgroundImage) right center no-repeat;
        background-size: cover;
        filter: blur(4px);
        transform: scale(1.1);
        opacity: 0.2;
        z-index: -1;
      }

      &:after {
        position: absolute;
        content: "";
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: linear-gradient(to right,
            rgba(28, 33, 46, 0.8) 60%,
            rgba(28, 33, 46, 0));
        z-index: -1;
      }
    }
  }
}
</style>