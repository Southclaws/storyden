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
  }),
  z.object({
    type: z.literal("library"),
  }),
]);
export type FeedSourceConfig = z.infer<typeof FeedSourceConfigSchema>;

export const DefaultFeedConfig = {
  layout: {
    type: "list",
  },
  source: {
    type: "threads",
  },
} as const;

export const FeedConfigSchema = z
  .object({
    layout: FeedLayoutConfigSchema,
    source: FeedSourceConfigSchema,
  })
  .default(DefaultFeedConfig);
export type FeedConfig = z.infer<typeof FeedConfigSchema>;
