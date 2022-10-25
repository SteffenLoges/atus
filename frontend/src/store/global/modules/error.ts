import { ref } from "vue";

export default () => {
  const error = ref<any>(null);

  const setError = (err: any, critical = false) =>
    (error.value = { err, critical });

  const clearError = () => (error.value = null);

  return {
    error,
    setError,
    clearError,
  };
};
