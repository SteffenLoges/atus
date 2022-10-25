<template>
  <Card title="Meta Files" variant="text" class="card-accent mb-4">
    <v-card-text>
      You can access meta files by using the following URL:
      <code class="d-block my-3 select-all" data-method="GET">{{ `${baseURL}/api/data/{RELEASE_UID}/{FILE}` }}</code>

      <v-alert type="info" class="mt-4 py-2">
        For security reasons, access to the original .torrent file is restricted to internal requests only.
      </v-alert>
    </v-card-text>
  </Card>
</template>


<script lang="ts">
import { defineComponent } from 'vue'
import { baseURL } from '@/utils/url'

export default defineComponent({
  setup() {
    return {
      baseURL,
    }
  },
})
</script>