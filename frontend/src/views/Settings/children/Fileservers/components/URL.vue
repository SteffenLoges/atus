<template>
  <v-alert type="info" class="mb-4" v-if="isSetup">
    <div class="mb-3">Make sure the
      <a target="blank" :href="dereferURL('https://github.com/SteffenLoges/atus-rutorrent-api')">{{ appName }} ruTorrent
        API</a>
      plugin is installed and enabled on your fileserver.
    </div>

    <small>
      <p class="mt-2">
        <b>Important:</b> If your fileserver requires basic authentication, add the username and password to the
        URL.<br /><i>Example:</i>
        <code class="ml-1">https://<b>myuser</b>:<b>mypassword</b>@fileserver.to/rutorrent</code>
      </p>
    </small>
  </v-alert>

  <TextField :modelValue="url" @update:modelValue="$emit('update:url', $event)" persistent-hint label="URL"
    persistent-placeholder placeholder="e.g. https://fileserver.to/rutorrent"
    :hint="`Enter the URL of the fileserver${isSetup ? ' you want to add' : ''}`" required />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { dereferURL } from "@/utils/url";

export default defineComponent({
  props: {
    isSetup: {
      type: Boolean,
      default: false,
    },
    url: {
      type: String,
      default: "",
    },
  },
  emits: ["update:url"],
  setup() {
    return {
      appName: import.meta.env.VITE_APP_NAME,
      dereferURL,
    };
  },
});
</script>
