<template>
  <div class="file text-truncate">
    <v-icon :icon="fileIcon" color="teal-darken-2" />
    {{ file.name }}
    <small class="ml-2 text-medium-emphasis" v-text="fileSize"></small>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, toRefs } from "vue";
import { bytesHumanReadable } from "@/utils/conversion";
import {
  mdiFileCheck, mdiImage, mdiFileMusic, mdiFileDocument, mdiZipBox,
  mdiFilePdfBox, mdiFileEye, mdiFileVideo, mdiFileXmlBox, mdiBookInformationVariant, mdiFileQuestion
} from "@mdi/js";

export default defineComponent({
  props: {
    file: {
      type: Object as PropType<IFile>,
      required: true,
    },
  },
  setup(props) {
    const { file } = toRefs(props);

    const fileSize = bytesHumanReadable(file.value.length);

    const fileIcon = (() => {
      const fileExt = file.value.name.split(".").pop();

      switch (fileExt) {
        case "png":
        case "jpg":
        case "jpeg":
          return mdiImage;
        case "mp4":
        case "avi":
        case "mkv":
          return mdiFileVideo;
        case "sfv":
          return mdiFileCheck;
        case "mp3":
        case "flac":
          return mdiFileMusic;
        case "pdf":
          return mdiFilePdfBox;
        case "doc":
        case "docx":
          return mdiFileDocument;
        case "zip":
        case "rar":
          return mdiZipBox;
        case "srt":
          return mdiFileEye;
        case "xml":
          return mdiFileXmlBox;
        case "nfo":
          return mdiBookInformationVariant;
      }

      if (fileExt && /^r[0-9]{2}$/.test(fileExt)) {
        return mdiZipBox;
      }

      return mdiFileQuestion;
    })();

    return {
      file,
      fileIcon,
      fileSize,
    };
  },
});
</script>


<style lang="scss" scoped>
.file {
  $file-height: 23px;

  height: $file-height;
  line-height: $file-height;
  position: relative;
  margin-left: 45px;
}
</style>