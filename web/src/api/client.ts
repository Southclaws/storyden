import { Options, buildPayload, cleanQuery, shouldLog } from "./common";

export const fetcher = async <T>({
  url,
  method = "GET",
  headers,
  params,
  data,
}: Options): Promise<T> => {
  const req = new Request(`/api${url}${cleanQuery(params)}`, {
    // NOTE: this is forced uppercase due to a bug somewhere in another part of
    // the code. It might be Orval as it generates all lowercase methods for all
    // requests, however this seems to work fine with every operation except the
    // PATCH calls. Not really sure if that's a browser issue or not though...
    method: method.toUpperCase(),
    mode: "cors",
    credentials: "include",
    ...(headers ? { headers } : {}),
    body: buildPayload(data),
  });

  const response = await fetch(req);

  if (!response.ok) {
    const data = await response
      .json()
      .catch(() => ({ error: "Failed to parse API response" }));

    if (shouldLog(response.status)) {
      console.warn({
        ...data,
        status: response.status,
        statusText: response.statusText,
      });
    }

    throw new Error(
      data.message ??
        `An unexpected error occurred:  ${response.status} ${response.statusText}`,
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
