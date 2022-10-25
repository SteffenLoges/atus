<template>
  <FormCard :loading="isLoading" title="Edit Categories" @submit="onSubmit">
    <v-card-text>
      <v-alert type="info" class="mb-4">
        Choose wich categories {{ appName }} will download.
        <small>
          <p class="mt-3">
            <strong>Includes:</strong> Release name <u>must</u> contain one of the words in this list. <br />
            <strong>Excludes:</strong> Release name <u>must not</u> contain one of the words in this list.
          </p>
          <p class="mt-2">
            Seperate words with a newline.<br />
            Leave blank to downlaod all files in the category.
            <br />
            <i>All values are case insensitive.</i>
          </p>
        </small>
      </v-alert>

      <v-card variant="text" v-for="(category, i) in categories" :class="{ 'mt-3': i > 0 }" :key="category.name"
        :title="category.name" class="card-accent">
        <v-card-text>
          <Category v-bind="category" @update:enabled="categories[i].enabled = $event"
            @update:includes="categories[i].includes = $event" @update:excludes="categories[i].excludes = $event"
            @update:maxSize="categories[i].maxSize = $event" />
        </v-card-text>
      </v-card>
    </v-card-text>

    <v-card-actions class="px-5 justify-end">
      <v-btn color="primary" type="submit">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref } from "vue";
import useGlobalStore from "@/store/global";
import { send } from "@/utils/websocket";
import { success } from "@/plugins/toast";
import Category from "./components/Category.vue";

export default defineComponent({
  components: {
    Category,
  },
  async setup() {
    const globalStore = useGlobalStore();
    const appName = import.meta.env.VITE_APP_NAME

    let isLoading = ref(false);
    let categories = ref<ICategory[]>([]);

    // --------------------------------------------------------------------------

    const r: IResponse<ICategory[]> = await send("SETTINGS__FILTERS_CATEGORIES__GET_ALL")
    categories.value = r.payload;

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      send("SETTINGS__FILTERS_CATEGORIES__SAVE", { categories: categories.value, })
        .then(() => success("Settings saved successfully"))
        .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
        .finally(() => isLoading.value = false);
    };

    // --------------------------------------------------------------------------

    return {
      appName,
      onSubmit,
      categories,
      isLoading,
    };
  },
});
</script>
