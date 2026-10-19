import { fireEvent, render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, test, vi } from "vitest";

import type {
  Account,
  Category,
  CategoryListOKResponse,
  NodeListResult,
  NodeWithChildren,
} from "@/api/openapi-schema";
import { Visibility } from "@/api/openapi-schema";
import { CategoryListTree } from "@/components/category/CategoryList/CategoryList";
import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { LibraryNavigationTree } from "@/components/site/Navigation/LibraryNavigationTree/LibraryNavigationTree";
import { libraryNavigationVisibility } from "@/components/site/Navigation/LibraryNavigationTree/LibraryNavigationTree.constants";
import { LinkButton } from "@/components/ui/link-button";
import { Card } from "@/components/ui/surface";

const router = vi.hoisted(() => ({ prefetch: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => router,
  usePathname: () => "/l",
  useParams: () => ({}),
}));

const dnd = vi.hoisted(() => ({
  draggables: [] as Array<Record<string, unknown>>,
  droppables: [] as Array<Record<string, unknown>>,
}));

const api = vi.hoisted(() => ({
  useNodeList: vi.fn(),
}));

vi.mock("@/api/openapi-client/nodes", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@/api/openapi-client/nodes")>()),
  useNodeList: api.useNodeList,
}));

vi.mock("@dnd-kit/core", () => ({
  useDraggable: (options: Record<string, unknown>) => {
    dnd.draggables.push(options);
    return {
      active: null,
      attributes: {},
      isDragging: false,
      listeners: {},
      over: null,
      setNodeRef: vi.fn(),
      transform: null,
    };
  },
  useDroppable: (options: Record<string, unknown>) => {
    dnd.droppables.push(options);
    return {
      setNodeRef: vi.fn(),
    };
  },
}));

vi.mock("@/auth", () => ({
  useSession: (initialSession?: Account) => initialSession,
}));

vi.mock("@/lib/category/events", () => ({
  useCategoryEvent: vi.fn(),
}));

vi.mock("@/components/category/CategoryMenu/CategoryMenu", () => ({
  CategoryMenu: ({ category }: { category: Category }) => (
    <button type="button">Manage {category.name}</button>
  ),
}));

vi.mock("@/components/category/CategoryCreate/CategoryCreateTrigger", () => ({
  CategoryCreateTrigger: () => <button type="button">Create category</button>,
}));

vi.mock("@/components/library/CreatePage", () => ({
  CreatePageAction: ({ parentSlug }: { parentSlug?: string }) => (
    <button type="button">
      {parentSlug ? `Create page under ${parentSlug}` : "Create page"}
    </button>
  ),
}));

vi.mock("@/components/library/LibraryPageMenu/LibraryPageMenu", () => ({
  LibraryPageMenu: ({ node }: { node: NodeWithChildren }) => (
    <button type="button">Manage {node.name}</button>
  ),
}));

function category(
  id: string,
  name: string,
  sort: number,
  parent?: string,
): Category {
  return {
    id,
    name,
    slug: id,
    sort,
    parent,
    children: [],
    colour: "#31a89f",
    description: "",
    postCount: 0,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  } as Category;
}

function page(
  id: string,
  name: string,
  children: NodeWithChildren[] = [],
  visibility: Visibility = Visibility.published,
): NodeWithChildren {
  return {
    id,
    name,
    slug: id,
    visibility,
    children,
  } as NodeWithChildren;
}

const categoryManager = {
  roles: [{ permissions: ["MANAGE_CATEGORIES"] }],
} as Account;

beforeEach(() => {
  dnd.draggables.length = 0;
  dnd.droppables.length = 0;
  api.useNodeList.mockImplementation((_params, options) => ({
    data: options?.swr?.fallbackData,
  }));
});

describe("sidebar navigation trees", () => {
  test("content cards warm a destination once across hover and keyboard focus", () => {
    const { rerender } = render(
      <Card
        id="content"
        title="Handbook"
        text="Preview"
        url="/l/handbook"
        prefetchOnIntent
      />,
    );

    expect(router.prefetch).not.toHaveBeenCalled();
    expect(
      screen.queryByRole("link", { name: "Preview" }),
    ).not.toBeInTheDocument();
    fireEvent.mouseEnter(screen.getByRole("link", { name: "Handbook" }));
    fireEvent.focus(screen.getByRole("link", { name: "Handbook" }));
    fireEvent.mouseEnter(screen.getByRole("link", { name: "Handbook" }));
    expect(router.prefetch).toHaveBeenCalledExactlyOnceWith("/l/handbook");

    rerender(
      <Card
        id="content"
        title="Discussion"
        text="Preview"
        url="/t/discussion"
        prefetchOnIntent
      />,
    );
    fireEvent.focus(screen.getByRole("link", { name: "Discussion" }));
    expect(router.prefetch).toHaveBeenCalledTimes(2);
    expect(router.prefetch).toHaveBeenLastCalledWith("/t/discussion");
  });

  test.each([
    { url: "/m/member", prefetchOnIntent: false },
    { url: "/l/disabled", prefetchOnIntent: true, disableAnchors: true },
    { url: "https://example.com/article", prefetchOnIntent: true },
  ])("cards do not prefetch ineligible destinations: $url", (props) => {
    render(<Card id="content" title="Content" {...props} />);
    const link = screen.getByRole("link", { name: "Content" });
    fireEvent.mouseEnter(link);
    fireEvent.focus(link);
    expect(router.prefetch).not.toHaveBeenCalled();
  });

  test.each([false, true, undefined])(
    "LinkButton forwards prefetch=%s to Next Link",
    (prefetch) => {
      const link = LinkButton({ href: "/l", children: "Library", prefetch });

      expect(link.props.prefetch).toBe(prefetch);
      expect(link.props.href).toBe("/l");
    },
  );

  test("category navigation renders and expands nested category routes", () => {
    const categories = [
      category("engineering", "Engineering", 1),
      category("frontend", "Frontend", 1, "engineering"),
      category("community", "Community", 2),
    ];

    render(
      <CategoryListTree
        categories={categories}
        currentPath="/d/frontend"
        initialSession={categoryManager}
        mutate={vi.fn() as never}
      />,
    );

    expect(screen.getByRole("link", { name: "Engineering" })).toHaveAttribute(
      "href",
      "/d/engineering",
    );
    expect(screen.getByRole("treeitem", { name: "Frontend" })).toHaveAttribute(
      "href",
      "/d/frontend",
    );
    expect(screen.getByRole("treeitem", { name: "Frontend" })).toHaveAttribute(
      "aria-current",
      "page",
    );
    expect(
      screen.getByRole("button", { name: "Create category" }),
    ).toBeVisible();
    expect(
      screen.getByRole("button", { name: "Manage Frontend" }),
    ).toBeVisible();

    expect(dnd.draggables).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: "engineering",
          data: expect.objectContaining({
            type: "category",
            categoryID: "engineering",
          }),
        }),
        expect.objectContaining({
          id: "frontend",
          data: expect.objectContaining({
            type: "category",
            categoryID: "frontend",
          }),
        }),
      ]),
    );
    expect(dnd.droppables).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: "frontend_above",
          data: expect.objectContaining({
            type: "category-divider",
            siblingCategoryID: "frontend",
            parentID: "engineering",
          }),
        }),
      ]),
    );
  });

  test("library navigation expands the current page ancestry and preserves node drag data", () => {
    const nodes = [
      page("handbook", "Handbook", [
        page("engineering", "Engineering", [page("frontend", "Frontend")]),
      ]),
      page("conduct", "Code of conduct"),
    ];

    render(
      <LibraryPageTree
        nodes={nodes}
        currentPath="/l/frontend"
        canManageLibrary
      />,
    );

    expect(screen.getByRole("link", { name: "Handbook" })).toHaveAttribute(
      "href",
      "/l/handbook",
    );
    expect(screen.getByRole("treeitem", { name: "Frontend" })).toHaveAttribute(
      "href",
      "/l/frontend",
    );
    expect(screen.getByRole("treeitem", { name: "Frontend" })).toHaveAttribute(
      "aria-current",
      "page",
    );
    expect(
      screen.getByRole("button", { name: "Create page under handbook" }),
    ).toBeVisible();
    expect(
      screen.getByRole("button", { name: "Manage Frontend" }),
    ).toBeVisible();

    expect(dnd.draggables).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: "handbook",
          data: expect.objectContaining({
            type: "node",
            parentID: null,
            context: "sidebar",
          }),
        }),
        expect.objectContaining({
          id: "frontend",
          data: expect.objectContaining({
            type: "node",
            parentID: "engineering",
            context: "sidebar",
          }),
        }),
      ]),
    );
    expect(dnd.droppables).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          id: "frontend_above",
          data: expect.objectContaining({
            type: "divider",
            parentID: "engineering",
            context: "sidebar",
          }),
        }),
      ]),
    );
  });

  test("library navigation renders authenticated draft sections from server fallback data", () => {
    const initialNodeList = {
      nodes: [
        page("handbook", "Handbook"),
        page("work-in-progress", "Work in progress", [], Visibility.draft),
      ],
    } as NodeListResult;

    render(
      <LibraryNavigationTree
        initialNodeList={initialNodeList}
        initialSession={categoryManager}
        visibility={libraryNavigationVisibility}
      />,
    );

    expect(
      screen.getByRole("heading", { name: "Library", level: 2 }),
    ).toBeVisible();
    expect(
      screen.getByRole("heading", { name: "Drafts", level: 3 }),
    ).toBeVisible();
    expect(screen.queryByRole("heading", { level: 1 })).not.toBeInTheDocument();
    expect(
      screen.getByRole("treeitem", { name: "Work in progress" }),
    ).toHaveAttribute("href", "/l/work-in-progress");
    expect(api.useNodeList).toHaveBeenCalledWith(
      { visibility: libraryNavigationVisibility },
      { swr: { fallbackData: initialNodeList } },
    );
  });
});
