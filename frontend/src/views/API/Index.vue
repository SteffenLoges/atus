<template>
  <v-container class="d-flex mt-4" style="max-width: 1200px">
    <ConfirmDialog :modelValue="showResetAuthTokenDialog" @always="showResetAuthTokenDialog = false"
      @confirm="resetAuthToken()">
      Are you sure you want to reset the API token?
    </ConfirmDialog>

    <div class="flex-grow-1">
      <Card title="Main Settings">
        <v-card-text>
          <Card :loading="isLoading" title="The API Authentication Token is:" variant="text" class="card-accent">
            <v-card-text>
              <TextField :modelValue="authToken" readonly copyable single-line hide-details />
            </v-card-text>

            <v-card-actions class="justify-end">
              <v-btn color="warning" :disabled="isLoading" @click="showResetAuthTokenDialog = true">
                Regenerate Token
              </v-btn>
            </v-card-actions>
          </Card>


          <Authorisation class="mt-4" />
          <MetaFiles class="mt-4" />
          <Releases class="mt-4" :authToken="authToken" />
        </v-card-text>
      </Card>

    </div>
  </v-container>
</template>


<script lang="ts">
import { defineComponent, ref } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";
import Authorisation from "./components/Authorisation.vue";
import MetaFiles from "./components/MetaFiles.vue"
import Releases from "./components/Releases.vue";

interface IAuthToken {
  authToken: string;
}

export default defineComponent({
  components: {
    Authorisation,
    MetaFiles,
    Releases
  },
  async setup() {
    const global = useGlobalStore();

    const isLoading = ref(false);
    const showResetAuthTokenDialog = ref(false);

    const authToken = ref("");

    const r: IResponse<IAuthToken> = await send("SETTINGS__API_MANAGE__GET_TOKEN")
    authToken.value = r.payload.authToken

    const resetAuthToken = () => {
      send("SETTINGS__API_MANAGE__RESET_API_TOKEN")
        .then(({ payload }: IResponse<IAuthToken>) => {
          authToken.value = payload.authToken;
          success("A new API token has been generated.");
        })
        .catch(({ payload, statusCode }: IResponse<string>) => global.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => isLoading.value = false)
    };

    return {
      authToken,
      showResetAuthTokenDialog,
      resetAuthToken,
      isLoading,
    };
  },
});
</script>
