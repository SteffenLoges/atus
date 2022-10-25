<template>
  <v-container class="d-flex flex-column mt-4" style="max-width: 1200px">
    <Filters class="mb-3" v-bind="filters" @update:filter="setFilter($event.name, $event.value)" />
    <div class="mb-4 mx-2">
      <Pagination v-model="filters.page" :pages="pages" size="small" variant="text" />
    </div>

    <v-progress-circular v-if="isLoading" indeterminate class="d-flex mx-auto my-9" color="primary" size="64" />
    <Card v-else variant="flat" color="transparent" class="w-100">
      <SlideDownTransition>
        <v-alert type="info" v-if="!isLoading && releases.length == 0">
          No releases found.
        </v-alert>
      </SlideDownTransition>

      <SlideDownTransition>
        <v-alert variant="tonal" color="primary" type="info" class="mb-4" v-if="sumNewReleases > 0">
          There are {{sumNewReleases}} new releases.
          <v-btn variant="text" size="small" @click="filters.page = 1">Show them</v-btn>
        </v-alert>
      </SlideDownTransition>

      <Release :class="{ 'mt-4': i > 0 }" v-for="(rls, i) of releases" :key="rls.uid" :release="rls" />
    </Card>

    <div class="mt-4 mx-2">
      <Pagination v-model="filters.page" :pages="pages" size="small" variant="text" />
    </div>
  </v-container>
</template>


<script lang="ts">
import { defineComponent, ref, computed, watch, onBeforeMount } from "vue";
import { storeToRefs } from "pinia"
import { useRouter, useRoute } from "vue-router";
import useGlobalStore from "@/store/global";
import { send } from "@/utils/websocket";
import { preserveTypeMerge } from "@/utils/helpers";
import Release from "./components/Release.vue";
import Filters from "./components/Filters.vue";

type IFilterName =
  | "name"
  | "category"
  | "state"
  | "perPage"
  | "page";

export default defineComponent({
  components: {
    Release,
    Filters,
  },
  async setup() {
    const globalStore = useGlobalStore();
    const { sumNewReleases } = storeToRefs(globalStore)
    const router = useRouter();
    const route = useRoute();
    const isLoading = ref(false);
    const releases = ref<IRelease[]>([]);
    const count = ref(0);

    // reset new releases counter on page load.
    // the user just loaded the page, so the releases are in the list
    // ToDo: this is a hacky way to do this. 
    //   - Replace the whole new releases counter with eventhanders in this component / composable
    onBeforeMount(() => globalStore.resetSumNewReleases())

    const filters = ref<{ [key in IFilterName]: any }>({
      name: "",
      page: 1,
      perPage: 10,
      category: "all",
      state: "all"
    });

    const setFilter = (name: IFilterName, value: any) => {
      filters.value.page = 1; // reset page on filter change
      filters.value[name] = value;
    };

    const pages = computed(() => Math.ceil(count.value / (filters.value.perPage as number)))

    const getReleases = async () => send("RELEASES__BROWSE__GET", {
      ...filters.value,
      offset: ((filters.value.page as number) - 1) * (filters.value.perPage as number),
      limit: filters.value.perPage,
    }).then(({ payload }: IResponse<{ count: number; releases: IRelease[] }>) => {
      releases.value = payload.releases || [];
      count.value = payload.count || 0;
      isLoading.value = false;
    });

    watch(filters, () => {
      router.push({ query: { ...route.query, ...filters.value } });
    }, { deep: true, immediate: true });

    watch(route, () => {
      isLoading.value = true
      filters.value = preserveTypeMerge(filters.value, route.query as { [key in IFilterName]: any })
      getReleases()
    }, { deep: true, immediate: true });

    // refresh list when new release is added, we are on the first page and there are no filters set
    watch(sumNewReleases, () => {
      if (sumNewReleases.value > 0 && filters.value.page == 1 && filters.value.state === "all" && filters.value.category === "all") {
        globalStore.resetSumNewReleases()
        getReleases()
      }
    })

    return {
      releases,
      isLoading,
      filters,
      pages,
      setFilter,
      sumNewReleases,
    };
  },
});
</script>
