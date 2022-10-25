import { createToast as _createToast } from "mosha-vue-toastify";

// maxiumum of 6 toast messages at a time
const MAX_TOASTS = 6;

const TIMEOUT = 5e3;

let activeToasts = 0;

interface IToast {
  type:
    | "info"
    | "danger"
    | "warning"
    | "success"
    | "default";
  title: string;
  message: string;
}

let queue: IToast[] = [];

// Warning: potential xss vulnerability - description renders html!
const createToast = (toast: IToast) => {
  if (activeToasts >= MAX_TOASTS) {
    queue.push(toast);
    return;
  }

  activeToasts++;

  _createToast(
    {
      title: toast.title,
      description: toast.message || undefined,
    },
    {
      type: toast.type,
      position: "bottom-right",
      showIcon: true,
      timeout: TIMEOUT,
      onClose: () => {
        activeToasts--;
        if (queue.length > 0) {
          const next = queue.shift();
          if (next) {
            setTimeout(() => createToast(next), 500);
          }
        }
      },
    }
  );
};

export default createToast;

export const info = (title: string, message = "") => {
  createToast({
    type: "info",
    title,
    message,
  });
};

export const success = (title: string, message = "") => {
  createToast({
    type: "success",
    title,
    message,
  });
};

export const warning = (title: string, message = "") => {
  createToast({
    type: "warning",
    title,
    message,
  });
};

export const error = (title: string, message = "") => {
  createToast({
    type: "danger",
    title,
    message,
  });
};
