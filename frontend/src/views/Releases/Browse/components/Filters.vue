<template>
  <v-card class="d-flex pa-2 browse-filters-wrapper flex-column flex-md-row">
    <div class="flex-grow-1">
      <TextField bgColor="blue-grey-darken-4" label="Release Name" clearable singleLine hideDetails variant="solo"
        :modelValue="name" @blur="updateFilter('name', $event.target.value)"
        @keypress.enter="updateFilter('name', $event.target.value)" @click:clear="updateFilter('name', '')">
        <template v-slot:prepend-inner>
          <v-icon :icon="mdiMagnify" />
        </template>
      </TextField>
    </div>
    <div class="mt-2 mt-md-0 ml-md-3">
      <v-select :modelValue="category as any" bgColor="blue-grey-darken-4" :items="categories" singleLine hideDetails
        variant="solo" style="min-width:220px" @update:modelValue="updateFilter('category', $event)">
        <template v-slot:prepend-inner>
          <v-icon :icon="mdiFileQuestionOutline" />
        </template>
      </v-select>
    </div>
    <div class="mt-2 mt-md-0 ml-md-3">
      <v-select :modelValue="state as any" bgColor="blue-grey-darken-4" :items="states" singleLine hideDetails
        variant="solo" style="min-width:220px" @update:modelValue="updateFilter('state', $event)">
        <template v-slot:prepend-inner>
          <v-icon :icon="mdiProgressCheck" />
        </template>
      </v-select>
    </div>
    <div class="mt-2 mt-md-0 ml-md-3">
      <v-select :modelValue="perPage as any" bgColor="blue-grey-darken-4" :items="perPageOptions" singleLine hideDetails
        variant="solo" style="min-width:120px" @update:modelValue="updateFilter('perPage', $event)">
        <template v-slot:prepend-inner>
          <v-icon :icon="mdiViewList" />
        </template>
      </v-select>
    </div>
  </v-card>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { mdiFileQuestionOutline, mdiMagnify, mdiProgressCheck, mdiViewList } from "@mdi/js";

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
    perPage: {
      type: Number,
      required: true,
    },
    state: {
      type: String,
      required: true,
    },
    category: {
      type: String,
      required: true,
    },
  },
  emits: ["update:filter"],
  setup(props, { emit }) {
    const updateFilter = (name: string, value: any) => emit("update:filter", { name, value });

    const states = [
      { title: "All", value: "all" },
      { title: "All but uploaded", value: "all_but_uploaded" },
      { title: "Uploaded", value: "uploaded" },
      { title: "General error", value: "general_error" },
      { title: "Upload error", value: "upload_error" },
    ];

    const perPageOptions = [
      { title: "10", value: 10 },
      { title: "20", value: 20 },
      { title: "30", value: 30 },
    ];

    const categories = [
      { title: "All", value: "all" },
      { title: "Apps", value: "app" },
      { title: "Audio", value: "audio" },
      { title: "Documentaries", value: "docu" },
      { title: "EBook", value: "ebook" },
      { title: "Games", value: "game" },
      { title: "Movies", value: "movie" },
      { title: "TV", value: "tv" },
      { title: "XXX", value: "xxx" },
      { title: "Unknown", value: "unknown" },
    ]

    return {
      states,
      perPageOptions,
      categories,
      updateFilter,
      mdiMagnify,
      mdiFileQuestionOutline,
      mdiProgressCheck,
      mdiViewList,
    };
  },
});
</script>


<style lang="scss">
.browse-filters-wrapper .v-field {
  box-shadow: none !important;
}
</style>

