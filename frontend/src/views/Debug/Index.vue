<template>
  <v-container class="d-flex mt-4" style="max-width: 1200px">
    <Card color="black" class="w-100">
      <div class="d-flex flex-column-reverse" style="height: clamp(250px, 50vh, 600px);overflow-y: scroll;">
        <div>
          <Entry v-for="(e, i) of entries" :key="i" v-bind="e" />
        </div>
      </div>
    </Card>
  </v-container>
</template> 


<script lang="ts">
import { defineComponent, ref, onBeforeUnmount } from "vue";
import { addEventHandler, removeEventHandler, send } from "@/utils/websocket";
import Entry from "./components/Entry.vue";

export default defineComponent({
  components: {
    Entry,
  },
  async setup() {
    let entries = ref<any>([] as any);

    let messageHandler = -1
    onBeforeUnmount(() => removeEventHandler(messageHandler))

    const { payload }: IResponse<IRelease> = await send("DEBUG__GET_CACHE");
    entries.value = payload;

    // it is theoretically possible that we are missing some messages between cache and messagehandler
    // ToDo: add message handler before cache call and merge the two arrays

    // ToDo: pause updates when user scrolls up

    messageHandler = addEventHandler("DEBUG_ENTRY", ({ payload }) => {
      if (entries.value.length >= 250) {
        entries.value.shift();
      }

      entries.value.push(payload);
    })

    return {
      entries,
    };
  },
});
</script>

