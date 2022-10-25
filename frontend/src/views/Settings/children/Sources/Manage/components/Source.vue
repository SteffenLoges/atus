<template>
  <v-card density="compact" variant="flat" class="rounded-0 wrapper-card px-2 card-accent">
    <v-card-title class="d-flex align-center">
      <div class="d-flex align-center" style="width: 22px">
        <img v-if="favicon" :src="getFileURL(`favicons/${favicon}`)" style="max-width: 20px; max-height: 20px" />
      </div>

      <div class="flex-grow-1 ml-2 text-truncate">{{ name }}</div>
      <v-spacer />
      <div class="d-flex">
        <v-btn class="mr-2" :class="`text-${enabled ? 'red' : 'green'}`" :icon="enabled ? mdiStop : mdiPlay"
          size="small" :disabled="isStopping" flat @click.prevent="onToggle()" />
        <div>
          <v-tooltip location="bottom" :disabled="!enabled" text="You can't modify a source while it's running">
            <template v-slot:activator="{ props }">
              <div v-bind="props">
                <v-btn :disabled="enabled" :to="{ name: 'settings_sources_edit', params: { uid } }" class="mr-1"
                  size="small" :icon="mdiPencil" flat />
                <v-btn :disabled="enabled" @click.prevent="$emit('delete')" size="small" :icon="mdiDelete" flat />
              </div>
            </template>
          </v-tooltip>
        </div>
      </div>
    </v-card-title>
    <v-card-text style="font-size: 0.75em !important">
      <v-row no-gutters class="statistics mt-lg-2">
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Times checked</span>
            <span class="text-high-emphasis">{{ timesChecked.toLocaleString() }}</span>
          </div>
        </v-col>

        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Last checked</span>
            <span class="text-high-emphasis">{{ lastCheckedHumanized }}</span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Next check</span>
            <span class="text-high-emphasis">{{ enabled ? nextCheckHumanized : "---" }}</span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Status</span>
            <span :class="status.class" v-text="status.text"></span>
          </div>
        </v-col>
      </v-row>
      <v-row no-gutters class="statistics mt-lg-2">
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Downloaded Releases</span>
            <span class="text-high-emphasis">{{ sumReleasesDownloaded.toLocaleString() }}</span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Downloaded .torrent Files</span>
            <span class="text-high-emphasis">{{ sumTorrentsDownloaded.toLocaleString() }}</span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Downloaded Images</span>
            <span class="text-high-emphasis">{{ sumImagesDownloaded.toLocaleString() }}</span>
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>



<script lang="ts">
import { defineComponent, toRefs, computed, ref, watch } from "vue";
import { getFileURL } from "@/utils/url";
import { isValid } from "@/utils/date";
import moment from "moment";
import { mdiPencil, mdiDelete, mdiPlay, mdiStop } from "@mdi/js";

export default defineComponent({
  props: {
    uid: {
      type: String,
      required: true,
    },
    favicon: {
      type: String,
      default: "",
    },
    name: {
      type: String,
      required: true,
    },
    enabled: {
      type: Boolean,
      required: true,
    },
    timesChecked: {
      type: Number,
      required: true,
    },
    sumTorrentsDownloaded: {
      type: Number,
      required: true,
    },
    sumImagesDownloaded: {
      type: Number,
      required: true,
    },
    sumReleasesDownloaded: {
      type: Number,
      required: true,
    },
    lastChecked: {
      type: String,
      required: true,
    },
    nextCheck: {
      type: String,
      required: true,
    },
  },
  emits: ["delete", "toggle"],
  setup(props, { emit }) {
    const { lastChecked, nextCheck, enabled } = toRefs(props);

    const lastCheckedHumanized = computed(() => {
      if (!isValid(lastChecked.value)) {
        return "Never";
      }

      return moment(lastChecked.value).fromNow();
    });

    const nextCheckHumanized = computed(() => {
      if (!isValid(nextCheck.value)) {
        return "---";
      }

      return moment(nextCheck.value).fromNow();
    });

    const isStopping = ref(false);

    watch(enabled, () => isStopping.value = false)

    const onToggle = () => {
      if (enabled.value) {
        isStopping.value = true;
      }

      emit("toggle", !enabled.value);
    };

    const status = computed(() => {
      if (isStopping.value && enabled.value) {
        return { text: "Stopping", class: "text-orange" };
      }

      if (enabled.value) {
        return { text: "Running", class: "text-green" };
      }

      return { text: "Stopped", class: "text-red" };
    });

    return {
      getFileURL,
      onToggle,
      isStopping,
      status,
      lastCheckedHumanized,
      nextCheckHumanized,
      mdiPencil,
      mdiDelete,
      mdiPlay,
      mdiStop,
    };
  },
});
</script>


<style lang="scss" scoped>
@use "vuetify/styles/settings/variables" as v;

.wrapper-card {
  border-width: 1px 0;
}

.statistics {
  .inner {
    background: darken(variables.$accent-color,
        7%) !important;
    border-radius: 5px;
    padding: 5px 10px;
    margin: 2px 0;

    @media #{(map-get(v.$display-breakpoints, "lg-and-up"))} {
      padding-left: 0;
      padding-right: 0;
      margin-left: 7px;
      margin-right: 7px;
    }
  }

  .title {
    display: inline-block;
    min-width: 160px;
    font-size: 1.1em;
    font-weight: 500;

    @media #{(map-get(v.$display-breakpoints, "lg-and-up"))} {
      min-width: auto;
    }
  }
}
</style>