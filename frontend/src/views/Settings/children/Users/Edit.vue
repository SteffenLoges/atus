<template>
  <FormCard :loading="isLoading" title="Edit User" @submit="onSubmit">
    <v-card-text>
      <v-alert v-if="currentUserUID === uid" type="info" class="mb-4">
        After editing your own account, you will be logged
        out and redirected to the login page.
      </v-alert>

      <TextField label="Name" v-model="username" required />
      <TextField type="password" label="Password" v-model="password" hint="Leave blank to keep the current password"
        class="mb-2" persistentHint />
      <TextField type="password" label="Confirm Password" v-model="passwordConfirm" />
    </v-card-text>

    <v-card-actions class="px-5 justify-end">
      <v-btn color="error" @click.prevent="$router.push({ name: 'settings_users_manage' })">Cancel</v-btn>
      <v-btn color="primary" type="submit">Save</v-btn>
    </v-card-actions>
  </FormCard>
</template>


<script lang="ts">
import { defineComponent, ref } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";
import { useRouter, useRoute } from "vue-router";
import useUserStore from "@/store/user";

interface IEditUser {
  name: string;
}

export default defineComponent({
  async setup() {
    const globalStore = useGlobalStore();
    const user = useUserStore();
    const router = useRouter();
    const route = useRoute();

    const isLoading = ref(false);
    const uid = route.params.uid as string;
    const username = ref("");
    const password = ref("");
    const passwordConfirm = ref("");

    // --------------------------------------------------------------------------

    const r: IResponse<IEditUser> = await send("SETTINGS__USERS__GET", { uid })
    username.value = r.payload.name

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      if (password.value !== passwordConfirm.value) {
        globalStore.setError("Passwords do not match");
        return;
      }

      send("SETTINGS__USERS__UPDATE", {
        uid,
        username: username.value,
        password: password.value,
      })
        .then(() => {
          success("Settings saved successfully");
          if (user.uid === uid) {
            user.deleteAuthToken();
            router.push({ name: "login" });
          } else {
            router.push({ name: "settings_users_manage" });
          }
        })
        .catch(({ payload, statusCode }: IResponse<string>) => globalStore.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => isLoading.value = false)
    };

    // --------------------------------------------------------------------------

    return {
      uid,
      currentUserUID: user.uid,
      username,
      password,
      passwordConfirm,
      onSubmit,
      isLoading,
    };
  },
});
</script>