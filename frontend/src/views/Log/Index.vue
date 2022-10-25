<template>
  <ConfirmDialog v-model="showDeleteConfirmDialog" @always="showDeleteConfirmDialog = false" @confirm="clearLog()">
    Are you sure?
  </ConfirmDialog>

  <v-container class="d-flex mt-4" style="max-width: 1200px">
    <div class="flex-grow-1">
      <Card :loading="isLoading" title="Log">
        <v-card-text>
          <v-row class="mb-2">
            <v-col cols="12" md="6">
              <v-select :model-value="filters.severity" :items="severitiesComputed" label="Severity" hideDetails
                @update:modelValue="setFilter('severity', $event)" />
            </v-col>
            <v-col cols="12" md="6">
              <v-select :model-value="filters.type" :items="typesComputed" label="Type" hideDetails
                @update:modelValue="setFilter('type', $event)" />
            </v-col>
          </v-row>

          <div class="mb-4">
            <Pagination v-model="filters.page" :pages="pages" size="small" variant="text" />
          </div>

          <v-alert v-if="!entries.length" type="info">No entries found.</v-alert>

          <Entry v-for="(entry, i) of entries" :key="i" :entry="entry" />

          <div class="mt-4">
            <Pagination v-model="filters.page" :pages="pages" size="small" variant="text" />
          </div>
        </v-card-text>

        <v-card-actions class="justify-end">
          <v-btn color="warning" :disabled="isLoading" @click="showDeleteConfirmDialog = true">
            Clear log
          </v-btn>
        </v-card-actions>
      </Card>
    </div>
  </v-container>
</template>


<script lang="ts">
import { computed, defineComponent, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { preserveTypeMerge } from "@/utils/helpers";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import useLog from "./composables/log";
import useGlobalStore from "@/store/global";
import { severities } from "./composables/log";
import Entry from "./components/Entry.vue";

type IFilterName = "perPage" | "page" | "type" | "severity";

export default defineComponent({
  components: {
    Entry,
  },
  async setup() {
    const router = useRouter();
    const route = useRoute();
    const global = useGlobalStore();
    const { entries, types, count, getEntries } = useLog();

    const isLoading = ref(false);
    const showDeleteConfirmDialog = ref(false);

    const filters = ref<{ [key in IFilterName]: any }>({
      page: 1,
      perPage: 25,
      type: "",
      severity: -1,
    });

    const severitiesComputed = computed(() => {
      return [
        { title: "All", value: -1 },
        ...severities.map((s, i) => ({
          title:
            s.text.charAt(0).toUpperCase() +
            s.text.slice(1).toLowerCase(),
          value: i,
        })),
      ];
    });

    const typesComputed = computed(() => {
      return [
        { title: "All", value: "" },
        ...types.value.map((t: string) => ({
          title: (
            t.charAt(0).toUpperCase() +
            t.slice(1).toLowerCase()
          ).replace(/_/g, " "),
          value: t,
        })),
      ];
    });

    const setFilter = (name: IFilterName, value: any) => {
      filters.value.page = 1
      filters.value[name] = value
    }

    const pages = computed(() => Math.ceil(count.value / (filters.value.perPage as number)))

    const clearLog = () => {
      isLoading.value = true;

      send("LOG__CLEAR")
        .then(() => {
          success("Log cleared");
          _getEntries();
        })
        .catch((err) => global.setError(err))
        .finally(() => isLoading.value = false)
    };

    const _getEntries = async () => {
      isLoading.value = true;

      await getEntries({
        offset: ((filters.value.page as number) - 1) * (filters.value.perPage as number),
        limit: filters.value.perPage,
        type: filters.value.type,
        severity: filters.value.severity,
      })
        .catch(({ payload, statusCode }: IResponse<string>) => global.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => isLoading.value = false)
    };

    watch(filters, () => {
      router.push({ query: { ...route.query, ...filters.value } });
    }, { deep: true, immediate: true });

    watch(route, () => {
      isLoading.value = true
      filters.value = preserveTypeMerge(filters.value, route.query as { [key in IFilterName]: any })
      _getEntries()
    }, { deep: true, immediate: true });

    return {
      entries,
      filters,
      severitiesComputed,
      typesComputed,
      setFilter,
      clearLog,
      pages,
      count,
      isLoading,
      showDeleteConfirmDialog,
      severities,
      types,
    };
  },
});
</script>
