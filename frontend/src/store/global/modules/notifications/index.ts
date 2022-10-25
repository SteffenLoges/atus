import { ref } from "vue";
import { info } from "@/plugins/toast";
import { addEventHandler } from "@/utils/websocket";

export default () => {
  // ToDo: Tracking of new releases is currently handled by checking the type of a notification
  // This is fine for now, but should be done in a more robust way in the future
  const sumNewReleases = ref(0);
  const resetSumNewReleases = () =>
    (sumNewReleases.value = 0);

  addEventHandler(
    "NOTIFICATION",
    ({ payload }: IResponse<INotification>) => {
      if (payload.type === "NEW_RELEASE") {
        sumNewReleases.value++;
      }

      info(payload.title, payload.message);
    }
  );

  return {
    sumNewReleases,
    resetSumNewReleases,
  };
};
