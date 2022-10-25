<template>
  <div class="folder text-no-wrap" :class="{ 'ml-5': depth > 0 }">
    <h4 @click="toggle" class="title text-truncate" :class="{ expanded }">
      <v-icon class="chevron" :icon="mdiChevronRight" />
      <v-icon color="yellow-darken-4" :icon="expanded ? mdiFolderOpen : mdiFolderPlus" />
      {{ folder.path }}
    </h4>
    <VSlideYTransition>
      <div v-show="expanded">
        <div v-if="folder.subFolders">
          <Folder ref="subFolders" v-for="f in folder.subFolders" :depth="depth + 1" :key="f.path" :folder="f" />
        </div>

        <div v-if="folder.files">
          <File v-for="file of folder.files" :key="file.name" :file="file" />
        </div>
      </div>
    </VSlideYTransition>
  </div>
</template>



<script lang="ts">
import { defineComponent, PropType, toRefs, ref } from "vue";
import { mdiFolderPlus, mdiFolderOpen, mdiChevronRight } from "@mdi/js";
import File from "./File.vue";

export default defineComponent({
  name: "Folder",
  components: {
    File,
  },
  props: {
    folder: {
      type: Object as PropType<IFolder>,
      required: true,
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const { folder } = toRefs(props);
    const expanded = ref(props.depth <= 2);
    const subFolders = ref<any>(null);

    let toggle = () => {
      if (expanded.value) {
        collapseAll();
        return;
      }

      expanded.value = true;
    };

    let collapseAll = () => {
      expanded.value = false;

      if (subFolders.value) {
        for (let subFolder of subFolders.value) {
          subFolder.collapseAll();
        }
      }
    };

    let expandAll = () => {
      expanded.value = true;

      if (subFolders.value) {
        for (let subFolder of subFolders.value) {
          subFolder.expandAll();
        }
      }
    };

    return {
      toggle,
      collapseAll,
      expandAll,
      folder,
      subFolders,
      expanded,
      mdiFolderPlus,
      mdiFolderOpen,
      mdiChevronRight
    };
  },
});
</script>

<style lang="scss" scoped>
.folder {
  margin-bottom: 2px;

  .title {
    cursor: pointer;

    &:not(.expanded) {
      opacity: 0.75;
    }

    .chevron {
      width: 25px;
      transition: transform 0.2s ease-in-out;
    }

    &.expanded {
      .chevron {
        transform: rotate(90deg);
      }
    }
  }
}
</style>