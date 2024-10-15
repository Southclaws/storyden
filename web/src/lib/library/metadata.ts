import { z } from "zod";

import { Metadata } from "@/api/openapi-schema";

export const CoverImageSchema = z.object({
  top: z.number(),
  left: z.number(),
});
export type CoverImage = z.infer<typeof CoverImageSchema>;

export const NodeMetadataSchema = z.object({
  coverImage: CoverImageSchema.optional(),
});
export type NodeMetadata = z.infer<typeof NodeMetadataSchema>;

export function parseNodeMetadata(
  metadata: Metadata | undefined,
): NodeMetadata {
  const parsed = NodeMetadataSchema.safeParse(metadata);
  if (parsed.success) {
    return parsed.data;
  }

  return {};
}
