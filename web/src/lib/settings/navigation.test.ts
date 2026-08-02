import { test } from "uvu";
import * as assert from "uvu/assert";

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
  assert.equal(
    NavigationConfigSchema.parse(undefined),
    DefaultNavigationConfig,
  );
});

test("built-in items may only appear once", () => {
  const result = NavigationConfigSchema.safeParse({
    items: [{ type: "categories" }, { type: "categories" }],
  });

  assert.is(result.success, false);
});

test("custom links use their stable ID for identity", () => {
  const navigation = NavigationConfigSchema.parse({ items: [] });
  const updated = addCustomNavigationLink(navigation, {
    id: "docs",
    label: "Documentation",
    href: "https://docs.example.com/",
  });

  assert.equal(updated.items, [
    {
      type: "custom-link",
      id: "docs",
      label: "Documentation",
      href: "https://docs.example.com/",
    },
  ]);
  assert.is(
    addCustomNavigationLink(updated, {
      id: "docs",
      label: "Duplicate",
      href: "/duplicate",
    }),
    updated,
  );
  assert.is(getNavigationItemKey(updated.items[0]!), "custom-link:docs");
});

test("items can be inserted, reordered, replaced, and removed", () => {
  const initial = NavigationConfigSchema.parse({
    items: [{ type: "categories" }, { type: "library" }],
  });
  const inserted = addBuiltInNavigationItem(initial, "members", 0);

  assert.equal(inserted.items, [
    { type: "categories" },
    { type: "members" },
    { type: "library" },
  ]);
  assert.is(addBuiltInNavigationItem(inserted, "library"), inserted);

  const reordered = reorderNavigationItem(inserted, "library", "categories");
  assert.equal(reordered.items, [
    { type: "library" },
    { type: "categories" },
    { type: "members" },
  ]);
  assert.is(reorderNavigationItem(reordered, "missing", "library"), reordered);

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
  assert.equal(replaced.items.at(-1), {
    type: "custom-link",
    id: "community-guide",
    label: "Start here",
    href: "/l/start-here",
  });

  assert.equal(removeNavigationItem(replaced, "members").items, [
    { type: "library" },
    { type: "categories" },
    {
      type: "custom-link",
      id: "community-guide",
      label: "Start here",
      href: "/l/start-here",
    },
  ]);
  assert.is(removeNavigationItem(replaced, "missing"), replaced);
});

test("navigation URLs accept safe local and HTTP links", () => {
  assert.ok(isSafeNavigationHref("/about"));
  assert.ok(isSafeNavigationHref("?view=latest"));
  assert.ok(isSafeNavigationHref("#community"));
  assert.ok(isSafeNavigationHref("https://example.com/docs"));
  assert.not.ok(isSafeNavigationHref("//example.com"));
  assert.not.ok(isSafeNavigationHref("javascript:alert(1)"));
  assert.not.ok(isSafeNavigationHref("data:text/html,hello"));

  assert.is(
    normaliseNavigationHref("example.com/docs"),
    "https://example.com/docs",
  );
  assert.is(normaliseNavigationHref("not a link"), undefined);
});

test.run();
