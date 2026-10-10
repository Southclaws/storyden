import assert from "node:assert/strict";
import { test } from "vitest";

import {
  DefaultNavigationConfig,
  NavigationConfigSchema,
  addBuiltInNavigationItem,
  addCustomNavigationLink,
  getNavigationItemKey,
  isSafeNavigationHref,
  normaliseNavigationHref,
  removeNavigationItem,
  reorderNavigationItem,
  replaceCustomNavigationLink,
} from "./navigation";

test("missing navigation configuration uses the current sidebar", () => {
  assert.deepStrictEqual(
    NavigationConfigSchema.parse(undefined),
    DefaultNavigationConfig,
  );
});

test("robots is available as an optional built-in navigation item", () => {
  const navigation = NavigationConfigSchema.parse({ items: [] });
  assert.deepStrictEqual(addBuiltInNavigationItem(navigation, "robots").items, [
    { type: "robots" },
  ]);
});

test("new thread is the first default item and can be added", () => {
  assert.deepStrictEqual(DefaultNavigationConfig.items[0], { type: "compose" });

  const navigation = NavigationConfigSchema.parse({
    items: [{ type: "categories" }],
  });
  assert.deepStrictEqual(
    addBuiltInNavigationItem(navigation, "compose").items,
    [{ type: "categories" }, { type: "compose" }],
  );
});

test("the default navigation contains only first-class items", () => {
  assert.deepStrictEqual(DefaultNavigationConfig.items, [
    { type: "compose" },
    { type: "categories" },
    { type: "library" },
    { type: "collections" },
    { type: "links" },
    { type: "members" },
    { type: "roles" },
  ]);
});

test("built-in items may only appear once", () => {
  const result = NavigationConfigSchema.safeParse({
    items: [{ type: "categories" }, { type: "categories" }],
  });

  assert.strictEqual(result.success, false);
});

test("custom links use their stable ID for identity", () => {
  const navigation = NavigationConfigSchema.parse({ items: [] });
  const updated = addCustomNavigationLink(navigation, {
    id: "docs",
    label: "Documentation",
    href: "https://docs.example.com/",
  });

  assert.deepStrictEqual(updated.items, [
    {
      type: "custom-link",
      id: "docs",
      label: "Documentation",
      href: "https://docs.example.com/",
    },
  ]);
  assert.strictEqual(
    addCustomNavigationLink(updated, {
      id: "docs",
      label: "Duplicate",
      href: "/duplicate",
    }),
    updated,
  );
  assert.strictEqual(
    getNavigationItemKey(updated.items[0]!),
    "custom-link:docs",
  );
});

test("items can be inserted, reordered, replaced, and removed", () => {
  const initial = NavigationConfigSchema.parse({
    items: [{ type: "categories" }, { type: "library" }],
  });
  const inserted = addBuiltInNavigationItem(initial, "members", 0);

  assert.deepStrictEqual(inserted.items, [
    { type: "categories" },
    { type: "members" },
    { type: "library" },
  ]);
  assert.strictEqual(addBuiltInNavigationItem(inserted, "library"), inserted);

  const reordered = reorderNavigationItem(inserted, "library", "categories");
  assert.deepStrictEqual(reordered.items, [
    { type: "library" },
    { type: "categories" },
    { type: "members" },
  ]);
  assert.strictEqual(
    reorderNavigationItem(reordered, "missing", "library"),
    reordered,
  );

  const withLink = addCustomNavigationLink(reordered, {
    id: "community-guide",
    label: "Community guide",
    href: "/l/community-guide",
  });
  const replaced = replaceCustomNavigationLink(withLink, {
    type: "custom-link",
    id: "community-guide",
    label: "Start here",
    href: "/l/start-here",
  });
  assert.deepStrictEqual(replaced.items.at(-1), {
    type: "custom-link",
    id: "community-guide",
    label: "Start here",
    href: "/l/start-here",
  });

  assert.deepStrictEqual(removeNavigationItem(replaced, "members").items, [
    { type: "library" },
    { type: "categories" },
    {
      type: "custom-link",
      id: "community-guide",
      label: "Start here",
      href: "/l/start-here",
    },
  ]);
  assert.strictEqual(removeNavigationItem(replaced, "missing"), replaced);
});

test("navigation URLs accept safe local and HTTP links", () => {
  assert.ok(isSafeNavigationHref("/about"));
  assert.ok(isSafeNavigationHref("?view=latest"));
  assert.ok(isSafeNavigationHref("#community"));
  assert.ok(isSafeNavigationHref("https://example.com/docs"));
  assert.ok(!isSafeNavigationHref("//example.com"));
  assert.ok(!isSafeNavigationHref("javascript:alert(1)"));
  assert.ok(!isSafeNavigationHref("data:text/html,hello"));

  assert.strictEqual(
    normaliseNavigationHref("example.com/docs"),
    "https://example.com/docs",
  );
  assert.strictEqual(normaliseNavigationHref("not a link"), undefined);
});
