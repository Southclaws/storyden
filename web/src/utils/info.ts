import { GetInfoOKResponse } from "src/api/openapi/schemas";
import { API_ADDRESS } from "src/config";

export async function getInfo(): Promise<GetInfoOKResponse> {
  const res = await fetch(`${API_ADDRESS}/api/v1/info`);
  if (!res.ok) {
    throw new Error(
      `failed to fetch API info endpoint: ${res.status} ${res.statusText}`
    );
  }

  const info = (await res.json()) as GetInfoOKResponse;

  return info;
}
