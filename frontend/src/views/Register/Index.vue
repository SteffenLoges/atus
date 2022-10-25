<template>
  <v-container fluid class="d-flex align-center justify-center h-100">
    <Card class="d-flex flex-grow-1 register-card px-0 py-3 flex-column">
      <template #title>
        <v-alert color="info" class="text-subtitle-1 mr-0 ml-4 pa-3 text-center">
          <h2>Welcome to <strong>{{appName}}</strong>!</h2>
          <p class="mt-1">
            <span class="d-none d-md-inline">Looks like you don't have an account yet.</span>
            Please register below.
          </p>
        </v-alert>
      </template>
      <v-row no-gutters class="px-4">
        <v-col cols="12" md="4" class="py-12 d-flex flex-column justify-space-evenly align-center">
          <img src="/src/assets/images/logo-big.png" class="logo mt-md-n6" />
          <img src="/src/assets/images/logo-text.png" class="logo-text mt-md-n6" />
        </v-col>
        <v-col cols="12" md="8" class="pa-0">
          <FormCard variant="text" @submit.prevent="onSubmit()" :disabled="loading" noConfirm>
            <v-card-text>
              <v-row>
                <v-col>
                  <TextField autocomplete="username" autofocus v-model="username" label="Username"
                    :error="registerError" dense hideDetails filled />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="mt-3">
                  <TextField autocomplete="current-password" v-model="password" type="password" label="Password"
                    :error="registerError" dense hideDetails filled />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="mt-3">
                  <TextField autocomplete="current-password" v-model="passwordConfirm" type="password"
                    label="Confirm Password" :error="registerError" dense hideDetails filled />
                </v-col>
              </v-row>
            </v-card-text>

            <v-card-actions class="px-4 pt-0 justify-center justify-md-start">
              <v-btn :loading="loading" :disabled="loading" type="submit" color="primary" block>
                Create Master Account
              </v-btn>
            </v-card-actions>
          </FormCard>
        </v-col>
      </v-row>
    </Card>
  </v-container>
</template>


<script lang="ts">
import { defineComponent, ref } from "vue";
import { fetchInternally } from "@/utils/fetch";
import { error } from "@/plugins/toast";
import useUserStore from "@/store/user";
import { useRouter } from "vue-router";

interface IRegisterResponse {
  token: string;
}

export default defineComponent({
  setup() {
    const appName = import.meta.env.VITE_APP_NAME;
    const router = useRouter();
    const userStore = useUserStore();

    const username = ref("");
    const password = ref("");
    const passwordConfirm = ref("");
    const loading = ref(false);
    const registerError = ref(false);

    const onSubmit = () => {

      if (password.value !== passwordConfirm.value) {
        error("Passwords do not match");
        return;
      }

      loading.value = true;

      fetchInternally("/user/register", {
        method: "POST",
        body: JSON.stringify({
          username: username.value,
          password: password.value,
        }),
      })
        .then((r: IRegisterResponse) => {
          userStore.setAuthToken(r.token);
          router.replace({ name: 'setup' });
        })
        .catch((r) => {
          r.text()
            .then((t: string) => error("Failed to create account", t))
            .catch(() => error("Failed to create account", "Unknown error"));
          registerError.value = true;
        })
        .finally(() => loading.value = false)
    };

    return {
      appName,
      username,
      password,
      passwordConfirm,
      loading,
      registerError,
      onSubmit,
    };
  },
});
</script>


<style lang="scss" scoped>
.register-card {
  max-width: 700px;
}

.logo {
  width: 120px;
}

.logo-text {
  width: 120px;
}
</style>