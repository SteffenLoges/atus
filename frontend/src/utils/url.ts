import useUserStore from "@/store/user";

export const isSSL = location.protocol === "https:";

export const host =
  import.meta.env.VITE_APP_BACKEND_HOST ||
  location.hostname;

export const port =
  import.meta.env.VITE_APP_BACKEND_PORT || location.port;

export const baseURL = `${location.protocol}//${host}${
  port ? `:${port}` : ""
}`;

export const getFileURL = (
  path: string,
  includeToken = true
) => {
  const { getAuthToken } = useUserStore();

  let url = `${baseURL}/api/data/${path}`;

  if (includeToken) {
    url += `?token=${getAuthToken()}`;
  }

  return url;
};

export const dereferURL = (url: string, noSplash = false) =>
  "https://" +
  (noSplash ? "nosplash." : "") +
  "open-dereferrer.com/" +
  encodeURIComponent(url);
