import { beforeEach, describe, expect, it, vi } from "vitest";

import {
  Asset,
  LinkReference,
  NodeWithChildren,
  PropertySchema,
  PropertyType,
} from "@/api/openapi-schema";
import { hydrateNode } from "@/lib/library/metadata";

import { createNodeStore } from "./store";

let consoleWarn: ReturnType<typeof vi.spyOn>;

beforeEach(() => {
  vi.spyOn(console, "debug").mockImplementation(() => {});
  consoleWarn = vi.spyOn(console, "warn").mockImplementation(() => {});
});

describe("LibraryPageScreen store", () => {
  it("updates draft fields without mutating the original node", () => {
    const original = hydrateNode(node("root"));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    store.getState().setName("Renamed");
    store.getState().setSlug("renamed");
    store.getState().setContent("<p>updated</p>");
    store.getState().setTags(["alpha", "beta"]);
    store.getState().setLink(linkReference("https://storyden.org"));
    store.getState().removeLink();

    const state = store.getState();

    expect(state.draft.name).toBe("Renamed");
    expect(state.draft.slug).toBe("renamed");
    expect(state.draft.content).toBe("<p>updated</p>");
    expect(state.draft.tags.map((tag) => tag.name)).toEqual([
      "alpha",
      "beta",
    ]);
    expect(state.draft.link).toBeUndefined();

    expect(state.original.name).toBe("Node root");
    expect(state.original.slug).toBe("node-root");
    expect(state.original.content).toBe("<p>hello</p>");
    expect(state.original.tags.map((tag) => tag.name)).toEqual(["general"]);
    expect(state.original.link?.url).toBe("https://example.com/root");
  });

  it("manages cover image and asset draft state", () => {
    const original = hydrateNode(node("root", { primary_image: undefined }));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    const cover = asset("cover");
    const gallery = asset("gallery");

    store.getState().setPrimaryImage({
      asset: cover,
      isReplacement: false,
      config: {
        top: 10,
        left: 20,
      },
    });
    store.getState().addAsset(gallery);

    expect(store.getState().draft.primary_image).toEqual(cover);
    expect(store.getState().draft.meta.coverImage).toEqual({
      top: 10,
      left: 20,
    });
    expect(store.getState().draft.assets).toEqual([gallery]);

    store.getState().removePrimaryImage();
    store.getState().removeAsset(gallery);

    expect(store.getState().draft.primary_image).toBeUndefined();
    expect(store.getState().draft.meta.coverImage).toBeNull();
    expect(store.getState().draft.assets).toEqual([]);
  });

  it("manages string properties by name and ID", () => {
    const original = hydrateNode(node("root"));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    store.getState().addProperty("Label", PropertyType.text, "second");

    const added = store
      .getState()
      .draft.properties.find((property) => property.name === "Label 1");
    expect(added).toMatchObject({
      name: "Label 1",
      type: PropertyType.text,
      sort: "5",
      value: "second",
    });

    store.getState().setPropertyName("Label 1", "Subtitle");
    store.getState().setPropertyValue("Subtitle", "new value");
    store.getState().removePropertyByName("Label");

    expect(store.getState().draft.properties).toEqual([
      expect.objectContaining({
        name: "Subtitle",
        value: "new value",
      }),
    ]);

    store.getState().removePropertyByID(added!.fid);

    expect(store.getState().draft.properties).toEqual([]);
  });

  it("manages child property schema and directory layout columns together", () => {
    const original = hydrateNode(
      node("root", {
        meta: {
          layout: {
            blocks: [
              {
                type: "directory",
                config: {
                  layout: "table",
                  columns: [{ fid: "field-1", hidden: false }],
                },
              },
            ],
          },
        },
      }),
    );
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    const newSchema: PropertySchema = {
      fid: "field-2",
      name: "Rating",
      sort: "1",
      type: PropertyType.text,
    };

    store.getState().addChildProperty(newSchema);
    store.getState().setChildPropertyName("field-2", "Score");
    store.getState().setChildPropertyHiddenState("field-2", true);

    expect(store.getState().draft.child_property_schema).toEqual([
      expect.objectContaining({ fid: "field-1", name: "Label" }),
      expect.objectContaining({ fid: "field-2", name: "Score" }),
    ]);
    expect(directoryColumns(store.getState().draft)).toEqual([
      { fid: "field-1", hidden: false },
      { fid: "field-2", hidden: true },
    ]);

    store.getState().removeChildPropertyByID("field-2");

    expect(store.getState().draft.child_property_schema).toEqual([
      expect.objectContaining({ fid: "field-1" }),
    ]);
    expect(directoryColumns(store.getState().draft)).toEqual([
      { fid: "field-1", hidden: false },
    ]);
  });

  it("updates child fixed fields and child property values", () => {
    const child = node("child", {
      name: "Child",
      description: "child description",
      link: linkReference("https://example.com/child"),
      properties: [
        {
          fid: "field-1",
          name: "Label",
          sort: "0",
          type: PropertyType.text,
          value: "old child value",
        },
      ],
    });
    const original = hydrateNode(node("root", { children: [child] }));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    store.getState().setChildPropertyValue("child", "fixed:name", "New child");
    store
      .getState()
      .setChildPropertyValue("child", "fixed:description", "New description");
    store
      .getState()
      .setChildPropertyValue("child", "fixed:link", "https://storyden.org");
    store
      .getState()
      .setChildPropertyValue("child", "field-1", "new child value");

    const updatedChild = store.getState().draft.children[0]!;

    expect(updatedChild.name).toBe("New child");
    expect(updatedChild.description).toBe("New description");
    expect(updatedChild.link?.url).toBe("https://storyden.org");
    expect(updatedChild.properties[0]?.value).toBe("new child value");

    store.getState().setChildPropertyValue("missing", "field-1", "ignored");
    store.getState().setChildPropertyValue("child", "missing", "ignored");

    expect(consoleWarn).toHaveBeenCalled();
  });

  it("manages library page layout blocks", () => {
    const original = hydrateNode(node("root"));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });

    store.getState().addBlock("tags", 1);
    store.getState().moveBlock("tags", 0);
    store.getState().overwriteBlock({
      type: "assets",
      config: {
        layout: "grid",
        gridSize: 3,
      },
    });
    store.getState().removeBlock("link");

    const blocks = store.getState().draft.meta.layout?.blocks;

    expect(blocks?.map((block) => block.type)).toEqual([
      "tags",
      "cover",
      "title",
      "content",
    ]);
    expect(blocks?.find((block) => block.type === "link")).toBeUndefined();

    store.getState().addBlock("assets");
    store.getState().overwriteBlock({
      type: "assets",
      config: {
        layout: "grid",
        gridSize: 3,
      },
    });

    expect(
      store.getState().draft.meta.layout?.blocks.find(
        (block) => block.type === "assets",
      ),
    ).toEqual({
      type: "assets",
      config: {
        layout: "grid",
        gridSize: 3,
      },
    });
  });

  it("skips commit callbacks when the draft is clean", async () => {
    const original = hydrateNode(node("root"));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });
    const callback = vi.fn();

    await store.getState().commit(callback);

    expect(callback).not.toHaveBeenCalled();
  });

  it("commits derived mutations and replaces original and draft with the saved node", async () => {
    const original = hydrateNode(node("root"));
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });
    const updated = hydrateNode(node("root", { name: "Saved" }));
    const callback = vi.fn().mockResolvedValue(updated);

    store.getState().setName("Saved");

    await store.getState().commit(callback);

    expect(callback).toHaveBeenCalledWith(
      expect.objectContaining({
        clean: false,
        nodeMutation: {
          name: "Saved",
        },
        childMutation: {},
        childPropertySchemaMutation: undefined,
      }),
    );
    expect(store.getState().original.name).toBe("Saved");
    expect(store.getState().draft.name).toBe("Saved");
  });

  it("commits child mutations and child property schema changes", async () => {
    const original = hydrateNode(
      node("root", {
        children: [node("child", { name: "Old child" })],
      }),
    );
    const store = createNodeStore({
      original,
      draft: structuredClone(original),
    });
    const updated = hydrateNode(
      node("root", {
        children: [node("child", { name: "New child" })],
        child_property_schema: [
          {
            fid: "field-1",
            name: "Renamed label",
            sort: "0",
            type: PropertyType.text,
          },
        ],
      }),
    );
    const callback = vi.fn().mockResolvedValue(updated);

    store.getState().setChildPropertyName("field-1", "Renamed label");
    store.getState().setChildPropertyValue("child", "fixed:name", "New child");

    await store.getState().commit(callback);

    expect(callback).toHaveBeenCalledWith(
      expect.objectContaining({
        clean: false,
        nodeMutation: {},
        childMutation: {
          child: [
            {
              name: "New child",
            },
          ],
        },
        childPropertySchemaMutation: [
          {
            fid: "field-1",
            name: "Renamed label",
            sort: "0",
            type: PropertyType.text,
          },
        ],
      }),
    );
    expect(store.getState().draft.child_property_schema[0]?.name).toBe(
      "Renamed label",
    );
    expect(store.getState().draft.children[0]?.name).toBe("New child");
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
    owner: profile("author"),
    parent: undefined,
    primary_image: asset(`cover-${id}`),
    assets: [],
    tags: [{ name: "general", colour: "white", item_count: 1 }],
    link: linkReference(`https://example.com/${id}`),
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
    child_property_schema: [
      {
        fid: "field-1",
        name: "Label",
        sort: "0",
        type: PropertyType.text,
      },
    ],
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
    id: `link-${url}`,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    domain: "example.com",
    slug: "example",
    url,
  };
}

function directoryColumns(node: NodeWithChildren) {
  const block = (node.meta as any).layout?.blocks.find(
    (block) => block.type === "directory",
  );

  if (!block || block.type !== "directory") {
    return [];
  }

  return block.config?.columns ?? [];
}
