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
  directory: "Directory",
  tags: "Tag list",
};

export const LibraryPageBlockTypeSchema = z.enum([
  "title",
  "cover",
  "link",
  "content",
  "assets",
  "properties",
  "directory",
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

export const LibraryPageBlockTypeAssetsLayoutSchema = z.enum(["strip", "grid"]);
export type LibraryPageBlockTypeAssetsLayout = z.infer<
  typeof LibraryPageBlockTypeAssetsLayoutSchema
>;

export const LibraryPageBlockTypeAssetsConfigSchema = z.object({
  layout: LibraryPageBlockTypeAssetsLayoutSchema,
  gridSize: z.number().optional(),
});
export type LibraryPageBlockTypeAssetsConfig = z.infer<
  typeof LibraryPageBlockTypeAssetsConfigSchema
>;

export const LibraryPageBlockTypeAssetsSchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.assets),
  config: LibraryPageBlockTypeAssetsConfigSchema.optional(),
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

export const LibraryPageBlockTypeDirectoryLayoutSchema = z.enum([
  "table",
  "grid",
]);
export type LibraryPageBlockTypeDirectoryLayout = z.infer<
  typeof LibraryPageBlockTypeDirectoryLayoutSchema
>;

export const LibraryPageBlockTypeDirectoryColumnSchema = z.object({
  fid: z.string(),
  hidden: z.boolean().default(false),
});
export type LibraryPageBlockTypeDirectoryColumn = z.infer<
  typeof LibraryPageBlockTypeDirectoryColumnSchema
>;

export const LibraryPageBlockTypeDirectoryConfigSchema = z.object({
  layout: LibraryPageBlockTypeDirectoryLayoutSchema,
  columns: z.array(LibraryPageBlockTypeDirectoryColumnSchema),
});
export type LibraryPageBlockTypeDirectoryConfig = z.infer<
  typeof LibraryPageBlockTypeDirectoryConfigSchema
>;

export const LibraryPageBlockTypeDirectorySchema = z.object({
  type: z.literal(LibraryPageBlockTypeSchema.Enum.directory),
  config: LibraryPageBlockTypeDirectoryConfigSchema.optional(),
});
export type LibraryPageBlockTypeDirectory = z.infer<
  typeof LibraryPageBlockTypeDirectorySchema
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
  LibraryPageBlockTypeDirectorySchema,
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
    { type: "link" as const },
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
