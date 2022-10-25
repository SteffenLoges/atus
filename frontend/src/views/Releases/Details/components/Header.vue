<template>
  <header>
    <div class="backdrop" :style="{ backgroundImage: `url(${coverURL})` }"></div>
    <div class="content-wrapper d-flex align-center py-15">
      <v-container class="content" fluid>
        <div class="release-name mb-5 mb-lg-3 text-center text-lg-left">
          <ReleaseName :name="release.name" />
        </div>
        <div class="d-flex flex-column flex-md-row">
          <div class="flex-grow-0 mt-lg-n11 d-flex justify-center mb-5 mb-lg-0" style="flex-basis: 330px">
            <div class="d-flex justify-center align-center cover-image mr-md-9">
              <img :src="coverURL" />
            </div>
          </div>
          <div class="d-flex flex-column flex-grow-1">
            <div class="text-caption text-medium-emphasis">
              <span class="d-inline-block" style="min-width: 45px">Found:</span>
              <DateTimeLive :date="release.added" />
              on {{ release.sourceName }}
            </div>
            <div class="mt-1 text-caption text-medium-emphasis">
              <span class="d-inline-block" style="min-width: 45px">Pre:</span>
              <DateTimeLive :date="release.pre" />
            </div>

            <div class="mt-5">
              <Category :name="release.category" :raw="release.categoryRaw" />
            </div>

            <div class="mt-4">
              <Size :value="release.size" />
            </div>

            <div class="mt-8" v-if="metaFiles">
              <ProgressChips :size="$vuetify.display.xlAndUp ? 'default' : 'small'" :metaFiles="metaFiles" />
            </div>

            <div class="mt-8 text-caption text-high-emphasis">
              Assigned fileserver: {{ assignedServer }}
            </div>

            <div class="mt-3">
              <State :state="state" :uploadDate="state.uploadDate" :downloadState="downloadState" />
            </div>

            <div class="mt-6 d-flex flex-column flex-sm-row">
              <v-btn v-if="torrentURL" :href="torrentURL" variant="tonal" size="small" class="mb-1 mb-sm-0">
                Download .torrent
              </v-btn>

              <v-btn variant="tonal" size="small" class="mb-1 mb-sm-0 ml-sm-2"
                :disabled="progress !== 100 || uploadInProgress" :loading="uploadInProgress" @click="$emit('upload')">
                {{ state.state === "UPLOADED" ? "Re-" : "" }}Upload
              </v-btn>

              <v-btn variant="tonal" size="small" class="mb-1 mb-sm-0 ml-sm-2" @click="$emit('delete')"
                :disabled="uploadInProgress">Delete</v-btn>
            </div>
          </div>
        </div>
      </v-container>
    </div>

    <ProgressBar :progress="progress" :height="6" class="progress-bar" />
  </header>
</template>


<script lang="ts">
import { defineComponent, PropType, toRefs, computed } from "vue";
import useMetaFiles from "../../composables/metaFiles";
import ReleaseName from "../../components/ReleaseName.vue";
import ProgressBar from "../../components/ProgressBar.vue";
import DateTimeLive from "../../components/DateTimeLive.vue";
import Category from "../../components/Category.vue";
import Size from "../../components/Size.vue";
import ProgressChips from "../../components/ProgressChips.vue";
import State from "../../components/State.vue";
import { getFileURL } from "@/utils/url";

export default defineComponent({
  components: {
    ReleaseName,
    ProgressBar,
    DateTimeLive,
    Category,
    Size,
    ProgressChips,
    State,
  },
  props: {
    release: {
      type: Object as PropType<IRelease>,
      required: true,
    },
    state: {
      type: Object as PropType<IReleaseState>,
      required: true,
    },
    downloadState: {
      type: Object as PropType<IDownloadState>,
    },
    metaFiles: {
      type: Array as PropType<IMetaFile[]>,
      required: true,
    },
    uploadInProgress: {
      type: Boolean,
      required: true,
    },
  },
  emits: ["delete", "upload"],
  setup(props) {
    const { metaFiles, downloadState, state, release } = toRefs(props);
    const { getCoverImage } = useMetaFiles();

    const coverURL = computed(() => getCoverImage(metaFiles.value));

    const progress = computed(() => {
      if (["UPLOADED", "DOWNLOADED"].includes(state.value.state)) {
        return 100;
      }

      return downloadState.value?.done || 0;
    });

    const torrentURL = computed(() => {
      const mf = metaFiles.value.find((file) => file.type === "TORRENT" && file.state === "PROCESSED")
      if (!mf) {
        return undefined;
      }

      return getFileURL(`${props.release.uid}/release.torrent`) + "&download"
    });

    const assignedServer = computed(() => {
      if (release.value.fileserverName) {
        return release.value.fileserverName;
      }

      if (state.value.state === "NEW") {
        return "To be determined";
      }

      return "Server unreachable";
    });

    return {
      progress,
      coverURL,
      torrentURL,
      assignedServer
    };
  },
});
</script>


<style lang="scss" scoped>
@use "vuetify/styles/settings/variables" as v;
$theme-color: #07091a;

.cover-image {
  object-fit: contain;
  max-height: 400px;
  max-width: 100%;

  img {
    max-height: 400px;
    max-width: 100%;
    border-radius: 5px;
    box-shadow: 0 3px 6px rgba(0, 0, 0, 0.16), 0 3px 6px rgba(0, 0, 0, 0.23);
  }
}

header {
  overflow: hidden;
  position: relative;

  .backdrop {
    position: absolute;
    $offset: -100px; // removes shimmering effect caused by blur
    top: $offset;
    left: $offset;
    right: $offset;
    bottom: $offset;
    filter: blur(10px);
    background: no-repeat center 40%;
    background-size: 100% auto;
  }

  .content-wrapper {
    position: relative;
    height: 100%;
    background-image:
      linear-gradient(to bottom, rgba($theme-color, 0.6) 10vh, $theme-color ),
      linear-gradient(to right, rgba($theme-color, 0.7) 20%, transparent 80%, rgba($theme-color, 0.7));
  }

  .content {
    padding: 0 8vw;
  }

  .progress-bar {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
  }
}

.release-name {
  @media #{(map-get(v.$display-breakpoints, "lg-and-up"))} {
    padding-left: 330px;
  }
}
</style>