<template>
  <v-container fluid class="d-flex align-center justify-center mt-4">
    <Card class="px-4 py-3" :title="`Thank you for using ${appName}!`">
      <v-card-text>
        <p>
          This software is and will always be open-source and free to use.<br>
          However, if you want to support the development of this project, you can do so by donating.
        </p>

        <v-row class="mt-4">
          <v-col cols="12" md="6">
            <v-btn class="text-none" size="large" block color="purple-darken-3"
              :href="dereferURL('https://www.paypal.com/donate/?hosted_button_id=YY94RAW5PWQEG', true)" target="_blank">
              <v-icon class="mr-1" left :icon="mdiAlphaPCircle" color="green-lighten-3" style="margin-top: 2px;" />
              Donate with PayPal
            </v-btn>
          </v-col>
          <v-col cols="12" md="6">
            <v-btn class="text-none" size="large" block color="purple-darken-3"
              :href="dereferURL('https://commerce.coinbase.com/checkout/fe39872c-d48d-4a89-8827-f9ada803c158', true)"
              target="_blank">
              <v-icon class="mr-1" left :icon="mdiBitcoin" color="green-lighten-3"
                style="margin-top: 2px;transform: rotate(-15deg)" />
              Donate with Crypto
            </v-btn>
          </v-col>
        </v-row>

        <v-alert type="info" class="mt-8 py-3">
          Please note that a donation is solely a way to support the development of this project.<br>
          It does not grant you any special privileges, nor does it guarantee a faster
          response to GitHub issues.
        </v-alert>
      </v-card-text>
    </Card>
  </v-container>


</template>


<script lang="ts">
import { defineComponent } from "vue";
// import { useHead } from "@vueuse/head";
import { dereferURL } from "@/utils/url";
import { mdiAlphaPCircle, mdiBitcoin } from "@mdi/js";

export default defineComponent({
  setup() {
    // useHead({
    //   script: [
    //     {
    //       src: "https://commerce.coinbase.com/v1/checkout.js?version=201807",
    //       defer: true,
    //     },
    //   ],
    // });

    return {
      appName: import.meta.env.VITE_APP_NAME,
      dereferURL,
      mdiAlphaPCircle,
      mdiBitcoin,
    }
  },
});
</script>