<template>
  <pre class="nfo" v-html="nfoComputed"></pre>
</template>
 

<script lang="ts">
import { defineComponent, toRefs, computed } from "vue";
import { dereferURL } from "@/utils/url";

export default defineComponent({
  props: {
    data: {
      type: String,
      required: true,
    },
    fontSize: {
      type: Number,
      default: 12,
    },
  },
  emits: ["update:fontSize"],
  setup(props, { emit }) {
    const { fontSize, data } = toRefs(props);

    // ToDo: find a better way to handle nfos
    // the current solution generates a lot of unnecessary markup
    let longestRow = 0;
    const nfoComputed = computed(() => {

      longestRow = 0;
      let n = data.value;

      if (!n) {
        return "";
      }

      // get longest row length
      let r = n.split("\n").map((l: string) => l.trimEnd());
      for (let row of r) {
        if (row.length > longestRow) {
          longestRow = row.length;
        }
      }

      // set initial font size
      // this should work as long as we are using fixed width fonts
      emit("update:fontSize", Math.max(Math.min(21 - Math.log(longestRow * 100), 20), 8));

      // equalize row lengths so we can center the nfo
      for (let i in r) {
        if (r[i].length < longestRow) {
          r[i] += " ".repeat(longestRow - r[i].length);
        }
      }

      // remove empty lines at the beginning
      for (let i = 0; i < r.length; i++) {
        if (r[i].trim() !== "") {
          break;
        }

        r.splice(i, 1);
      }

      // remove empty lines at the end
      for (let i = r.length - 1; i >= 0; i--) {
        if (r[i].trim() !== "") {
          break;
        }

        r.splice(i, 1);
      }

      n = r.join("\n");

      // find urls and store them in a map for later
      let urlMap: string[] = [];
      n = n.replace(/(https?:\/\/|www.)[^()<>"'\s]+/g, (match) => {
        urlMap.push(match);
        return `mappedURL${urlMap.length - 1}`;
      });

      // highlight text
      n = n.replace(/[a-z0-9 ]{1,}[:@;\-/,.&!)()]?/gim, '<span class="alpha">$&</span>');

      // remove useless tags
      n = n.replace(/<span class="alpha">(\s+)<\/span>/g, "$1");

      // make links clickable
      n = n.replace(/mappedURL([0-9]+)/g,
        (match, index) => `<a href="${dereferURL(urlMap[index])}" target="_blank">${urlMap[index]}</a>`
      );



      return n
    });

    return {
      nfoComputed,
      fontSize,
    };
  },
});
</script>




<style lang="scss">
pre.nfo {
  $font-color: #eee;
  $ascii-color: #888;

  // transition: font-size 0.3s ease-in-out;
  font-family: monospace, Courier, "Lucida Console" !important;
  font-size: calc(v-bind(fontSize) * 1px);
  line-height: calc(v-bind(fontSize) * 1px);
  margin: auto;
  color: $ascii-color;
  text-shadow: -1px -1px 0px $ascii-color;
  text-align: center;
  user-select: none;

  span.alpha {
    color: $font-color;
    text-shadow: none;
    user-select: text;
  }

  a {
    color: skyblue;
    text-decoration: none;
  }
}
</style>