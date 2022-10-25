<template>
  <Card title="NFO" class="overflow-auto" :loading="isLoading" v-bind="$attrs">
    <template #title-actions>
      <v-btn size="small" :icon="mdiMagnifyMinus" :disabled="!isProcessed || fontSize <= 6" @click="fontSize -= 0.5" />
      <v-btn size="small" :icon="mdiMagnifyPlus" :disabled="!isProcessed || fontSize >= 30" @click="fontSize += 0.5" />
      <v-btn :disabled="!isProcessed" :href="`${nfoURL}&download`" size="small" :icon="mdiDownload" />
    </template>

    <v-card-text>
      <v-alert v-if="!isProcessed" type="info">
        NFO is currently being processed.
      </v-alert>
      <NFO v-else-if="nfoData" :data="nfoData" :fontSize="fontSize" @update:fontSize="fontSize = $event" />
    </v-card-text>
  </Card>
</template>


<script lang="ts">
import { defineComponent, PropType, ref, toRefs, watch, computed } from "vue";
import useGlobalStore from "@/store/global";
import { getFileURL } from "@/utils/url";
import { mdiMagnifyPlus, mdiMagnifyMinus, mdiDownload } from "@mdi/js";

export default defineComponent({
  props: {
    metaFiles: {
      type: Array as PropType<IMetaFile[]>,
      required: true,
    },
  },
  setup(props) {
    const { metaFiles } = toRefs(props);
    const globalStore = useGlobalStore();

    const isLoading = ref(false)
    const nfoData = ref("")
    const fontSize = ref(12)
    const isProcessed = computed(() => metaFiles.value[0].state === "PROCESSED")
    const nfoURL = computed(() => getFileURL(`${metaFiles.value[0].releaseUID}/${metaFiles.value[0].fileName}`))

    watch(isProcessed, (v) => {
      if (v) {
        fetch(nfoURL.value)
          .then((res) => res.text())
          .then((text) => nfoData.value = text)
          .catch((err) => globalStore.setError(err))
          .finally(() => isLoading.value = false);
      }
    }, { immediate: true });

    return {
      isLoading,
      nfoURL,
      isProcessed,
      nfoData,
      fontSize,
      getFileURL,
      mdiMagnifyPlus,
      mdiMagnifyMinus,
      mdiDownload,
    };
  },
});
</script>
