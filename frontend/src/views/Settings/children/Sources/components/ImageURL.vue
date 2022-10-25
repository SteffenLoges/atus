<template>
  <v-alert v-if="imagePathAutoDetected" type="success" class="my-2">
    Good News!
    <br />
    <br />
    {{ appName }} automatically detected a valid image path. You can proceed with the next step.
  </v-alert>

  <Switch v-if="imagePathAutoDetected" v-model="overrideAutoDetectedImagePath" @change="onAutoDetectChange()"
    label="Override Auto Detected Settings" />

  <SlideDownTransition>
    <div v-if="overrideAutoDetectedImagePath ||      !imagePathAutoDetected">
      <v-alert v-if="overrideAutoDetectedImagePath" type="warning">
        {{ appName }} already detected a valid image using the <code>{{ orgImagePath }}</code> path.
        <p class="mt-2">You shouln't override this path unless you know what you're doing.</p>
      </v-alert>

      <Switch :modelValue="imagePathUseAsKey" @update:modelValue="$emit('update:imagePathUseAsKey', $event)" :disabled="imagePathAutoDetected &&
      !overrideAutoDetectedImagePath" :true-value="false" :false-value="true" @change="onAutoDetectChange()"
        label="Use a URL instead of a path" persistent-hint />

      <v-alert type="info" class="mb-3">
        <SlideTransition>
          <div v-if="imagePathUseAsKey">
            <div class="text-subtitle-2 mb-2">Path Info:</div>
            Paths are XML-Tags, relative to the item, concatenated with an underscore.
            <br />
            <br />
            <small><i>Examples:</i><br />
              Path for
              <code>&lt;item&gt;&lt;cover&gt;http://tracker.to/bitbucket/f-12345-1.jpg&lt;/cover&gt;&lt;/item&gt;</code>
              is <code>cover</code>
              <br />
              Path for
              <code>&lt;item&gt;&lt;links&gt;&lt;thumbnail&gt;http://tracker.to/bitbucket/f-12345-1.jpg&lt;/thumbnail&gt;&lt;/links&gt;&lt;/item&gt;</code>
              is <code>links_thumbnail</code>
            </small>

            <div v-if="!imagePathAutoDetected && isSetup" class="mt-4">
              <strong class="text-red-lighten-2">
                {{ appName }} tried to auto detect the path but wasn't able to find a valid image.
                <br />This option will most likely <u>not work</u>, you should use a URL instead.
              </strong>
            </div>
          </div>
          <div v-else>
            <div class="text-subtitle-2 mb-2">URL Info:</div>
            Specify the URL to an image.<br /><br />
            Available variables are: <code>{id}</code> and <code>{title}</code>. <br /><br />
            <small><i>Examples:</i><br />
              If the URL is <code>http://tracker.to/bitbucket/f-12345-1.jpg</code>
              use <code>http://tracker.to/bitbucket/f-{id}-1.jpg</code>
              <br />
              If the URL is <code>http://tracker.to/covers/1234/Torrent.Name.jpg</code>
              use <code>http://tracker.to/covers/{id}/{name}.jpg</code><br />
            </small>
          </div>
        </SlideTransition>
      </v-alert>

      <TextField :modelValue="imagePath" @update:modelValue="$emit('update:imagePath', $event)" persistent-hint
        :disabled="imagePathAutoDetected &&
        !overrideAutoDetectedImagePath" :label="imagePathUseAsKey ? 'Path' : 'URL'" persistent-placeholder
        :placeholder="`e.g. ${imagePathUseAsKey ? 'link': 'https://awesome-tracker.to/bitbucket/f-{id}-1.jpg'}`"
        hint="Leave blank to ignore images from this source">
      </TextField>
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
    imagePath: {
      type: String,
      default: "",
    },
    imagePathAutoDetected: {
      type: Boolean,
      default: false,
    },
    imagePathUseAsKey: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:imagePath", "update:imagePathUseAsKey"],
  setup(props, { emit }) {
    const overrideAutoDetectedImagePath = ref(false);
    const { imagePath, imagePathUseAsKey, imagePathAutoDetected } = toRefs(props);
    const orgImagePath = imagePath.value;

    const onAutoDetectChange = () => {
      if (imagePathUseAsKey.value) {
        emit("update:imagePath", orgImagePath);
        return;
      }

      if (imagePathAutoDetected.value && !overrideAutoDetectedImagePath.value) {
        imagePathUseAsKey.value = true;
      }

      emit("update:imagePath", "");
    };

    return {
      appName: import.meta.env.VITE_APP_NAME,
      overrideAutoDetectedImagePath,
      orgImagePath,
      onAutoDetectChange,
    };
  },
});
</script>
