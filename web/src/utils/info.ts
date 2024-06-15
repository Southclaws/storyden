import { GetInfoOKResponse } from "@/api/openapi/schemas/getInfoOKResponse";
import { server } from "@/api/server";

import { FALLBACK_COLOUR } from "./colour";

export async function getInfo(): Promise<GetInfoOKResponse> {
  try {
    return await server<GetInfoOKResponse>({ url: "/v1/info" });
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
