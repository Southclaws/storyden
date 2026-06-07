import { describe, expect, it } from "vitest";

import {
  NodeVersion,
  NodeWithChildren,
  PropertyType,
} from "@/api/openapi-schema";
import { hydrateNode } from "@/lib/library/metadata";

import { buildNodeVersionMutation, overlayNodeVersion } from "./versionedEdit";

describe("versionedEdit helpers", () => {
  it("builds node version patches from only version-supported node fields", () => {
    const patch = buildNodeVersionMutation({
      clean: false,
      nodeMutation: {
        name: "Name",
        slug: "slug",
        description: "Description",
        content: "<p>body</p>",
        properties: [
          {
            fid: "field-1",
            name: "Label",
            value: "Value",
            type: PropertyType.text,
          },
        ],
        tags: ["ignored"],
        hide_child_tree: true,
      },
      childMutation: {
        child: [{ name: "ignored child" }],
      },
      childPropertySchemaMutation: [
        {
          fid: "ignored",
          name: "Ignored",
          sort: "0",
          type: PropertyType.text,
        },
      ],
    });

    expect(patch).toEqual({
      name: "Name",
      slug: "slug",
      description: "Description",
      content: "<p>body</p>",
      properties: [
        {
          fid: "field-1",
          name: "Label",
          value: "Value",
          type: PropertyType.text,
        },
      ],
    });
  });

  it("overlays proposed version fields onto a hydrated node without mutating the original", () => {
    const original = hydrateNode(node("root"));
    const version = {
      id: "version-1",
      node_id: "root",
      author: profile("author"),
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
      status: "draft",
      name: "Proposed name",
      slug: "proposed-slug",
      description: "Proposed description",
      content: "<p>proposed</p>",
      meta: {},
      properties: [
        {
          name: "Release year",
          value: "1992",
        },
      ],
    } satisfies NodeVersion;

    const draft = overlayNodeVersion(original, version);

    expect(draft).toMatchObject({
      name: "Proposed name",
      slug: "proposed-slug",
      description: "Proposed description",
      content: "<p>proposed</p>",
      properties: [
        {
          fid: "new_field_0",
          name: "Release year",
          sort: "0",
          type: PropertyType.text,
          value: "1992",
        },
      ],
    });
    expect(draft.meta.layout?.blocks.map((block) => block.type)).toEqual([
      "cover",
      "title",
      "content",
      "link",
    ]);
    expect(original.name).toBe("Node root");
    expect(original.properties[0]?.value).toBe("old value");
  });

  it("clears nullable fields and overlays full snapshot fields", () => {
    const original = hydrateNode(node("root"));
    const version = {
      id: "version-1",
      node_id: "root",
      author: profile("author"),
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
      status: "draft",
      name: "Snapshot name",
      slug: "snapshot-slug",
      description: null,
      content: null,
      meta: {},
      properties: [],
    } satisfies NodeVersion;

    const draft = overlayNodeVersion(original, version);

    expect(draft.name).toBe("Snapshot name");
    expect(draft.slug).toBe("snapshot-slug");
    expect(draft.description).toBe("");
    expect(draft.content).toBeUndefined();
    expect(draft.properties).toEqual([]);
  });
});

function node(
  id: string,
  overrides: Partial<NodeWithChildren> = {},
): NodeWithChildren {
  return {
    id,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    name: `Node ${id}`,
    slug: `node-${id}`,
    description: "description",
    content: "<p>hello</p>",
    owner: profile("owner"),
    assets: [],
    tags: [],
    visibility: "published",
    hide_child_tree: false,
    meta: {},
    properties: [
      {
        fid: "field-1",
        name: "Label",
        sort: "0",
        type: PropertyType.text,
        value: "old value",
      },
    ],
    child_property_schema: [],
    children: [],
    recomentations: [],
    ...overrides,
  };
}

function profile(id: string) {
  return {
    id,
    handle: id,
    name: `Member ${id}`,
    joined: "2026-01-01T00:00:00Z",
    roles: [],
  };
}
