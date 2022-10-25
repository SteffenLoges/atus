<template>
  <Card title="Files" class="overflow-auto">
    <template #title-actions>
      <v-btn size="small" :icon="mdiCollapseAll" @click="($refs.folder as any).collapseAll()" />
      <v-btn size="small" :icon="mdiExpandAll" @click="($refs.folder as any).expandAll()" />
    </template>

    <v-card-text class="pt-3 pb-5 pr-6">
      <Folder v-if="files" ref="folder" :folder="files" />
    </v-card-text>
  </Card>
</template>


<script lang="ts">
import { defineComponent, ref, toRefs } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import Folder from "./components/Folder.vue";
import { mdiExpandAll, mdiCollapseAll } from "@mdi/js";

export default defineComponent({
  components: {
    Folder,
  },
  props: {
    uid: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const globalStore = useGlobalStore();
    const { uid } = toRefs(props);
    let files = ref<IFolder>();

    send("RELEASE__DETAILS__GET_FILES", { uid: uid.value, })
      .then(({ payload }: IResponse<IFolder>) => files.value = payload)
      .catch((e: IResponse<string>) => globalStore.setError(e))

    return {
      files,
      mdiExpandAll,
      mdiCollapseAll,
    };
  },
});
</script>



