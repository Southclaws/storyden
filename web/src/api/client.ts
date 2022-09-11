import { API_ADDRESS } from "src/config";

type Options = {
  url: string;
  method: "get" | "post" | "put" | "delete" | "patch";
  headers?: Record<string, string>;
  params?: any;
  data?: unknown;
  responseType?: string;
};

export const fetcher = async <T>({
  url,
  method,
  headers,
  params,
  data,
}: Options): Promise<T> => {
  const req = new Request(
    `${API_ADDRESS}${url}` + new URLSearchParams(params),
    {
      method,
      mode: "cors",
      credentials: "include",
      ...(headers ? { headers } : {}),
      ...(data ? { body: JSON.stringify(data) } : {}),
    }
  );

  console.log(req);

  const response = await fetch(req);

  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.error);
  }

  return response.json();
};

export default fetcher;
