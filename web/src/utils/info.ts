import { Info } from "@/api/openapi-schema";
import { getInfo as getInfoAPI } from "@/api/openapi-server/misc";

import { FALLBACK_COLOUR } from "./colour";

export async function getInfo(): Promise<Info> {
  try {
    const { data } = await getInfoAPI();
    return data;
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
