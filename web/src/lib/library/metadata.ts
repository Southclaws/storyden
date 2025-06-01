import { z } from "zod";

import { Metadata, Node, NodeWithChildren } from "@/api/openapi-schema";

export const CoverImageSchema = z.object({
  top: z.number(),
  left: z.number(),
});
export type CoverImage = z.infer<typeof CoverImageSchema>;

export const LibraryPageBlockName: Record<LibraryPageBlockType, string> = {
  title: "Title",
  cover: "Cover image",
  link: "External link",
  content: "Rich text content",
  assets: "Gallery",
  properties: "Page properties",
  table: "Child page table",
  tags: "Tag list",
};

export const LibraryPageBlockTypeSchema = z.enum([
  "title",
  "cover",
  "link",
  "content",
  "assets",
  "properties",
  "table",
  "tags",
]);
export type LibraryPageBlockType = z.infer<typeof LibraryPageBlockTypeSchema>;

export const LibraryPageBlockSchema = z.object({
  type: LibraryPageBlockTypeSchema,
});
export type LibraryPageBlock = z.infer<typeof LibraryPageBlockSchema>;

export const NodeLayoutSchema = z.object({
  blocks: z.array(LibraryPageBlockSchema),
});
export type NodeLayout = z.infer<typeof NodeLayoutSchema>;

export const NodeMetadataSchema = z.object({
  coverImage: CoverImageSchema.optional(),
  layout: NodeLayoutSchema.optional(),
});
export type NodeMetadata = z.infer<typeof NodeMetadataSchema>;

const DefaultLayout: NodeLayout = {
  blocks: [
    { type: "cover" as const },
    { type: "title" as const },
    { type: "content" as const },
    { type: "assets" as const },
    { type: "tags" as const },
    { type: "link" as const },
    { type: "properties" as const },
    { type: "table" as const },
  ],
};

export function parseNodeMetadata(
  metadata: Metadata | undefined,
): NodeMetadata {
  const parsed = NodeMetadataSchema.safeParse(metadata);
  if (parsed.success) {
    if (!parsed.data.layout) {
      parsed.data.layout = DefaultLayout;
    }

    return parsed.data;
  }

  return {
    layout: DefaultLayout,
  };
}

export function hydrateNode<T extends Node | NodeWithChildren>(
  node: T,
): WithMetadata<T> {
  const meta = parseNodeMetadata(node.meta);
  return {
    ...node,
    meta,
  };
}

export type WithMetadata<T extends Node | NodeWithChildren> = T & {
  meta: NodeMetadata;
};
