import { z } from "zod";

export const FeedBlockTypeSchema = z.enum([
  "title",
  "subtitle",
  "content",
  "cover",
  "categories",
  "threads",
  "quick-share",
  "library",
]);
export type FeedBlockType = z.infer<typeof FeedBlockTypeSchema>;

const FeedTitleBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.title),
});

const FeedSubtitleBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.subtitle),
});

const FeedContentBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.content),
});

const FeedCoverBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.cover),
});

const FeedCategoriesBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.categories),
  layout: z.enum(["list", "grid"]).default("list"),
});

const FeedThreadsBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.threads),
  source: z.enum(["all", "uncategorised"]).default("uncategorised"),
});

const FeedQuickShareBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum["quick-share"]),
  showCategorySelect: z.boolean().default(true),
});

const FeedLibraryBlockSchema = z.object({
  type: z.literal(FeedBlockTypeSchema.enum.library),
  node: z.string().optional(),
  layout: z.enum(["list", "grid"]).default("list"),
});

export const FeedBlockSchema = z.discriminatedUnion("type", [
  FeedTitleBlockSchema,
  FeedSubtitleBlockSchema,
  FeedContentBlockSchema,
  FeedCoverBlockSchema,
  FeedCategoriesBlockSchema,
  FeedThreadsBlockSchema,
  FeedQuickShareBlockSchema,
  FeedLibraryBlockSchema,
]);
export type FeedBlock = z.infer<typeof FeedBlockSchema>;

const OrderedFeedConfigSchema = z
  .object({
    blocks: z.array(FeedBlockSchema),
  })
  .superRefine(({ blocks }, context) => {
    const seen = new Set<FeedBlockType>();

    blocks.forEach((block, index) => {
      if (seen.has(block.type)) {
        context.addIssue({
          code: z.ZodIssueCode.custom,
          message: `Feed block "${block.type}" may only appear once.`,
          path: ["blocks", index],
        });
      }
      seen.add(block.type);
    });
  });
export type FeedConfig = z.infer<typeof OrderedFeedConfigSchema>;

const LegacyFeedLayoutConfigSchema = z.object({
  type: z.enum(["list", "grid"]),
});

const LegacyFeedSourceConfigSchema = z.discriminatedUnion("type", [
  z.object({
    type: z.literal("threads"),
    quickShare: z.enum(["enabled", "disabled"]).default("enabled"),
  }),
  z.object({
    type: z.literal("library"),
    node: z.string().optional(),
  }),
  z.object({
    type: z.literal("categories"),
    threadListMode: z
      .enum(["none", "all", "uncategorised"])
      .default("uncategorised"),
    quickShare: z.enum(["enabled", "disabled"]).default("enabled"),
  }),
]);

const LegacyFeedConfigSchema = z.object({
  layout: LegacyFeedLayoutConfigSchema,
  source: LegacyFeedSourceConfigSchema,
});

export const DefaultFeedConfig = {
  blocks: [
    { type: "title" },
    { type: "subtitle" },
    { type: "content" },
    { type: "categories", layout: "grid" },
    { type: "threads", source: "uncategorised" },
  ],
} satisfies FeedConfig;

export const FeedConfigSchema = z
  .preprocess((value) => {
    const current = OrderedFeedConfigSchema.safeParse(value);
    if (current.success) {
      return current.data;
    }

    const legacy = LegacyFeedConfigSchema.safeParse(value);
    if (legacy.success) {
      return migrateLegacyFeedConfig(legacy.data);
    }

    return value;
  }, OrderedFeedConfigSchema)
  .default(DefaultFeedConfig);

type LegacyFeedConfig = z.infer<typeof LegacyFeedConfigSchema>;

function migrateLegacyFeedConfig(feed: LegacyFeedConfig): FeedConfig {
  const { source } = feed;

  switch (source.type) {
    case "threads":
      return {
        blocks: [
          ...(source.quickShare === "enabled"
            ? ([
                {
                  type: "quick-share" as const,
                  showCategorySelect: true,
                },
              ] satisfies FeedBlock[])
            : []),
          { type: "threads", source: "all" },
        ],
      };

    case "categories":
      return {
        blocks: [
          { type: "categories", layout: feed.layout.type },
          ...(source.threadListMode !== "none" &&
          source.quickShare === "enabled"
            ? ([
                {
                  type: "quick-share" as const,
                  showCategorySelect: source.threadListMode === "all",
                },
              ] satisfies FeedBlock[])
            : []),
          ...(source.threadListMode === "none"
            ? []
            : ([
                {
                  type: "threads" as const,
                  source: source.threadListMode,
                },
              ] satisfies FeedBlock[])),
        ],
      };

    case "library":
      return {
        blocks: [
          {
            type: "library",
            ...(source.node ? { node: source.node } : {}),
            layout: feed.layout.type,
          },
        ],
      };
  }
}

export const FeedBlockName: Record<FeedBlockType, string> = {
  title: "Title",
  subtitle: "Subtitle",
  content: "Content",
  cover: "Cover image",
  categories: "Discussion categories",
  threads: "Thread list",
  "quick-share": "Quick share",
  library: "Library page",
};

export const AllFeedBlockTypes = FeedBlockTypeSchema.options;

export function createFeedBlock(type: FeedBlockType): FeedBlock {
  switch (type) {
    case "categories":
      return { type, layout: "list" };
    case "threads":
      return { type, source: "uncategorised" };
    case "quick-share":
      return { type, showCategorySelect: true };
    case "library":
      return { type, layout: "list" };
    default:
      return { type };
  }
}

export function addFeedBlock(
  feed: FeedConfig,
  type: FeedBlockType,
  afterIndex?: number,
): FeedConfig {
  if (feed.blocks.some((block) => block.type === type)) {
    return feed;
  }

  const blocks = [...feed.blocks];
  const insertionIndex =
    afterIndex === undefined
      ? blocks.length
      : Math.min(Math.max(afterIndex + 1, 0), blocks.length);

  blocks.splice(insertionIndex, 0, createFeedBlock(type));

  return { blocks };
}

export function reorderFeedBlock(
  feed: FeedConfig,
  active: FeedBlockType,
  over: FeedBlockType,
): FeedConfig {
  const activeIndex = feed.blocks.findIndex((block) => block.type === active);
  const overIndex = feed.blocks.findIndex((block) => block.type === over);
  if (activeIndex === -1 || overIndex === -1 || activeIndex === overIndex) {
    return feed;
  }

  const blocks = [...feed.blocks];
  const [moved] = blocks.splice(activeIndex, 1);
  if (!moved) {
    return feed;
  }

  blocks.splice(overIndex, 0, moved);

  return { blocks };
}

export function replaceFeedBlock(
  feed: FeedConfig,
  updated: FeedBlock,
): FeedConfig {
  const index = feed.blocks.findIndex((block) => block.type === updated.type);
  if (index === -1) {
    return feed;
  }

  const blocks = [...feed.blocks];
  blocks[index] = updated;

  return { blocks };
}

export function removeFeedBlock(
  feed: FeedConfig,
  type: FeedBlockType,
): FeedConfig {
  if (!feed.blocks.some((block) => block.type === type)) {
    return feed;
  }

  return {
    blocks: feed.blocks.filter((block) => block.type !== type),
  };
}
