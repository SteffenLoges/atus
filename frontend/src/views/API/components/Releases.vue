<template>
  <Card title="Releases" variant="text" class="card-accent mb-4">
    <v-card-text>
      Allows you to get releases from the API.
      <code class="d-block my-3 py-1 select-all d-flex" data-method="GET">
        <div class="flex-grow-1 align-self-center">{{ `${baseURL}/api/releases` }}</div>
        <v-btn color="blue" variant="flat" size="small" class="text-none" target="_blank" :href="`${baseURL}/api/releases?token=${authToken}`">Example</v-btn>  
      </code>

      <Card title="Parameters" variant="tonal" color="blue-grey" class="mt-4">
        <v-card-text>
          <v-alert type="info" class="py-3 mb-4">
            All parameters are optional and case-insensitive.
          </v-alert>

          <v-table>
            <thead>
              <tr>
                <th style="width: 160px">Parameter</th>
                <th style="width: 220px">Valid Values</th>
                <th>Example</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="q of queryParams" :key="q.name">
                <td>
                  {{ q.name }}
                  <div v-if="q.default">
                    <small class="text-disabled">Default: {{ q.default }}</small>
                  </div>
                </td>
                <td>
                  <ul>
                    <li v-for="v of q.validValues" :key="v">{{ v }}</li>
                  </ul>
                </td>
                <td class="examples">
                  <ul>
                    <li v-for="e of q.examples" :key="e"><code>{{ q.name }}=<span v-html=e></span></code></li>
                  </ul>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card-text>
      </Card>
    </v-card-text>
  </Card>
</template>


<script lang="ts">
import { defineComponent } from 'vue'
import { baseURL } from '@/utils/url'

export default defineComponent({
  props: {
    authToken: {
      type: String,
      required: true,
    },
  },
  setup() {
    return {
      baseURL,
      queryParams: [
        {
          name: 'name',
          validValues: ['any text'],
          examples: [
            'The.Release.Name-Group',
            '%Part.Of.Release.Name%  <small>(must be URL encoded)</small>',
          ],
        },
        {
          name: 'category',
          validValues: [
            'Movie',
            'TV',
            'Docu',
            'App',
            'Game',
            'Audio',
            'EBook',
            'XXX',
            'Unknown',
          ],
          examples: [
            'Movies',
            'Movies,TV',
          ],
        },
        {
          name: 'state',
          validValues: [
            'new',
            'download_init',
            'downloading',
            'downloaded',
            'uploaded',
            'upload_error',
            'general_error',
          ],
          examples: [
            'downloading',
            'downloading,uploaded',
          ],
        },
        {
          name: 'order_by',
          validValues: [
            'name',
            'category',
            'state',
            'size',
            'added',
            'pre',
          ],
          default: 'added',
          examples: [
            'name',
          ],
        },
        {
          name: 'order',
          validValues: ['asc', 'desc'],
          default: 'desc',
          examples: [
            'asc',
          ],
        },
        {
          name: 'limit',
          validValues: ['1 - 100'],
          default: '25',
          examples: [
            '10',
          ],
        },
        {
          name: 'offset',
          validValues: ['any positive integer'],
          default: '0',
          examples: [
            '100',
          ],
        },
      ]
    }
  },
})
</script>


<style lang="scss" scoped>
td {
  vertical-align: top;
  padding-top: 0.5rem !important;
  padding-bottom: 0.5rem !important;

  &.examples li:not(:last-child) {
    margin-bottom: 0.5rem !important;
  }
}
</style>