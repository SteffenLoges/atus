<template>
  <ConfirmDialog v-model="showDialog" @always="showDialog = false" @confirm="onConfirm()"
    title="You have unsaved changes">
    Are you sure you want to leave?
  </ConfirmDialog>
</template> 


<script lang="ts">
import { ref, defineComponent, onBeforeUnmount, toRefs } from "vue";
import { useRouter } from "vue-router";
import { onBeforeRouteLeave } from "vue-router";

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:modelValue"],
  setup(props, { emit }) {
    const { modelValue: showDialogOnLeave } = toRefs(props);
    let allowLeave = ref(false);
    let showDialog = ref(false);
    let leaveTo = ref({ path: "/", });

    const router = useRouter();
    const onConfirm = () => {
      allowLeave.value = true;
      router.push(leaveTo.value);
    };

    // -- routeGuard ------------------------------------------------

    onBeforeRouteLeave((to, from, next) => {
      if (!allowLeave.value && showDialogOnLeave.value) {
        leaveTo.value = to;
        showDialog.value = true;
        next(false);
        return;
      }
      next();
    });

    // -- Handle browser events -------------------------------------

    const onBeforeUnload = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      e.returnValue = "You have unsaved changes. Are you sure you want to leave?";
      return false;
    };

    window.addEventListener("beforeunload", onBeforeUnload);
    onBeforeUnmount(() => window.removeEventListener("beforeunload", onBeforeUnload));

    // --------------------------------------------------------------

    return {
      allowLeave,
      showDialog,
      onConfirm,
    };
  },
});
</script>
