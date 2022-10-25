<template>
  <v-alert v-if="metaPathAutoDetected" type="success" class="my-2">
    Good News!
    <p class="mt-2">{{ appName }} automatically detected a valid .torrent path. You can proceed with the next step.</p>
  </v-alert>

  <Switch v-if="metaPathAutoDetected" v-model="overrideAutoDetectedMetaPath" @change="onAutoDetectChange()"
    label="Override Auto Detected Settings" />

  <SlideDownTransition>
    <div v-if="overrideAutoDetectedMetaPath || !metaPathAutoDetected">
      <v-alert v-if="overrideAutoDetectedMetaPath" type="warning">
        {{ appName }} already detected a valid .torrent file using the <code>{{ orgMetaPath }}</code> path.
        <p class="mt-2">You shouln't override this path unless you know what you're doing.</p>
      </v-alert>

      <Switch :modelValue="metaPathUseAsKey" @update:modelValue="$emit('update:metaPathUseAsKey', $event)"
        :disabled="metaPathAutoDetected && !overrideAutoDetectedMetaPath" :true-value="false" :false-value="true"
        @change="onAutoDetectChange()" label="Use a URL instead of a path" persistent-hint />

      <v-alert type="info" class="mb-3">
        <SlideTransition>
          <div v-if="metaPathUseAsKey">
            <div class="text-subtitle-2 mb-2">Path Info:</div>
            Paths are XML-Tags, relative to the item, concatenated with an underscore.
            <br />
            <br />
            <small><i>Examples:</i><br />
              Path for
              <code>&lt;item&gt;&lt;link&gt;https://tracker.to/download?id=1234&lt;/link&gt;&lt;/item&gt;</code>
              is <code>link</code>
              <br />
              Path for
              <code>&lt;item&gt;&lt;links&gt;&lt;torrent&gt;https://tracker.to/download?id=1234&lt;/torrent&gt;&lt;/links&gt;&lt;/item&gt;</code>
              is <code>links_torrent</code>
            </small>

            <div v-if="!metaPathAutoDetected && isSetup" class="mt-4">
              <strong class="text-red-lighten-2">
                {{ appName }} tried to auto detect the path but wasn't able to find a valid .torrent file.<br />
                This option will most likely <u>not work</u>, you should use a URL instead.
              </strong>
            </div>
          </div>
          <div v-else>
            <div class="text-subtitle-2 mb-2">URL Info:</div>
            Specify the URL to a .torrent file.<br /><br />
            Available variables are: <code>{id}</code> and
            <code>{title}</code>. <br /><br />

            <small><i>Examples:</i><br />
              If the URL is <code>https://tracker.to/download.php?torrent=1234</code>
              use <code>https://tracker.to/download.php?torrent={id}</code>
              <br />
              If the URL is <code>https://tracker.to/download/1234/Torrent.Name.torrent</code>
              use <code>https://tracker.to/download/{id}/{title}.torrent</code><br />

              If the URL is <code>https://tracker.to/download?id=1234&amp;passkey=a1b2d3</code>
              use <code>https://tracker.to/download?id={id}&amp;passkey=a1b2d3</code>
            </small>
          </div>
        </SlideTransition>
      </v-alert>

      <TextField :modelValue="metaPath" @update:modelValue="$emit('update:metaPath', $event)" persistent-hint
        :disabled="metaPathAutoDetected && !overrideAutoDetectedMetaPath" :label="metaPathUseAsKey ? 'Path' : 'URL'"
        persistent-placeholder
        :placeholder="`e.g. ${metaPathUseAsKey ? 'link' : 'https://tracker.to/download.php?torrent={id}'}`" hide-details
        required />
    </div>
  </SlideDownTransition>
</template>


<script lang="ts">
import { defineComponent, ref, toRefs } from "vue";

export default defineComponent({
  props: {
    isSetup: {
      type: Boolean,
      default: false,
    },
    metaPath: {
      type: String,
      default: "",
    },
    metaPathUseAsKey: {
      type: Boolean,
      default: false,
    },
    metaPathAutoDetected: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:metaPath", "update:metaPathUseAsKey"],
  setup(props, { emit }) {
    const overrideAutoDetectedMetaPath = ref(false);
    const { metaPathUseAsKey, metaPathAutoDetected, metaPath } = toRefs(props);
    const orgMetaPath = metaPath.value;

    const onAutoDetectChange = () => {
      if (metaPathUseAsKey.value) {
        emit("update:metaPath", orgMetaPath);
        return;
      }

      if (metaPathAutoDetected.value && !overrideAutoDetectedMetaPath.value) {
        metaPathUseAsKey.value = true;
      }

      emit("update:metaPath", metaPath.value);
    };

    return {
      appName: import.meta.env.VITE_APP_NAME,
      onAutoDetectChange,
      overrideAutoDetectedMetaPath,
      orgMetaPath,
    };
  },
});
</script>
