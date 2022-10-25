<template>
  <div style="position: relative">
    <FormCard :loading="isLoading" title="Add New User" @submit="onSubmit">
      <v-card-text>
        <TextField label="Name" v-model="username" required />
        <TextField type="password" label="Password" v-model="password" required />
        <TextField type="password" label="Confirm Password" v-model="passwordConfirm" required />
      </v-card-text>

      <v-card-actions class="justify-end">
        <v-btn color="error" @click.prevent="$router.push({ name: 'settings_users_manage' })">Cancel</v-btn>
        <v-btn color="primary" type="submit">Save</v-btn>
      </v-card-actions>
    </FormCard>
  </div>
</template>


<script lang="ts">
import { defineComponent, ref } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import { success } from "@/plugins/toast";
import { useRouter } from "vue-router";

export default defineComponent({
  setup() {
    const globalStore = useGlobalStore();
    const router = useRouter();

    const isLoading = ref(false);

    const username = ref("");
    const password = ref("");
    const passwordConfirm = ref("");

    // --------------------------------------------------------------------------

    const onSubmit = () => {
      isLoading.value = true;

      if (password.value !== passwordConfirm.value) {
        globalStore.setError("Passwords do not match");
        return;
      }

      send("SETTINGS__USERS__ADD", {
        username: username.value,
        password: password.value,
      })
        .then(() => {
          success("User added successfully");
          router.push({ name: "settings_users_manage" });
        })
        .catch(({ payload, statusCode }: IResponse<string>) => globalStore.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => (isLoading.value = false));
    };

    // --------------------------------------------------------------------------

    return {
      isLoading,
      username,
      password,
      passwordConfirm,
      onSubmit,
    };
  },
});
</script>