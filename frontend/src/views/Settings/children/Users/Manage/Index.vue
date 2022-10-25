<template>
  <ConfirmDialog :modelValue="deleteUID != ''" @cancel="deleteUID = ''" @confirm="onDeleteConfirm(deleteUID)">
    Are you sure you want to delete this user?
  </ConfirmDialog>

  <Card :loading="isLoading" title="Users">
    <v-card-text>
      <v-table>
        <thead>
          <tr>
            <th>UID</th>
            <th>Name</th>
            <th class="d-none d-md-table-cell">Last Login</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <User v-for="(user, i) in users" :class="{ 'mt-1': i > 0 }" :key="user.uid" v-bind="user"
            @delete="deleteUID = user.uid" />
        </tbody>
      </v-table>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="primary" :to="{ name: 'settings_users_add' }">Add New User</v-btn>
    </v-card-actions>
  </Card>
</template>


<script lang="ts">
import { ref, defineComponent } from "vue";
import { send } from "@/utils/websocket";
import useGlobalStore from "@/store/global";
import User from "./components/User.vue";
import { success } from "@/plugins/toast";

export default defineComponent({
  components: {
    User,
  },
  async setup() {
    const globalStore = useGlobalStore();
    const isLoading = ref(true);
    const users = ref<IUser[]>([]);

    // --------------------------------------------------------------------------

    const load = () =>
      send("SETTINGS__USERS__GET_ALL")
        .then(({ payload }: IResponse<IUser[]>) => (users.value = payload))
        .catch(({ payload, statusCode }: IResponse<string>) => globalStore.setError(`Server returned status code ${statusCode} with message: ${payload}`))
        .finally(() => isLoading.value = false)

    await load();

    // --------------------------------------------------------------------------

    const deleteUID = ref("");
    const onDeleteConfirm = (uid: string) => {
      deleteUID.value = "";
      isLoading.value = true;

      send("SETTINGS__USERS__DELETE", { uid })
        .then(() => {
          load().then(() => isLoading.value = false)
          success("User deleted successfully");
        })
        .catch(({ payload }: IResponse<string>) => {
          isLoading.value = false;
          globalStore.setError(payload);
        });
    };

    // --------------------------------------------------------------------------

    return {
      users,
      isLoading,
      deleteUID,
      onDeleteConfirm,
    };
  },
});
</script>
