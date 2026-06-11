import { act, render } from "@testing-library/react";
import { PropsWithChildren } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import {
  Asset,
  LinkReference,
  NodeWithChildren,
  PropertyType,
} from "@/api/openapi-schema";

import { LibraryPageProvider, useLibraryPageContext } from "./Context";
import { LibraryPageAutosaveController } from "./LibraryPageAutosaveController";
import { NodeStoreAPI } from "./store";

const mocks = vi.hoisted(() => ({
  nodeUpdate: vi.fn(),
  nodeUpdateChildrenPropertySchema: vi.fn(),
  nodeVersionCreate: vi.fn(),
  nodeVersionUpdate: vi.fn(),
  revalidate: vi.fn(),
}));

vi.mock("@/api/openapi-client/nodes", () => ({
  nodeUpdate: mocks.nodeUpdate,
  nodeUpdateChildrenPropertySchema: mocks.nodeUpdateChildrenPropertySchema,
  nodeVersionCreate: mocks.nodeVersionCreate,
  nodeVersionUpdate: mocks.nodeVersionUpdate,
}));

vi.mock("@/lib/library/library", () => ({
  useLibraryMutation: () => ({
    revalidate: mocks.revalidate,
  }),
}));

vi.mock("./useEditState", () => ({
  useEditState: () => ({
    editMode: "direct",
    proposalVersion: undefined,
    setProposalVersion: vi.fn(),
  }),
}));

const consoleDebug = vi.spyOn(console, "debug").mockImplementation(() => {});
const consoleWarn = vi.spyOn(console, "warn").mockImplementation(() => {});

describe("LibraryPageProvider direct autosave", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    mocks.nodeUpdateChildrenPropertySchema.mockResolvedValue(undefined);
    mocks.revalidate.mockResolvedValue(undefined);
  });

  afterEach(() => {
    vi.runOnlyPendingTimers();
    vi.useRealTimers();
  });

  it("debounces draft changes and persists node mutations through the direct node update endpoint", async () => {
    const root = node("root");
    const updated = node("root", { name: "Renamed" });
    mocks.nodeUpdate.mockResolvedValue(updated);

    const store = await renderProvider(root);

    act(() => {
      store.getState().setName("Renamed");
    });

    expect(mocks.nodeUpdate).not.toHaveBeenCalled();

    await flushAutosave();

    expect(mocks.nodeUpdate).toHaveBeenCalledTimes(1);
    expect(mocks.nodeUpdate).toHaveBeenCalledWith("root", {
      name: "Renamed",
    });
    expect(mocks.revalidate).toHaveBeenCalledWith(updated);
    expect(store.getState().original.name).toBe("Renamed");
    expect(store.getState().draft.name).toBe("Renamed");
  });

  it("does not call update endpoints when the derived mutation is clean", async () => {
    const root = node("root", { slug: "valid-slug" });
    const store = await renderProvider(root);

    act(() => {
      store.getState().setSlug("invalid slug");
    });

    await flushAutosave();

    expect(mocks.nodeUpdate).not.toHaveBeenCalled();
    expect(mocks.nodeUpdateChildrenPropertySchema).not.toHaveBeenCalled();
    expect(mocks.revalidate).not.toHaveBeenCalled();
  });

  it("persists child property schema and child node mutations before revalidating the parent node", async () => {
    const child = node("child", {
      name: "Old child",
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
    const root = node("root", { children: [child] });
    mocks.nodeUpdate.mockImplementation(async (id: string, mutation: any) => {
      if (id === "child") {
        return node("child", {
          name: mutation.name ?? "New child",
          properties: mutation.properties ?? child.properties,
        });
      }

      return node("root", {
        children: [
          node("child", {
            name: "New child",
            properties: [
              {
                fid: "field-1",
                name: "Label",
                sort: "0",
                type: PropertyType.text,
                value: "new child value",
              },
            ],
          }),
        ],
        child_property_schema: [
          {
            fid: "field-1",
            name: "Renamed label",
            sort: "0",
            type: PropertyType.text,
          },
        ],
      });
    });

    const store = await renderProvider(root);

    act(() => {
      store.getState().setChildPropertyName("field-1", "Renamed label");
      store
        .getState()
        .setChildPropertyValue("child", "fixed:name", "New child");
      store
        .getState()
        .setChildPropertyValue("child", "field-1", "new child value");
    });

    await flushAutosave();

    expect(mocks.nodeUpdateChildrenPropertySchema).toHaveBeenCalledWith(
      "root",
      [
        {
          fid: "field-1",
          name: "Renamed label",
          sort: "0",
          type: PropertyType.text,
        },
      ],
    );
    expect(mocks.nodeUpdate).toHaveBeenCalledWith("child", {
      name: "New child",
      properties: [
        {
          fid: "field-1",
          name: "Label",
          value: "new child value",
          type: PropertyType.text,
        },
      ],
    });
    expect(mocks.nodeUpdate).toHaveBeenCalledWith("root", {});
    expect(mocks.revalidate).toHaveBeenCalledTimes(1);
  });

  it("replaces the browser URL after a slug autosave succeeds", async () => {
    const root = node("root", { slug: "old-slug" });
    const updated = node("root", { slug: "new-slug" });
    const replaceState = vi
      .spyOn(window.history, "replaceState")
      .mockImplementation(() => {});
    mocks.nodeUpdate.mockResolvedValue(updated);

    const store = await renderProvider(root);

    act(() => {
      store.getState().setSlug("new-slug");
    });

    await flushAutosave();

    expect(mocks.nodeUpdate).toHaveBeenCalledWith("root", {
      slug: "new-slug",
    });
    expect(replaceState).toHaveBeenCalledWith(
      null,
      "",
      "/l/new-slug?edit=true",
    );

    replaceState.mockRestore();
  });
});

afterEach(() => {
  consoleDebug.mockClear();
  consoleWarn.mockClear();
});

function ProviderHarness({
  node,
  children,
}: PropsWithChildren<{ node: NodeWithChildren }>) {
  return (
    <LibraryPageProvider node={node}>
      <LibraryPageAutosaveController />
      {children}
    </LibraryPageProvider>
  );
}

function StoreProbe({ onStore }: { onStore: (store: NodeStoreAPI) => void }) {
  const { store } = useLibraryPageContext();

  onStore(store);

  return null;
}

async function renderProvider(node: NodeWithChildren) {
  let store: NodeStoreAPI | undefined;

  render(
    <ProviderHarness node={node}>
      <StoreProbe onStore={(s) => (store = s)} />
    </ProviderHarness>,
  );

  expect(store).toBeDefined();

  return store!;
}

async function flushAutosave() {
  await act(async () => {
    await vi.advanceTimersByTimeAsync(500);
  });
  await Promise.resolve();
}

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
