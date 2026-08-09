import { describe, expect, test, vi } from "vitest";

import type { Account } from "@/api/openapi-schema";

import { libraryNavigationVisibility } from "../LibraryNavigationTree/LibraryNavigationTree.constants";

import { NavigationPane } from "./NavigationPane";

const server = vi.hoisted(() => ({
  categoryListCached: vi.fn(async () => ({ data: { categories: [] } })),
  nodeListCached: vi.fn(async () => ({ data: { nodes: [] } })),
}));

vi.mock("@/lib/category/server-category-list", () => ({
  categoryListCached: server.categoryListCached,
}));

vi.mock("@/lib/library/server-node-list", () => ({
  nodeListCached: server.nodeListCached,
}));

describe("NavigationPane", () => {
  test("requests the authenticated library tree during server rendering", async () => {
    await NavigationPane({ initialSession: {} as Account });

    expect(server.nodeListCached).toHaveBeenCalledWith({
      visibility: libraryNavigationVisibility,
    });
  });

  test("requests the public library tree for guests", async () => {
    await NavigationPane({});

    expect(server.nodeListCached).toHaveBeenCalledWith(undefined);
  });
});
