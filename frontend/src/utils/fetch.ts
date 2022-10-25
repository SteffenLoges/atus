import useUserStore from "@/store/user";
import { baseURL } from "@/utils/url";

export const fetchInternally = async (
  path: string,
  init?: RequestInit
) => {
  const { getAuthToken } = useUserStore();

  let headers = (init?.headers as Headers) || new Headers();
  headers.append("Content-Type", "application/json");

  const token = localStorage.getItem("token");
  if (token) {
    headers.append(
      "Authorization",
      `Bearer ${getAuthToken()}`
    );
  }

  const controller = new AbortController();

  const timeout = setTimeout(() => controller.abort(), 5e3);

  const options: RequestInit = {
    ...init,
    headers,
    signal: controller.signal,
  };

  const res = await fetch(
    `${baseURL}/frontend/api${path}`,
    options
  );

  clearTimeout(timeout);

  if (res.status === 200) {
    return res.json();
  }

  throw res;
};
