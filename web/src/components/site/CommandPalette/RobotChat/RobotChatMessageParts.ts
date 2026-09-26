import { Asset } from "@/api/openapi-schema";
import { StorydenUIMessage } from "@/api/robots-types";

export function buildUserMessageParts(
  text: string,
  assets: readonly Asset[],
): StorydenUIMessage["parts"] {
  return [
    ...(text ? [{ type: "text" as const, text }] : []),
    ...assets.map((asset) => ({
      type: "data-storyden-asset" as const,
      data: { asset_id: asset.id },
    })),
  ];
}
