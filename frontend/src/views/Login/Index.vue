<template>
  <v-container fluid class="d-flex align-center justify-center h-100">
    <Card class="d-flex flex-grow-1 login-card px-4 py-3">
      <v-row no-gutters>
        <v-col cols="12" md="4" class="py-10 d-flex flex-column justify-space-evenly align-center">
          <img src="/src/assets/images/logo-big.png" class="logo" />
          <img src="/src/assets/images/logo-text.png" class="logo-text" />
        </v-col>
        <v-col cols="12" md="8" class="pa-0">
          <FormCard variant="text" @submit.prevent="onSubmit()" :disabled="loading" noConfirm title="Login">
            <v-card-text>
              <v-row>
                <v-col>
                  <TextField autocomplete="username" autofocus v-model="username" label="Username" :error="loginError"
                    dense hideDetails filled />
                </v-col>
              </v-row>
              <v-row>
                <v-col class="mt-3">
                  <TextField autocomplete="current-password" v-model="password" type="password" label="Password"
                    :error="loginError" dense hideDetails filled />
                </v-col>
              </v-row>
            </v-card-text>

            <v-card-actions class="px-4 pt-0 justify-center justify-md-start">
              <v-btn :loading="loading" :disabled="loading" type="submit" color="primary" block>Login</v-btn>
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

interface ILoginResponse {
  token: string;
}

export default defineComponent({
  setup() {
    const router = useRouter();
    const userStore = useUserStore();

    const username = ref("");
    const password = ref("");
    const loading = ref(false);
    const loginError = ref(false);

    const onSubmit = () => {
      loading.value = true;

      fetchInternally("/user/login", {
        method: "POST",
        body: JSON.stringify({
          username: username.value,
          password: password.value,
        }),
      })
        .then((r: ILoginResponse) => {
          userStore.setAuthToken(r.token);
          router.replace("/");
        })
        .catch((r) => {
          r.text()
            .then((t: string) => error("Login failed", t))
            .catch(() => error("Login failed", "Unknown error"))
          loginError.value = true
        })
        .finally(() => loading.value = false)
    };

    return {
      username,
      password,
      loading,
      loginError,
      onSubmit,
    };
  },
});
</script>


<style lang="scss" scoped>
.login-card {
  max-width: 700px;
}

.logo {
  width: 120px;
}

.logo-text {
  width: 120px;
}
</style>