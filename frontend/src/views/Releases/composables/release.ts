import { computed, ref } from "vue";
import useMetaFiles from "./metaFiles";
import {
  addEventHandler,
  removeEventHandler,
} from "@/utils/websocket";

interface IHandlerMessage<T> {
  uid: string;
  data: T;
}

export default (
  uid: string,
  initialState: IReleaseState,
  initialMetaFiles: IMetaFile[],
  initialDownloadState: IDownloadState | undefined
) => {
  const { getCoverImage } = useMetaFiles();

  const state = ref<IReleaseState>(initialState);
  const metaFiles = ref(initialMetaFiles);
  const downloadState = ref(initialDownloadState);

  const progress = computed(() => {
    if (
      ["UPLOADED", "DOWNLOADED"].includes(state.value.state)
    ) {
      return 100;
    }

    return downloadState.value?.done || 0;
  });

  const coverURL = computed(() =>
    getCoverImage(metaFiles.value)
  );

  const backgroundImage = computed(
    () => `url(${coverURL.value})`
  );

  let eventHandlers: number[] = [];
  const addEventHandlers = () => {
    // base state
    eventHandlers.push(
      addEventHandler(
        "RELEASE_DETAILS__STATE",
        ({
          payload,
        }: IResponse<IHandlerMessage<IReleaseState>>) => {
          if (payload.uid === uid) {
            state.value = payload.data;
          }
        }
      )
    );

    // download state
    eventHandlers.push(
      addEventHandler(
        "RELEASE__DETAILS__DOWNLOAD_STATE",
        ({
          payload,
        }: IResponse<IHandlerMessage<IDownloadState>>) => {
          if (payload.uid === uid) {
            downloadState.value = payload.data;
          }
        }
      )
    );

    // meta files
    eventHandlers.push(
      addEventHandler(
        "RELEASE__DETAILS__META_FILES",
        ({
          payload,
        }: IResponse<IHandlerMessage<IMetaFile[]>>) => {
          if (payload.uid === uid) {
            metaFiles.value = payload.data;
          }
        }
      )
    );
  };

  const removeEventHandlers = () =>
    eventHandlers.forEach(removeEventHandler);

  return {
    progress,
    coverURL,
    backgroundImage,
    state,
    metaFiles,
    downloadState,
    addEventHandlers,
    removeEventHandlers,
  };
};
