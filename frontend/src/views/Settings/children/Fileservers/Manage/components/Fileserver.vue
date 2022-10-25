<template>
  <v-card density="compact" variant="flat" class="rounded-0 wrapper-card px-2 card-accent">
    <v-card-title class="d-flex align-center">
      <div class="flex-grow-1 ml-4 text-truncate">{{ name }}</div>
      <v-spacer />
      <div class="d-flex">
        <div class="d-flex">
          <v-btn class="mr-2" :class="`text-${enabled ? 'red' : 'green'}-lighten-2`" :icon="enabled ? mdiStop : mdiPlay"
            size="small" :disabled="isStopping" flat @click.prevent="$emit('toggle', !enabled)" />
          <div>
            <v-tooltip location="bottom" :disabled="!enabled" text="You can't modify a fileserver while it's running">
              <template v-slot:activator="{ props }">
                <div v-bind="props">
                  <v-btn :disabled="enabled" :to="{ name: 'settings_fileservers_edit', params: { uid } }" class="mr-1"
                    size="small" :icon="mdiPencil" flat />
                  <v-btn :disabled="enabled" @click.prevent="$emit('delete')" size="small" :icon="mdiDelete" flat />
                </div>
              </template>
            </v-tooltip>
          </div>
        </div>
      </div>
    </v-card-title>
    <v-card-text style="font-size: 0.75em !important">
      <v-row no-gutters class="statistics mt-lg-2">
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Server Load</span>
            <span class="text-high-emphasis">
              <span v-if="!serverLoad">---</span>
              <span v-else v-for="(load, i) of serverLoad" :key="i" :class="{ 'ml-2': i > 0 }">
                {{ load.toFixed(2) }}
              </span>
            </span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Storage</span>
            <span class="text-high-emphasis">
              <template v-if="!diskFreeSpace || !diskTotalSpace">---</template>
              <template v-else>
                {{ bytesHumanReadable(diskFreeSpace) }}
                /
                {{ bytesHumanReadable(diskTotalSpace) }}
              </template>
            </span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Uploaded Releases</span>
            <span class="text-high-emphasis">{{ filesDownloaded.toLocaleString() }}</span>
          </div>
        </v-col>
        <v-col cols="12" lg="3" class="text-lg-center">
          <div class="inner">
            <span class="text-medium-emphasis d-lg-block title">Status</span>
            <span :class="status.class" v-text="status.text"></span>
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>



<script lang="ts">
import { defineComponent, PropType, ref, toRefs, computed, watch } from "vue";
import { bytesHumanReadable } from "@/utils/conversion";
import { mdiDelete, mdiPencil, mdiPlay, mdiStop } from "@mdi/js";

export default defineComponent({
  props: {
    uid: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    enabled: {
      type: Boolean,
      required: true,
    },
    filesDownloaded: {
      type: Number,
      required: true,
    },
    serverLoad: {
      type: Array as PropType<number[]>,
    },
    diskFreeSpace: {
      type: Number,
    },
    diskTotalSpace: {
      type: Number,
    },
  },
  emits: ["delete", "toggle"],
  setup(props, { emit }) {
    const { enabled } = toRefs(props);
    const isStopping = ref(false);

    watch(enabled, () => (isStopping.value = false));

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
      onToggle,
      status,
      isStopping,
      bytesHumanReadable,
      mdiDelete,
      mdiPencil,
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