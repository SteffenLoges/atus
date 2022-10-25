import { defineStore } from "pinia";
import { ref } from "vue";
import jwt_decode from "jwt-decode";
import { fetchInternally } from "@/utils/fetch";

export default defineStore("user", () => {
  const isAuthenticated = ref(false);
  const uid = ref("");
  const name = ref("");

  const getAuthToken = () =>
    localStorage.getItem("token") || "";

  const setAuthToken = (token: string) => {
    localStorage.setItem("token", token);
    const jwt: IJWTContent = jwt_decode(token);
    if (!jwt?.uid) {
      console.error("Invalid JWT", jwt);
      return;
    }
    uid.value = jwt.uid;
    name.value = jwt.name;
    isAuthenticated.value = true;
  };

  const authenticate = async () => {
    const { refreshToken }: IAuthResponse =
      await fetchInternally("/user/auth");

    setAuthToken(refreshToken);
  };

  const deleteAuthToken = () => {
    localStorage.removeItem("token");
    uid.value = "";
    name.value = "";
    isAuthenticated.value = false;
  };
  // --------------

  return {
    uid,
    name,
    isAuthenticated,
    authenticate,
    getAuthToken,
    setAuthToken,
    deleteAuthToken,
  };
});
