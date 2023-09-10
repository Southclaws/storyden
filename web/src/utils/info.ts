import { GetInfoOKResponse } from "src/api/openapi/schemas";
import { API_ADDRESS } from "src/config";

import { FALLBACK_COLOUR } from "./colour";

export async function getInfo(): Promise<GetInfoOKResponse> {
  try {
    const res = await fetch(`${API_ADDRESS}/api/v1/info`);

    if (!res.ok) {
      throw new Error(
        `failed to fetch API info endpoint: ${res.status} ${res.statusText}`,
      );
    }

    const info = (await res.json()) as GetInfoOKResponse;

    return info;
  } catch (e) {
    console.error(e);
    return {
      title: "Storyden",
      description: "A forum for the modern age.",
      accent_colour: FALLBACK_COLOUR,
      onboarding_status: "complete",
    };
  }
}
