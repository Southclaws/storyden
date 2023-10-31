import { API_ADDRESS } from "src/config";

import { Options, buildPayload, cleanQuery } from "./common";

export const server = async <T>({
  url,
  params,
  method = "get",
  data,
}: Options): Promise<T> => {
  const address = `${API_ADDRESS}/api${url}${cleanQuery(params)}`;
  const _method = method.toUpperCase();

  const response = await fetch(address, {
    method: _method,
    headers: { "Content-Type": "application/json" },
    body: buildPayload(data),
  });

  if (!response.ok) {
    const data = await response
      .json()
      .catch(() => ({ error: "Failed to parse API response" }));
    console.warn(data);
    throw new Error(
      data.message ??
        `An unexpected error occurred: ${response.status} ${response.statusText}`,
    );
  }

  // NOTE: The API code generator returns empty responses where there is no
  // response type specified with a content type so this is the easy way to
  // escape that code path and exit easily.
  if (response.headers.get("content-length") === "0") {
    return undefined as T;
  }

  return response.json();
};
