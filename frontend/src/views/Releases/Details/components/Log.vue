<template>
  <Card title="Log" :loading="isLoading">
    <v-card-text class="pt-0">
      <Pagination :modelValue="page" :pages="pages" size="small" variant="text" @update:modelValue="setPage($event)" />

      <v-alert v-if="!entries.length" type="info">
        No entries found.
      </v-alert>

      <Entry v-for="(entry, i) of entries" :key="i" :showRelease="false" :entry="entry" />
    </v-card-text>
  </Card>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onBeforeUnmount } from "vue";
import useGlobalStore from "@/store/global";
import Entry from "../../../Log/components/Entry.vue";
import useLog from "../../../Log/composables/log";

export default defineComponent({
  components: {
    Entry,
  },
  props: {
    uid: {
      type: String,
      required: true,
    },
  },
  async setup(props) {
    const globalStore = useGlobalStore();
    const { entries, count, getEntries } = useLog();

    const isLoading = ref(true);

    const perPage = 10;
    const page = ref(1);

    const pages = computed(() => Math.ceil(count.value / perPage));

    const setPage = (p: number) => {
      page.value = p;
      _getEntries();
    };

    let refreshInterval: number;
    onBeforeUnmount(() => clearInterval(refreshInterval));

    const getEntriesSilent = () =>
      getEntries({
        releaseUID: props.uid,
        offset: (page.value - 1) * perPage,
        limit: perPage,
        severity: -1,
      })
        .catch(({ payload, statusCode }: IResponse<string>) => globalStore.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => isLoading.value = false);

    const _getEntries = async () => {
      isLoading.value = true;
      await getEntriesSilent();
    };

    await _getEntries();

    refreshInterval = window.setInterval(() => {
      if (page.value === 1) {
        getEntriesSilent();
      }
    }, 5e3);

    return {
      entries,
      page,
      setPage,
      isLoading,
      pages,
    };
  },
});
</script>
