import { ref, computed } from "vue";
import { addEventHandler } from "@/utils/websocket";

export default () => {
  const setupStatus = ref<ISetupStatus>({} as ISetupStatus);

  const isSetupComplete = computed(() =>
    Object.values(setupStatus.value).every((value) => value)
  );

  addEventHandler(
    "SETUP_STATUS",
    ({ payload }: IResponse<ISetupStatus>) =>
      (setupStatus.value = payload)
  );

  return {
    setupStatus,
    isSetupComplete,
  };
};
