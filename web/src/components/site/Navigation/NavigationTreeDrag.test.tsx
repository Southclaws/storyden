import { act, render } from "@testing-library/react";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, test, vi } from "vitest";

import type { Category, NodeWithChildren } from "@/api/openapi-schema";
import { DndProvider } from "@/lib/dragdrop/provider";

type DndContextProps = {
  onDragEnd?: (event: Record<string, unknown>) => Promise<void>;
};

const mocks = vi.hoisted(() => ({
  contextProps: null as DndContextProps | null,
  emitCategory: vi.fn(),
  moveNode: vi.fn(),
  revalidate: vi.fn(),
}));

vi.mock("@dnd-kit/core", () => ({
  DndContext: ({ children, ...props }: { children: ReactNode }) => {
    mocks.contextProps = props;
    return children;
  },
  DragOverlay: ({ children }: { children: ReactNode }) => children,
  MouseSensor: class MouseSensor {},
  TouchSensor: class TouchSensor {},
  getClientRect: vi.fn(),
  pointerWithin: vi.fn(() => []),
  rectIntersection: vi.fn(() => []),
  useSensor: vi.fn(() => ({})),
  useSensors: vi.fn(() => []),
}));

vi.mock("@/api/client", () => ({
  handle: async (
    operation: () => Promise<void>,
    options?: { cleanup?: () => Promise<void> },
  ) => {
    await operation();
    await options?.cleanup?.();
  },
}));

vi.mock("@/lib/library/library", () => ({
  useLibraryMutation: () => ({
    moveNode: mocks.moveNode,
    revalidate: mocks.revalidate,
  }),
}));

vi.mock("@/lib/category/events", () => ({
  useEmitCategoryEvent: () => mocks.emitCategory,
}));

vi.mock("@/lib/feed/events", () => ({
  useEmitFeedBlockEvent: () => vi.fn(),
}));

vi.mock("@/lib/library/events", () => ({
  useEmitLibraryBlockEvent: () => vi.fn(),
}));

vi.mock("@/lib/navigation/events", () => ({
  useEmitNavigationItemEvent: () => vi.fn(),
}));

function category(id: string): Category {
  return { id, slug: id, name: id, children: [] } as unknown as Category;
}

function page(id: string): NodeWithChildren {
  return {
    id,
    slug: id,
    name: id,
    children: [],
  } as unknown as NodeWithChildren;
}

function getDragEnd() {
  const onDragEnd = mocks.contextProps?.onDragEnd;
  if (!onDragEnd) {
    throw new Error("DndProvider did not register onDragEnd");
  }
  return onDragEnd;
}

beforeEach(() => {
  mocks.contextProps = null;
});

describe("sidebar tree drop behavior", () => {
  test("category dividers preserve sibling order and parent placement", async () => {
    render(<DndProvider>navigation</DndProvider>);

    await act(() =>
      getDragEnd()({
        active: {
          data: {
            current: { type: "category", category: category("music") },
          },
        },
        over: {
          data: {
            current: {
              type: "category-divider",
              direction: "above",
              siblingCategoryID: "engineering",
              parentID: "technology",
            },
          },
        },
      }),
    );

    expect(mocks.emitCategory).toHaveBeenCalledWith(
      "category:reorder-category",
      {
        categorySlug: "music",
        targetCategory: "engineering",
        direction: "above",
        newParent: "technology",
      },
    );
  });

  test("dropping a category on another category reparents it", async () => {
    render(<DndProvider>navigation</DndProvider>);

    await act(() =>
      getDragEnd()({
        active: {
          data: {
            current: { type: "category", category: category("music") },
          },
        },
        over: {
          data: {
            current: {
              type: "category",
              categoryID: "technology",
            },
          },
        },
      }),
    );

    expect(mocks.emitCategory).toHaveBeenCalledWith(
      "category:reorder-category",
      {
        categorySlug: "music",
        targetCategory: "technology",
        direction: "inside",
        newParent: "technology",
      },
    );
  });

  test("library dividers preserve order and the destination parent", async () => {
    render(<DndProvider>navigation</DndProvider>);

    await act(() =>
      getDragEnd()({
        active: {
          data: {
            current: {
              type: "node",
              node: page("handbook"),
              parentID: null,
              context: "sidebar",
            },
          },
        },
        over: {
          data: {
            current: {
              type: "divider",
              direction: "below",
              siblingNode: page("api"),
              parentID: "documentation",
              context: "sidebar",
            },
          },
        },
      }),
    );

    expect(mocks.moveNode).toHaveBeenCalledWith(
      "handbook",
      "api",
      "below",
      "documentation",
      "documentation",
    );
    expect(mocks.revalidate).toHaveBeenCalledOnce();
  });

  test("dropping a library page on another page reparents it", async () => {
    render(<DndProvider>navigation</DndProvider>);

    await act(() =>
      getDragEnd()({
        active: {
          data: {
            current: {
              type: "node",
              node: page("handbook"),
              parentID: null,
              context: "sidebar",
            },
          },
        },
        over: {
          data: {
            current: {
              type: "node",
              node: page("documentation"),
              parentID: null,
              context: "sidebar",
            },
          },
        },
      }),
    );

    expect(mocks.moveNode).toHaveBeenCalledWith(
      "handbook",
      "documentation",
      "inside",
      "documentation",
      undefined,
    );
  });
});
