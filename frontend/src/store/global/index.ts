import { defineStore } from "pinia";
import useNotifications from "./modules/notifications";
import useDrawer from "./modules/drawer";
import useFileservers from "./modules/fileservers";
import useError from "./modules/error";
import useSetup from "./modules/setup";

export default defineStore("global", () => ({
  ...useNotifications(),
  ...useSetup(),
  ...useDrawer(),
  ...useFileservers(),
  ...useError(),
}));
