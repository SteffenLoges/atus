<template>
  <tr>
    <td class="uid">{{ uid }}</td>
    <td class="name">{{ name }}</td>
    <td class="last-login d-none d-md-table-cell">{{ lastLoginHumanized }}</td>
    <td class="options text-right">
      <v-btn :to="{ name: 'settings_users_edit', params: { uid } }" class="mr-2" size="small" :icon="mdiPencil" flat />
      <v-tooltip location="bottom" :disabled="!isMaster && currentUserUID !== uid"
        :text="isMaster ? 'The master user can\'t be deleted' : 'You can\'t delete yourself'">
        <template v-slot:activator="{ props }">
          <span v-bind="props">
            <v-btn @click.prevent="$emit('delete')" :disabled="isMaster || currentUserUID === uid" size="small"
              :icon="mdiDelete" flat />
          </span>
        </template>
      </v-tooltip>
    </td>
  </tr>
</template>



<script lang="ts">
import { defineComponent, toRefs, computed } from "vue";
import { storeToRefs } from "pinia";
import { isValid } from "@/utils/date";
import useUserStore from "@/store/user";
import moment from "moment";
import { mdiDelete, mdiPencil } from "@mdi/js";

export default defineComponent({
  props: {
    uid: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      default: "",
    },
    lastLogin: {
      type: String,
      default: "",
    },
    isMaster: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["delete"],
  setup(props) {
    const { lastLogin } = toRefs(props);
    const userStore = useUserStore();
    const { uid: currentUserUID } = storeToRefs(userStore);

    const lastLoginHumanized = computed(() => {
      if (!isValid(lastLogin.value)) {
        return "Never";
      }

      return moment(lastLogin.value).fromNow();
    });

    return {
      lastLoginHumanized,
      currentUserUID,
      mdiDelete,
      mdiPencil,
    };
  },
});
</script>


<style lang="scss" scoped>
.uid {
  width: 170px;
}

.options {
  width: 150px;
}
</style>