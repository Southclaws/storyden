import { DndContext } from "@dnd-kit/core";
import {
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { PropsWithChildren } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { NodeWithChildren } from "@/api/openapi-schema";
import {
  LibraryPageBlockTypeDirectory,
  hydrateNode,
} from "@/lib/library/metadata";

import { NodeStoreAPI, createNodeStore } from "../../store";

import { DirectoryRichCard } from "./DirectoryRichCard";
import { LibraryPageDirectoryBlockGrid } from "./LibraryPageDirectoryBlockGrid";

let store: NodeStoreAPI;
vi.mock("../../Context", () => ({
  useLibraryPageContext: () => ({ nodeID: "parent", store }),
}));
vi.mock("../../useEditState", () => ({
  useEditState: () => ({ isDirectEditing: true }),
}));
vi.mock("@/components/library/LibraryPageMenu/LibraryPageMenu", () => ({
  LibraryPageMenu: ({ children }: PropsWithChildren) => children,
}));
vi.mock("next/navigation", () => ({
  useRouter: () => ({ prefetch: vi.fn() }),
}));

const schema = [
  { fid: "style", name: "Style", type: "text", sort: "0" },
] as const;
const child: NodeWithChildren = {
  id: "font",
  slug: "font",
  name: "Test font",
  description: "A test font",
  createdAt: "2026-11-01T00:00:00Z",
  updatedAt: "2026-11-01T00:00:00Z",
  owner: {
    id: "author",
    handle: "author",
    name: "Author",
    joined: "2026-11-01T00:00:00Z",
    roles: [],
  },
  visibility: "published",
  hide_child_tree: false,
  meta: {},
  assets: [],
  tags: [],
  children: [],
  child_property_schema: [],
  recomentations: [],
  properties: [{ ...schema[0], value: "Monospace" }],
};

describe("Directory card draft editing", () => {
  beforeEach(() => {
    const parent = hydrateNode({
      ...child,
      id: "parent",
      slug: "parent",
      name: "Fonts",
      children: [child],
      properties: [],
      child_property_schema: [...schema],
    });
    store = createNodeStore({
      original: parent,
      draft: structuredClone(parent),
    });
    vi.spyOn(console, "debug").mockImplementation(() => {});
  });

  it.each(["minimal", "row", "box", "fill"] as const)(
    "limits directory highlight targets to names in %s cards",
    (card_style) => {
      const { container } = render(
        <DirectoryRichCard
          node={{
            ...child,
            primary_image: {
              id: "cover",
              path: "/assets/font.png",
              filename: "font.png",
              mime_type: "image/png",
              width: 1200,
              height: 800,
            },
          }}
          layout="grid"
          config={{
            layout: "grid",
            card_style,
            columns: [
              { fid: "fixed:name", hidden: false },
              { fid: "fixed:description", hidden: false },
            ],
          }}
          columns={[]}
          values={{}}
          editing={false}
          onFieldValueChange={vi.fn()}
        />,
      );

      const description = screen.getByText("A test font");
      expect(description.closest("[data-search-highlight]")).toBeNull();
      expect(
        Array.from(container.querySelectorAll("[data-search-highlight]")).map(
          (element) => element.textContent,
        ),
      ).toEqual(["Test font"]);
    },
  );

  it.each(["rows", "grid", "gallery"] as const)(
    "keeps badge typing in the %s card while the server snapshot is unchanged",
    async (layout) => {
      const user = userEvent.setup();
      const block: LibraryPageBlockTypeDirectory = {
        type: "directory",
        config: {
          layout,
          card_style: "box",
          columns: [
            { fid: "fixed:name", hidden: false },
            { fid: "style", hidden: false, placement: "footer-start" },
          ],
        },
      };
      const card = () => (
        <DndContext>
          <LibraryPageDirectoryBlockGrid
            nodes={[child]}
            block={block}
            currentChildPropertySchema={[...schema]}
            layout={layout}
          />
        </DndContext>
      );
      const { rerender } = render(card());
      const badge = screen.getByRole("textbox", {
        name: "Style for Test font",
      });
      await user.click(badge);
      await user.keyboard("{End} display");

      expect(badge).toHaveValue("Monospace display");
      expect(store.getState().draft.children[0]?.properties[0]?.value).toBe(
        "Monospace display",
      );
      expect(child.properties[0]?.value).toBe("Monospace");
      rerender(card());
      expect(
        screen.getByRole("textbox", { name: "Style for Test font" }),
      ).toHaveValue("Monospace display");
    },
  );

  it.each([
    ["rows", "row"],
    ["grid", "minimal"],
    ["gallery", "minimal"],
    ["grid", "row"],
    ["grid", "box"],
    ["grid", "fill"],
    ["gallery", "box"],
  ] as const)(
    "keeps the full directory media while editing a lighter draft in %s/%s",
    async (layout, card_style) => {
      const user = userEvent.setup();
      const fullChild: NodeWithChildren = {
        ...child,
        content: "<p>Full page content from the directory entry.</p>",
        primary_image: {
          id: "cover",
          path: "/assets/font-cover",
          filename: "font.png",
          mime_type: "image/png",
          width: 1200,
          height: 800,
        },
      };
      const block: LibraryPageBlockTypeDirectory = {
        type: "directory",
        config: {
          layout,
          card_style,
          show_content: true,
          columns: [
            { fid: "fixed:name", hidden: false },
            { fid: "fixed:primary_image", hidden: false },
            { fid: "style", hidden: false, placement: "footer-start" },
          ],
        },
      };
      const { container } = render(
        <DndContext>
          <LibraryPageDirectoryBlockGrid
            nodes={[fullChild]}
            block={block}
            currentChildPropertySchema={[...schema]}
            layout={layout}
          />
        </DndContext>,
      );

      expect(container.querySelector("img")).not.toBeNull();
      let expectedOverlay = false;

      if (layout === "grid" && card_style === "fill") {
        expectedOverlay = true;
      }

      expect(
        screen.getByRole("button", { name: "Move Test font" }),
      ).toHaveAttribute("data-overlay", String(expectedOverlay));
      const editor = screen.getByRole("textbox", {
        name: "Style for Test font",
      });
      await user.click(editor);
      await user.keyboard("{End} display");
      expect(editor).toHaveValue("Monospace display");
      expect(container.querySelector("img")).not.toBeNull();

      expect(
        await screen.findByText("Full page content from the directory entry."),
      ).toBeInTheDocument();
    },
  );

  it.each(["grid", "gallery"] as const)(
    "uses Surface for Minimal %s, including missing and hidden covers",
    (layout) => {
      const config = {
        layout,
        card_style: "minimal",
        columns: [{ fid: "fixed:name", hidden: false }],
      } as const;
      function card(
        node: NodeWithChildren,
        columns: NonNullable<
          LibraryPageBlockTypeDirectory["config"]
        >["columns"] = [...config.columns],
      ) {
        return (
          <DndContext>
            <LibraryPageDirectoryBlockGrid
              nodes={[node]}
              block={{ type: "directory", config: { ...config, columns } }}
              currentChildPropertySchema={[...schema]}
              layout={layout}
            />
          </DndContext>
        );
      }
      const { container, rerender } = render(card(child));
      expect(
        screen.getByRole("article", { name: "Test font" }),
      ).toBeInTheDocument();
      expect(
        container.querySelector(".rich-card__root--shape_minimal"),
      ).not.toBeNull();
      expect(
        container.querySelector(".rich-card__mediaMissing svg"),
      ).not.toBeNull();
      expect(
        screen.getByRole("textbox", { name: "Name for Test font" }),
      ).toHaveValue("Test font");

      rerender(card(child, []));
      expect(
        container.querySelector(".rich-card__mediaMissing"),
      ).toHaveTextContent("Test font");
      expect(screen.queryByRole("heading")).toBeNull();

      const withCover: NodeWithChildren = {
        ...child,
        primary_image: {
          id: "cover",
          path: "/assets/font.png",
          filename: "font.png",
          mime_type: "image/png",
          width: 400,
          height: 1600,
        },
      };
      rerender(card(withCover, [{ fid: "fixed:primary_image", hidden: true }]));
      expect(container.querySelector("img")).toBeNull();
      expect(container.querySelector(".rich-card__mediaContainer")).toBeNull();

      rerender(card(withCover, []));
      const image = container.querySelector("img");
      expect(image).toHaveAttribute("loading", "lazy");
      expect(image).toHaveAttribute("width", "400");
      expect(image?.closest("a")).toBeNull();

      if (layout === "gallery") {
        expect(image).toHaveStyle({ objectFit: "contain" });
        expect(image?.parentElement).not.toHaveStyle({
          aspectRatio: "1.3333333333333333",
        });
      }

      if (layout === "grid") {
        expect(image).toHaveStyle({ objectFit: "cover" });
        expect(image?.parentElement).toHaveStyle({
          aspectRatio: "1.3333333333333333",
        });
      }
    },
  );

  it.each([
    ["rows", "row"],
    ["grid", "minimal"],
    ["grid", "box"],
    ["gallery", "box"],
  ] as const)(
    "edits details and configures properties in %s/%s",
    async (layout, card_style) => {
      const user = userEvent.setup();
      const block: LibraryPageBlockTypeDirectory = {
        type: "directory",
        config: {
          layout,
          card_style,
          columns: [
            { fid: "fixed:name", hidden: false },
            { fid: "style", hidden: false, placement: "details" },
          ],
        },
      };
      const draft = structuredClone(store.getState().draft);
      draft.meta.layout = { blocks: [block] };
      store.setState({ draft });

      render(
        <DndContext>
          <LibraryPageDirectoryBlockGrid
            nodes={[child]}
            block={block}
            currentChildPropertySchema={[...schema]}
            layout={layout}
          />
        </DndContext>,
      );

      const editor = screen.getByRole("textbox", {
        name: "Style for Test font",
      });
      await user.click(editor);
      await user.keyboard("{End} display");
      expect(editor).toHaveValue("Monospace display");
      expect(store.getState().draft.children[0]?.properties[0]?.value).toBe(
        "Monospace display",
      );

      const trigger = screen.getByRole("button", { name: "Configure Style" });
      expect(trigger.closest("dt")?.nextElementSibling?.tagName).toBe("DD");
      await user.click(trigger);
      const menu = screen.getByRole("menu");
      await waitFor(() => expect(menu).toHaveFocus());
      fireEvent.change(
        within(menu).getByRole("textbox", { name: "Name for Style" }),
        {
          target: { value: "Classification" },
        },
      );
      expect(store.getState().draft.child_property_schema[0]?.name).toBe(
        "Classification",
      );

      await user.keyboard("{Home}");
      await waitFor(() =>
        expect(
          screen.getByRole("menuitem", { name: /Hide column/ }),
        ).toHaveAttribute("data-highlighted"),
      );
      await user.keyboard("{Enter}");
      await waitFor(() =>
        expect(store.getState().draft.meta.layout?.blocks[0]).toMatchObject({
          config: {
            columns: expect.arrayContaining([
              { fid: "style", hidden: true, placement: "details" },
            ]),
          },
        }),
      );

      await user.click(trigger);
      await waitFor(() => expect(screen.getByRole("menu")).toHaveFocus());
      await user.keyboard("{End}");
      await waitFor(() =>
        expect(
          screen.getByRole("menuitem", { name: /Delete/ }),
        ).toHaveAttribute("data-highlighted"),
      );
      await user.keyboard("{Enter}");
      await waitFor(() =>
        expect(store.getState().draft.child_property_schema).toHaveLength(0),
      );
      expect(store.getState().draft.meta.layout?.blocks[0]).toMatchObject({
        config: { columns: [{ fid: "fixed:name", hidden: false }] },
      });
    },
  );
});
