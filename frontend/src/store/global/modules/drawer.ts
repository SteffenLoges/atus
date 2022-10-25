import { ref } from "vue";

export default () => {
  const forceDrawer = ref(false);

  const setForceDrawer = (v: boolean) =>
    (forceDrawer.value = v);

  const toggleForceDrawer = () =>
    (forceDrawer.value = !forceDrawer.value);

  return {
    forceDrawer,
    toggleForceDrawer,
    setForceDrawer,
  };
};
