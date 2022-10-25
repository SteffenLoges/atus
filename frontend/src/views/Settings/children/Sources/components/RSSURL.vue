<template>
  <TextField :modelValue="rssURL" @update:modelValue="$emit('update:rssURL', $event)" persistent-hint
    label="RSS-Feed-URL" persistent-placeholder placeholder="e.g. https://awseome-tracker.to/rss.php?passkey=123456789"
    hint="Enter the RSS-Feed-URL of the tracker you want to add" required />

  <Switch :modelValue="requiresCookies" @update:modelValue="$emit('update:requiresCookies', $event)" hide-details
    label="This tracker requires cookies" />

  <VSlideYTransition>
    <v-table density="compact" color="transparent" v-show="requiresCookies">
      <tbody>
        <tr v-for="(cookie, i) in cookiesComputed" :key="i">
          <td class="pl-0 py-2" style="width: 250px">
            <TextField density="compact" variant="outlined" v-model.lazy="cookie.name" placeholder="Name"
              hide-details />
          </td>
          <td class="px-0 py-2">
            <TextField density="compact" variant="outlined" v-model.lazy="cookie.value" placeholder="Value"
              hide-details />
          </td>
          <td style="width: 60px">
            <v-btn size="small" :icon="mdiDelete" flat @click="deleteCookie(i)" />
          </td>
        </tr>
        <tr>
          <td class="pa-0" colspan="3">
            <v-btn variant="text" block :prepend-icon="mdiPlus" @click="addCookie()">Add more</v-btn>
          </td>
        </tr>
      </tbody>
    </v-table>
  </VSlideYTransition>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs, computed } from "vue";
import { mdiDelete, mdiPlus } from "@mdi/js";

export default defineComponent({
  props: {
    isSetup: {
      type: Boolean,
      default: false,
    },
    rssURL: {
      type: String,
      required: true,
    },
    requiresCookies: {
      type: Boolean,
      required: true,
    },
    cookies: {
      type: Array as PropType<ICookie[]>,
      required: true,
    },
  },
  emits: [
    "update:rssURL",
    "update:requiresCookies",
    "update:cookies",
  ],
  setup(props, { emit }) {
    const { cookies } = toRefs(props);

    const cookiesComputed = computed({
      get: () => cookies.value,
      set: (v) => emit("update:cookies", v),
    });

    const addCookie = () => (cookiesComputed.value = [...cookiesComputed.value, { name: "", value: "" },]);
    const deleteCookie = (i: number) => cookiesComputed.value.splice(i, 1);

    return {
      cookies,
      cookiesComputed,
      addCookie,
      deleteCookie,
      mdiDelete,
      mdiPlus,
    };
  },
});
</script>
