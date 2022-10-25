<template>
  <video ref="player" :id="id" class="video-js">
    <source :src="src" :type="type" />
    <slot>
      <p class="vjs-no-js">
        To view this video please enable JavaScript, and consider upgrading to a web browser that
        <a href="https://videojs.com/html5-video-support/" target="_blank">supports HTML5 video</a>
      </p>
    </slot>
  </video>
</template>


<script lang="ts">
import { defineComponent, onMounted, onBeforeUnmount } from "vue";
import useUserStore from "@/store/user";
import videojs from "video.js";
import "video.js/dist/video-js.css";

export default defineComponent({
  components: {},
  props: {
    src: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: "application/dash+xml",
    },
    id: {
      type: String,
      default: "player",
    },
  },
  setup(props) {
    const { getAuthToken } = useUserStore();

    let player: videojs.Player;

    onMounted(() => {
      // @ts-ignore
      videojs.Vhs.xhr.beforeRequest = (options) => {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${getAuthToken()}`,
        };
      };

      var options: videojs.PlayerOptions = {
        fluid: true,
        preload: "auto",
        controls: true,
        controlBar: {
          pictureInPictureToggle: false,
        },
      };

      player = videojs(props.id, options);
    });

    onBeforeUnmount(() => {
      if (player) {
        player.dispose();
      }
    });
  },
});
</script>




