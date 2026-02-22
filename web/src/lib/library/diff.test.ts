import { test } from "uvu";
import * as assert from "uvu/assert";

import {
  type Asset,
  type LinkReference,
  type NodeWithChildren,
  PropertyType,
} from "@/api/openapi-schema";

import { deriveMutationFromDifference } from "./diff";

const consoleDebug = console.debug;
const consoleWarn = console.warn;

test.before(() => {
  console.debug = () => {};
  console.warn = () => {};
});

test.after(() => {
  console.debug = consoleDebug;
  console.warn = consoleWarn;
});

function node(
  id: string,
  overrides: Partial<NodeWithChildren> = {},
): NodeWithChildren {
  return {
    id,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    name: `node-${id}`,
    slug: `node-${id}`,
    content: "<p>hello</p>",
    tags: [{ name: "general" }],
    primary_image: asset(`img-${id}`),
    link: linkReference("https://example.com"),
    description: "desc",
    properties: [
      {
        fid: "field-1",
        name: "Label",
        value: "Value",
        type: PropertyType.text,
      },
    ],
    hide_child_tree: false,
    meta: {},
    child_property_schema: [
      {
        fid: "field-1",
        name: "Label",
        sort: "0",
        type: PropertyType.text,
      },
    ],
    children: [],
    ...overrides,
  } as NodeWithChildren;
}

function asset(id: string): Asset {
  return {
    id,
    filename: `${id}.png`,
    height: 64,
    mime_type: "image/png",
    path: `/assets/${id}`,
    width: 64,
  };
}

function linkReference(url: string): LinkReference {
  return {
    id: "link-1",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    domain: "example.com",
    slug: "example",
    url,
  } as LinkReference;
}

test("deriveMutationFromDifference returns clean when nodes are identical", () => {
  const current = node("root");
  const updated = node("root");

  const result = deriveMutationFromDifference(current, updated);

  assert.ok(result.clean);
  assert.equal(result.nodeMutation, {});
  assert.equal(result.childMutation, {});
  assert.is(result.childPropertySchemaMutation, undefined);
});

test("deriveMutationFromDifference skips invalid slug changes", () => {
  const current = node("root", { slug: "valid-slug" });
  const updated = node("root", { slug: "invalid slug" });

  const result = deriveMutationFromDifference(current, updated);

  assert.ok(result.clean);
  assert.equal(result.nodeMutation, {});
});

test("deriveMutationFromDifference sets url to null when cleared", () => {
  const current = node("root", {
    link: linkReference("https://example.com/path"),
  });
  const updated = node("root", { link: undefined });

  const result = deriveMutationFromDifference(current, updated);

  assert.not.ok(result.clean);
  assert.equal(result.nodeMutation, { url: null });
});

test("deriveMutationFromDifference does not clear url on invalid updated url", () => {
  const current = node("root", {
    link: linkReference("https://example.com/path"),
  });
  const updated = node("root", {
    link: {
      ...linkReference("https://example.com/path"),
      url: "not a valid url",
    },
  });

  const result = deriveMutationFromDifference(current, updated);

  assert.ok(result.clean);
  assert.equal(result.nodeMutation, {});
});

test("deriveMutationFromDifference sets primary image to null when cleared", () => {
  const current = node("root", { primary_image: asset("img-1") });
  const updated = node("root", { primary_image: undefined });

  const result = deriveMutationFromDifference(current, updated);

  assert.not.ok(result.clean);
  assert.equal(result.nodeMutation, { primary_image_asset_id: null });
});

test("deriveMutationFromDifference builds child mutations recursively", () => {
  const currentChild = node("child-1", { name: "child-old" });
  const updatedChild = node("child-1", { name: "child-new" });

  const current = node("root", { children: [currentChild] });
  const updated = node("root", { children: [updatedChild] });

  const result = deriveMutationFromDifference(current, updated);

  assert.not.ok(result.clean);
  assert.equal(result.nodeMutation, {});
  assert.equal(result.childMutation["child-1"], [{ name: "child-new" }]);
});

test("deriveMutationFromDifference omits fid for new child property schema fields", () => {
  const current = node("root", { child_property_schema: [] });
  const updated = node("root", {
    child_property_schema: [
      {
        fid: "new_field_1",
        name: "New Field",
        sort: "1",
        type: PropertyType.text,
      },
    ],
  });

  const result = deriveMutationFromDifference(current, updated);

  assert.not.ok(result.clean);
  assert.equal(result.childPropertySchemaMutation, [
    {
      name: "New Field",
      sort: "1",
      type: PropertyType.text,
    },
  ]);
});

test.run();
