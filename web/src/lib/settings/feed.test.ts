import assert from "node:assert/strict";
import { test } from "vitest";

import {
  DefaultFeedConfig,
  FeedConfigSchema,
  addFeedBlock,
  removeFeedBlock,
  reorderFeedBlock,
  replaceFeedBlock,
} from "./feed";

test("missing feed configuration uses the fresh-install default", () => {
  const feed = FeedConfigSchema.parse(undefined);

  assert.deepStrictEqual(feed, DefaultFeedConfig);
  assert.deepStrictEqual(feed.blocks, [
    { type: "title" },
    { type: "subtitle" },
    { type: "content" },
    { type: "categories", layout: "grid" },
    { type: "threads", source: "uncategorised" },
  ]);
});

test("feed configuration preserves ordered blocks and applies block defaults", () => {
  const feed = FeedConfigSchema.parse({
    blocks: [
      { type: "title" },
      { type: "categories" },
      { type: "threads" },
      { type: "quick-share" },
      { type: "library" },
    ],
  });

  assert.deepStrictEqual(feed.blocks, [
    { type: "title" },
    { type: "categories", layout: "list" },
    { type: "threads", source: "uncategorised" },
    { type: "quick-share", showCategorySelect: true },
    { type: "library", layout: "list" },
  ]);
});

test("feed configuration rejects duplicate block types", () => {
  const result = FeedConfigSchema.safeParse({
    blocks: [{ type: "title" }, { type: "title" }],
  });

  assert.strictEqual(result.success, false);
});

test("legacy thread feeds preserve Quick Share visibility", () => {
  const enabled = FeedConfigSchema.parse({
    layout: { type: "list" },
    source: { type: "threads", quickShare: "enabled" },
  });
  const disabled = FeedConfigSchema.parse({
    layout: { type: "grid" },
    source: { type: "threads", quickShare: "disabled" },
  });

  assert.deepStrictEqual(enabled.blocks, [
    { type: "quick-share", showCategorySelect: true },
    { type: "threads", source: "all" },
  ]);
  assert.deepStrictEqual(disabled.blocks, [{ type: "threads", source: "all" }]);
});

test("legacy category feeds migrate their layout and thread source", () => {
  const feed = FeedConfigSchema.parse({
    layout: { type: "grid" },
    source: {
      type: "categories",
      quickShare: "enabled",
      threadListMode: "all",
    },
  });

  assert.deepStrictEqual(feed.blocks, [
    { type: "categories", layout: "grid" },
    { type: "quick-share", showCategorySelect: true },
    { type: "threads", source: "all" },
  ]);
});

test("legacy category feeds hide the category picker for uncategorised threads", () => {
  const feed = FeedConfigSchema.parse({
    layout: { type: "list" },
    source: { type: "categories" },
  });

  assert.deepStrictEqual(feed.blocks, [
    { type: "categories", layout: "list" },
    { type: "quick-share", showCategorySelect: false },
    { type: "threads", source: "uncategorised" },
  ]);
});

test("legacy category feeds without a thread list do not render Quick Share", () => {
  const feed = FeedConfigSchema.parse({
    layout: { type: "grid" },
    source: {
      type: "categories",
      quickShare: "enabled",
      threadListMode: "none",
    },
  });

  assert.deepStrictEqual(feed.blocks, [{ type: "categories", layout: "grid" }]);
});

test("legacy category feeds preserve a disabled Quick Share", () => {
  const feed = FeedConfigSchema.parse({
    layout: { type: "list" },
    source: {
      type: "categories",
      quickShare: "disabled",
      threadListMode: "all",
    },
  });

  assert.deepStrictEqual(feed.blocks, [
    { type: "categories", layout: "list" },
    { type: "threads", source: "all" },
  ]);
});

test("legacy library feeds migrate their page and layout", () => {
  const selectedPage = FeedConfigSchema.parse({
    layout: { type: "list" },
    source: { type: "library", node: "page-id" },
  });
  const libraryRoot = FeedConfigSchema.parse({
    layout: { type: "grid" },
    source: { type: "library" },
  });

  assert.deepStrictEqual(selectedPage.blocks, [
    { type: "library", node: "page-id", layout: "list" },
  ]);
  assert.deepStrictEqual(libraryRoot.blocks, [
    { type: "library", layout: "grid" },
  ]);
});

test("blocks can be inserted after another block or at the end", () => {
  const feed = FeedConfigSchema.parse({
    blocks: [{ type: "title" }, { type: "threads" }],
  });

  assert.deepStrictEqual(addFeedBlock(feed, "content", 0).blocks, [
    { type: "title" },
    { type: "content" },
    { type: "threads", source: "uncategorised" },
  ]);
  assert.deepStrictEqual(addFeedBlock(feed, "quick-share").blocks, [
    { type: "title" },
    { type: "threads", source: "uncategorised" },
    { type: "quick-share", showCategorySelect: true },
  ]);
});

test("block insertion clamps invalid indices and ignores duplicates", () => {
  const feed = FeedConfigSchema.parse({
    blocks: [{ type: "title" }, { type: "threads" }],
  });

  assert.deepStrictEqual(addFeedBlock(feed, "content", -10).blocks, [
    { type: "content" },
    { type: "title" },
    { type: "threads", source: "uncategorised" },
  ]);
  assert.deepStrictEqual(addFeedBlock(feed, "content", 100).blocks, [
    { type: "title" },
    { type: "threads", source: "uncategorised" },
    { type: "content" },
  ]);
  assert.strictEqual(addFeedBlock(feed, "title"), feed);
});

test("blocks reorder consistently in both directions", () => {
  const feed = FeedConfigSchema.parse({
    blocks: [
      { type: "title" },
      { type: "content" },
      { type: "categories" },
      { type: "threads" },
    ],
  });

  assert.deepStrictEqual(reorderFeedBlock(feed, "title", "categories").blocks, [
    { type: "content" },
    { type: "categories", layout: "list" },
    { type: "title" },
    { type: "threads", source: "uncategorised" },
  ]);
  assert.deepStrictEqual(reorderFeedBlock(feed, "threads", "content").blocks, [
    { type: "title" },
    { type: "threads", source: "uncategorised" },
    { type: "content" },
    { type: "categories", layout: "list" },
  ]);
  assert.strictEqual(reorderFeedBlock(feed, "title", "title"), feed);
});

test("blocks can be replaced and removed without changing their neighbours", () => {
  const feed = FeedConfigSchema.parse({
    blocks: [
      { type: "categories", layout: "list" },
      { type: "threads", source: "uncategorised" },
    ],
  });

  const replaced = replaceFeedBlock(feed, {
    type: "categories",
    layout: "grid",
  });
  assert.deepStrictEqual(replaced.blocks, [
    { type: "categories", layout: "grid" },
    { type: "threads", source: "uncategorised" },
  ]);
  assert.deepStrictEqual(removeFeedBlock(replaced, "categories").blocks, [
    { type: "threads", source: "uncategorised" },
  ]);
  assert.strictEqual(removeFeedBlock(feed, "cover"), feed);
});
