import { z } from "zod";

export const FeedLayoutConfigSchema = z.union([
  z.object({
    type: z.literal("list"),
  }),
  z.object({
    type: z.literal("grid"),
  }),
]);
export type FeedLayoutConfig = z.infer<typeof FeedLayoutConfigSchema>;

export const FeedSourceConfigSchema = z.union([
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
export type FeedSourceConfig = z.infer<typeof FeedSourceConfigSchema>;

export const FeedConfigSchema = z.object({
  layout: FeedLayoutConfigSchema,
  source: FeedSourceConfigSchema,
});
export type FeedConfig = z.infer<typeof FeedConfigSchema>;

export const DefaultFeedConfig: FeedConfig = {
  layout: {
    type: "list",
  },
  source: {
    type: "threads",
    quickShare: "enabled",
  },
};
