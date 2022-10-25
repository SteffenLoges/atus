<template>
  <ConfirmDialog :modelValue="showDeleteConfirmDialog" @cancel="showDeleteConfirmDialog = false"
    @confirm="onDeleteConfirm()">
    Are you sure you want to delete this release?
    <p class="mt-3">Files on fileservers will NOT be deleted.</p>
  </ConfirmDialog>

  <ConfirmDialog :modelValue="showUploadConfirmDialog" @cancel="showUploadConfirmDialog = false"
    @confirm="onUploadConfirm()">
    Are you sure you want to upload this release?
  </ConfirmDialog>

  <ErrorDialog :modelValue="uploadErrorMessage !== null" @dismiss="uploadErrorMessage = null" title="Upload failed">
    <v-alert type="error" variant="text" density="compact" class="mb-3">
      The file could not be uploaded.<br />
      Check your settings and try again.
    </v-alert>

    <small>The Server responded with the following message:</small>
    <v-alert type="info" variant="tonal" border="start" density="compact" :icon="false">
      <span class="text-grey" style="white-space: pre-wrap">{{ uploadErrorMessage }}</span>
    </v-alert>
  </ErrorDialog>


  <template v-if="release">
    <Header :release="release" :state="state" :downloadState="downloadState" :metaFiles="metaFiles"
      :uploadInProgress="uploadInProgress" @delete="showDeleteConfirmDialog = true"
      @upload="showUploadConfirmDialog = true" />

    <Sample v-if="sampleVideoMetaFiles.length > 0 && sampleVideoMetaFiles[0].state !== 'ERROR'"
      :metaFiles="sampleVideoMetaFiles" />

    <Images v-if="imageMetaFiles.length > 0" :metaFiles="imageMetaFiles" />

    <section class="py-3 py-lg-8">
      <v-container fluid>
        <v-row class="justify-space-evenly">
          <v-col cols="12" lg="6" class="d-flex flex-grow-1 order-lg-1" v-if="nfoMetaFiles.length > 0"
            style="max-width: 800px">
            <NFOContainer class="h-100 w-100" :metaFiles="nfoMetaFiles" />
          </v-col>

          <v-col cols="12" lg="6" class="d-flex flex-grow-1" style="max-width: 800px">
            <Files class="h-100 w-100" :uid="release.uid" />
          </v-col>
        </v-row>
      </v-container>
    </section>

    <section class="pt-8 pb-4">
      <v-container fluid>
        <Log :uid="release.uid" />
      </v-container>
    </section>
  </template>
</template>


<script lang="ts">
import { defineComponent, defineAsyncComponent, ref, computed, onBeforeUnmount } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useHead } from "@vueuse/head"
import { send } from "@/utils/websocket"
import useGlobalStore from "@/store/global";
import useRelease from "../composables/release";
import { IMAGE_TYPES } from "../composables/metaFiles";
import { success } from "@/plugins/toast";
import Header from "./components/Header.vue";
import Files from "../components/Files/Index.vue";
import Log from "./components/Log.vue";
const Sample = defineAsyncComponent(() => import("./components/Sample.vue"));
const Images = defineAsyncComponent(() => import("./components/Images.vue"));
const NFOContainer = defineAsyncComponent(() => import("./components/NFOContainer.vue"));

export default defineComponent({
  components: {
    Header,
    NFOContainer,
    Images,
    Files,
    Sample,
    Log,
  },
  async setup() {
    const router = useRouter();
    const route = useRoute();
    const globalStore = useGlobalStore();
    const title = ref("");
    useHead({ title })

    let release = ref<IRelease>();
    const uid = route.params.uid.toString();

    let _removeMessageHandlers = () => { };
    onBeforeUnmount(() => _removeMessageHandlers());

    const { payload }: IResponse<IRelease> = await send("RELEASE__DETAILS__GET", { uid })
    release.value = payload;
    title.value = release.value.name;

    const { state, downloadState, coverURL, metaFiles, addEventHandlers, removeEventHandlers } =
      useRelease(uid, release.value.state, release.value.metaFiles, release.value.downloadState)

    addEventHandlers();
    _removeMessageHandlers = removeEventHandlers;

    const nfoMetaFiles = computed(() => metaFiles.value.filter(({ type }) => type === "NFO"))
    const imageMetaFiles = computed(() => metaFiles.value.filter(({ type }) => IMAGE_TYPES.includes(type)))
    const sampleVideoMetaFiles = computed(() => metaFiles.value.filter(({ type }) => type === "SAMPLE_VIDEO"))

    const showDeleteConfirmDialog = ref(false);
    const onDeleteConfirm = () => send("RELEASE__DELETE", { uid })
      .then(() => {
        success("Release deleted");
        router.push({ name: "releases_browse" });
      })
      .catch(({ payload }: IResponse<string>) => globalStore.setError(payload))
      .finally(() => showDeleteConfirmDialog.value = false);

    const showUploadConfirmDialog = ref(false);
    const uploadErrorMessage = ref<string | null>(null)
    const uploadInProgress = ref(false);
    const onUploadConfirm = () => {
      uploadInProgress.value = true;
      showUploadConfirmDialog.value = false;

      send("RELEASE__UPLOAD", { uid })
        .then(() => success("Release uploaded successfully"))
        .catch(({ payload }: IResponse<string>) => uploadErrorMessage.value = payload)
        .finally(() => uploadInProgress.value = false);
    }

    return {
      coverURL,
      release,
      metaFiles,
      state,
      downloadState,
      nfoMetaFiles,
      imageMetaFiles,
      sampleVideoMetaFiles,
      showDeleteConfirmDialog,
      onDeleteConfirm,
      showUploadConfirmDialog,
      onUploadConfirm,
      uploadErrorMessage,
      uploadInProgress
    };
  },
});
</script>


<style lang="scss" scoped>
section:nth-of-type(2n) {
  background: #0c0e17;
}
</style>