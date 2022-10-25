import { ref } from "vue";
import { addEventHandler } from "@/utils/websocket";

export default () => {
  let fileserverStatistics = ref<IFileserver[]>([]);

  addEventHandler(
    "FILESERVER_STATISTICS",
    ({ payload }: IResponse<IFileserver[]>) =>
      (fileserverStatistics.value = payload)
  );

  return {
    fileserverStatistics,
  };
};
