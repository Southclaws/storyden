import { z } from "zod";

import { Metadata } from "@/api/openapi-schema";

const RoleMetadataParseSchema = z.object({
  bold: z.boolean().optional(),
  italic: z.boolean().optional(),
  coloured: z.boolean().optional(),
});

export const RoleMetadataSchema = z.object({
  bold: z.boolean(),
  italic: z.boolean(),
  coloured: z.boolean(),
});
export type RoleMetadata = z.infer<typeof RoleMetadataSchema>;

export const DefaultRoleMetadata: RoleMetadata = {
  bold: false,
  italic: false,
  coloured: false,
};

export function parseRoleMetadata(meta: Metadata | undefined): RoleMetadata {
  const parsed = RoleMetadataParseSchema.safeParse(meta ?? {});

  if (!parsed.success && meta) {
    console.warn("Failed to parse role metadata:", parsed.error);
  }

  return {
    bold: parsed.success ? parsed.data.bold ?? DefaultRoleMetadata.bold : false,
    italic: parsed.success
      ? parsed.data.italic ?? DefaultRoleMetadata.italic
      : false,
    coloured: parsed.success
      ? parsed.data.coloured ?? DefaultRoleMetadata.coloured
      : false,
  };
}

export function writeRoleMetadata(
  existing: Metadata | undefined,
  next: RoleMetadata,
): Metadata {
  return {
    ...(existing ?? {}),
    bold: next.bold,
    italic: next.italic,
    coloured: next.coloured,
  };
}
