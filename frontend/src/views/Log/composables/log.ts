import { ref } from "vue";
import { send } from "@/utils/websocket";

export const severities = [
  { text: "DEBUG", color: "blue-grey-darken-4" },
  { text: "INFO", color: "blue-grey-darken-4" },
  { text: "WARNING", color: "yellow-darken-4" },
  { text: "ERROR", color: "red-accent-4" },
  { text: "FATAL", color: "red-darken-4" },
];

export default function useLog() {
  const entries = ref<ILogEntry[]>([]);
  const types = ref<string[]>([]);
  const count = ref(0);

  const getEntries = (params: Object) =>
    send("LOG__GET", params).then(
      ({ payload }: IResponse<ILogResponse>) => {
        entries.value = payload.entries || [];
        count.value = payload.count || 0;
        types.value = payload.types || [];
      }
    );

  return {
    entries,
    types,
    count,
    getEntries,
  };
}
