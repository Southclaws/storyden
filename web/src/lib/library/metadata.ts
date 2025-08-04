import { z } from "zod";

import { Metadata, Node, NodeWithChildren, React } from "@/api/openapi-schema";

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

// -
// Block type schemas
// -

export const LibraryPageBlockTypeTitleSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.title),
});
export type LibraryPageBlockTypeTitle = z.infer<
  typeof LibraryPageBlockTypeTitleSchema
>;

export const LibraryPageBlockTypeCoverSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.cover),
});
export type LibraryPageBlockTypeCover = z.infer<
  typeof LibraryPageBlockTypeCoverSchema
>;

export const LibraryPageBlockTypeLinkSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.link),
});
export type LibraryPageBlockTypeLink = z.infer<
  typeof LibraryPageBlockTypeLinkSchema
>;

export const LibraryPageBlockTypeContentSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.content),
});
export type LibraryPageBlockTypeContent = z.infer<
  typeof LibraryPageBlockTypeContentSchema
>;

export const LibraryPageBlockTypeAssetsSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.assets),
});
export type LibraryPageBlockTypeAssets = z.infer<
  typeof LibraryPageBlockTypeAssetsSchema
>;

export const LibraryPageBlockTypePropertiesSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.properties),
});
export type LibraryPageBlockTypeProperties = z.infer<
  typeof LibraryPageBlockTypePropertiesSchema
>;

export const LibraryPageBlockTypeTableColumnSchema = z.object({
  fid: z.string(),
  hidden: z.boolean(),
});
export type LibraryPageBlockTypeTableColumn = z.infer<
  typeof LibraryPageBlockTypeTableColumnSchema
>;

export const LibraryPageBlockTypeTableConfigSchema = z.object({
  columns: z.array(LibraryPageBlockTypeTableColumnSchema),
});
export type LibraryPageBlockTypeTableConfig = z.infer<
  typeof LibraryPageBlockTypeTableConfigSchema
>;

export const LibraryPageBlockTypeTableSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.table),
  config: LibraryPageBlockTypeTableConfigSchema.optional(),
});
export type LibraryPageBlockTypeTable = z.infer<
  typeof LibraryPageBlockTypeTableSchema
>;

export const LibraryPageBlockTypeTagsSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.tags),
});
export type LibraryPageBlockTypeTags = z.infer<
  typeof LibraryPageBlockTypeTagsSchema
>;

// -
// Block type union
// -

export const LibraryPageBlockSchema = z.union([
  LibraryPageBlockTypeTitleSchema,
  LibraryPageBlockTypeCoverSchema,
  LibraryPageBlockTypeLinkSchema,
  LibraryPageBlockTypeContentSchema,
  LibraryPageBlockTypeAssetsSchema,
  LibraryPageBlockTypePropertiesSchema,
  LibraryPageBlockTypeTableSchema,
  LibraryPageBlockTypeTagsSchema,
]);
export type LibraryPageBlock = z.infer<typeof LibraryPageBlockSchema>;

export const NodeLayoutSchema = z.object({
  blocks: z.array(LibraryPageBlockSchema),
});
export type NodeLayout = z.infer<typeof NodeLayoutSchema>;

export const NodeMetadataSchema = z.object({
  coverImage: CoverImageSchema.optional().nullable(),
  layout: NodeLayoutSchema.optional(),
});
export type NodeMetadata = z.infer<typeof NodeMetadataSchema>;

export const DefaultLayout: NodeLayout = {
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
