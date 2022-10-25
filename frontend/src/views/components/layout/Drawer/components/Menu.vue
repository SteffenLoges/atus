<template>
  <v-list class="main-nav pa-0" nav variant="flat" density="compact" openStrategy="single" bgColor="transparent"
    v-model:opened="open">
    <template v-for="item of items" :key="item.text">
      <v-divider v-if="item.divider" class="mt-2 mb-2" />

      <v-list-group v-else-if="item.children" :value="item.title" :expandIcon="mdiChevronDown"
        :collapseIcon="mdiChevronUp">
        <template v-slot:activator="{ props }">
          <v-list-item v-bind="props" :prepend-icon="item.icon" :title="item.title" :disabled="item.disabled" />
        </template>
        <v-list-item v-for="child of item.children" :key="child.title" :title="child.title" :to="child.to" exact
          :disabled="child.disabled" />
      </v-list-group>

      <v-list-item v-else inact :to="item.to" :prepend-icon="item.icon" :class="item.class" :title="item.title"
        :disabled="item.disabled" />
    </template>
  </v-list>
</template>

<script lang="ts">
import { defineComponent, computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { RouteLocationRaw } from "vue-router";
import useGlobalStore from "@/store/global";
import {
  mdiBookOpenPageVariant, mdiCloudUpload, mdiCodeTags, mdiFolderStar, mdiStar,
  mdiVideo, mdiBug, mdiFilter, mdiSourceBranch, mdiServer, mdiAccount, mdiChevronDown,
  mdiChevronUp
} from "@mdi/js";

interface IMenuItem {
  title?: string;
  icon?: string;
  to?: RouteLocationRaw;
  class?: string;
  children?: IMenuItem[];
  disabled?: boolean;
  divider?: true;
}


export default defineComponent({
  setup() {
    const globalStore = useGlobalStore();
    const { isSetupComplete } = storeToRefs(globalStore);

    const open = ref<any>([]);

    let items = computed(() => {

      let ret: IMenuItem[] = []

      if (!isSetupComplete.value) {
        ret.push({
          title: "Getting Started",
          icon: mdiStar,
          class: "text-yellow-accent-3",
          to: {
            name: "setup"
          },
        })

        ret.push({ divider: true })
      }

      ret.push({
        title: "Releases",
        icon: mdiFolderStar,
        // disabled: !isSetupComplete.value,
        to: {
          name: "releases_browse",
        },
      })


      ret.push({ divider: true })

      ret.push({
        title: "Sources",
        icon: mdiSourceBranch,
        children: [
          {
            title: "Manage",
            to: {
              name: "settings_sources_manage",
            },
          },
          {
            title: "Add New",
            to: {
              name: "settings_sources_add",
            },
          },
        ],
      },
        {
          title: "Fileservers",
          icon: mdiServer,
          children: [
            {
              title: "Manage",
              to: {
                name: "settings_fileservers_manage",
              },
            },
            {
              title: "Add New",
              to: {
                name: "settings_fileservers_add",
              },
            },
            {
              title: "Global Settings",
              to: {
                name: "settings_fileservers_settings",
              },
            },
          ],
        },

        {
          title: "Filters",
          icon: mdiFilter,
          children: [
            {
              title: "Catgories",
              to: {
                name: "settings_filters_categories",
              },
            },
            {
              title: "Miscellaneous",
              to: {
                name: "settings_filters_misc",
              },
            },
          ],
        },
        {
          title: "Upload Settings",
          icon: mdiCloudUpload,
          to: {
            name: "settings_destinations",
          },
        },
        {
          title: "Samples",
          icon: mdiVideo,
          to: {
            name: "settings_samples",
          },
        },
        {
          title: "Users",
          icon: mdiAccount,
          to: {
            name: "settings_users_manage",
          },
        },
        {
          title: "Log",
          icon: mdiBookOpenPageVariant,
          to: {
            name: "log",
          },
        },
        {
          title: "Debug",
          icon: mdiBug,
          to: {
            name: "debug",
          },
        },
        {
          title: "API",
          icon: mdiCodeTags,
          to: {
            name: "api",
          },
        })

      return ret
    });

    return {
      items,
      open,
      mdiFolderStar,
      mdiChevronDown,
      mdiChevronUp,
    };
  },
});
</script>


<style lang="scss">
.main-nav {
  overflow: hidden;

  .v-list-item,
  .v-list-group__items {
    border-radius: 0 !important;
  }

  >.v-list-item,
  >.v-list-group>.v-list-item {
    &:not(&--active) {
      background: transparent;
    }

    padding-left: 13px;

    .v-icon {
      margin-left: 3px !important;
      margin-right: 21px !important;
    }
  }

  >.v-list-group .v-list-group__items .v-list-item {
    padding-left: 61px !important;
  }

  .v-list-item .v-list-item__append .v-icon {
    margin-right: 5px !important;
  }
}
</style>